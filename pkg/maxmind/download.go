package maxmind

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/docker/go-units"
	"github.com/mholt/archiver/v3"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// DownloadDB download a database
func (c *Client) DownloadDB(eid EditionID, dlDir string) ([]os.FileInfo, error) {
	var err error

	// Check download directory
	if dlDir == "" {
		dlDir = filepath.Dir(os.Args[0])
	}
	dlDir, err = filepath.Abs(path.Clean(dlDir))
	if err != nil {
		return nil, errors.Wrap(err, "Cannot get absolute path of download directory")
	}
	if err := os.MkdirAll(dlDir, 0755); err != nil {
		return nil, errors.Wrap(err, "Cannot create download directory")
	}
	if err := isDirWriteable(dlDir); err != nil {
		return nil, errors.Wrap(err, "Download directory is not writable")
	}

	// Retrieve expected hash
	expHash, err := c.expectedHash(eid)
	if err != nil {
		return nil, errors.Wrap(err, "Cannot get archive MD5 hash")
	}

	// Download DB archive
	archive := path.Join(c.workDir, eid.Filename())
	if err := c.downloadArchive(eid, expHash, archive); err != nil {
		return nil, err
	}

	// Create MD5 file
	md5file := path.Join(c.workDir, fmt.Sprintf(".%s.%s", eid.Filename(), "md5"))
	if err := createFile(md5file, expHash); err != nil {
		return nil, errors.Errorf("Cannot create MD5 file %s", md5file)
	}

	// Extract DB from archive
	dbs, err := c.extractArchive(eid, archive, dlDir)
	if err != nil {
		return nil, err
	}

	return dbs, nil
}

func (c *Client) expectedHash(eid EditionID) (string, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/app/geoip_download", c.baseURL), nil)
	if err != nil {
		return "", errors.Wrap(err, "Request failed")
	}

	q := req.URL.Query()
	q.Add("license_key", c.licenseKey)
	q.Add("edition_id", eid.String())
	q.Add("suffix", fmt.Sprintf("%s.md5", eid.Suffix().String()))
	req.URL.RawQuery = q.Encode()

	if c.userAgent != "" {
		req.Header.Add("User-Agent", c.userAgent)
	}

	res, err := c.http.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", errors.Errorf("Received invalid status code %d: %s", res.StatusCode, res.Body)
	}

	md5, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", errors.Wrap(err, "Cannot download MD5 file")
	}

	return string(md5), nil
}

func (c *Client) currentHash(eid EditionID) (string, error) {
	md5file := path.Join(c.workDir, fmt.Sprintf(".%s.%s", eid.Filename(), "md5"))
	if _, err := os.Stat(md5file); os.IsNotExist(err) {
		return "", nil
	} else if err != nil {
		return "", err
	}
	curHash, err := ioutil.ReadFile(md5file)
	if err != nil {
		return "", errors.Wrap(err, "Cannot read current archive hash")
	}
	return string(curHash), nil
}

func (c *Client) downloadArchive(eid EditionID, expHash string, archive string) error {
	if _, err := os.Stat(archive); err == nil {
		curHash, err := checksum(archive)
		if err != nil {
			return errors.Wrap(err, "Cannot get archive checksum")
		}
		if expHash == curHash {
			c.log.Debug().
				Str("edition_id", eid.String()).
				Str("hash", expHash).
				Msgf("Archive already downloaded and valid. Skipping download")
			return nil
		}
	}

	c.log.Info().
		Str("edition_id", eid.String()).
		Msgf("Downloading %s archive...", filepath.Base(archive))

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/app/geoip_download", c.baseURL), nil)
	if err != nil {
		return errors.Wrap(err, "Request failed")
	}

	q := req.URL.Query()
	q.Add("license_key", c.licenseKey)
	q.Add("edition_id", eid.String())
	q.Add("suffix", eid.Suffix().String())
	req.URL.RawQuery = q.Encode()

	if c.userAgent != "" {
		req.Header.Add("User-Agent", c.userAgent)
	}

	res, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return errors.Errorf("Received invalid status code %d: %s", res.StatusCode, res.Body)
	}

	out, err := os.Create(archive)
	if err != nil {
		return errors.Wrap(err, "Cannot create archive file")
	}
	defer out.Close()

	_, err = io.Copy(out, res.Body)
	if err != nil {
		return errors.Wrap(err, "Cannot download archive")
	}

	curHash, err := checksum(archive)
	if err != nil {
		return errors.Wrap(err, "Cannot get archive checksum")
	}

	if expHash != curHash {
		return errors.Errorf("MD5 of downloaded archive (%s) does not match expected md5 (%s)", curHash, expHash)
	}

	return nil
}

func (c *Client) extractArchive(eid EditionID, archive string, dlDir string) ([]os.FileInfo, error) {
	var dbs []os.FileInfo
	err := archiver.Walk(archive, func(f archiver.File) error {
		if f.IsDir() {
			return nil
		}
		if filepath.Ext(f.Name()) != ".csv" && filepath.Ext(f.Name()) != ".mmdb" {
			return nil
		}

		expHashRaw := md5.New()
		fileReader := io.TeeReader(f, expHashRaw)
		expHash := hex.EncodeToString(expHashRaw.Sum(nil))

		sublog := log.With().
			Str("edition_id", eid.String()).
			Str("db_name", f.Name()).
			Str("db_size", units.HumanSize(float64(f.Size()))).
			Time("db_modtime", f.ModTime()).
			Str("db_hash", expHash).
			Logger()

		dbpath := path.Join(dlDir, f.Name())
		if fileExists(dbpath) {
			curHash, err := checksum(dbpath)
			if err != nil {
				return err
			}
			if expHash == curHash {
				sublog.Debug().Msg("Database is already up to date")
				return nil
			}
		}

		sublog.Debug().Msg("Extracting database")
		dbfile, err := os.Create(dbpath)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("Cannot create database file %s", f.Name()))
		}
		defer dbfile.Close()

		_, err = io.Copy(dbfile, fileReader)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("Cannot extract database file %s", f.Name()))
		}

		if err = os.Chtimes(dbpath, f.ModTime(), f.ModTime()); err != nil {
			sublog.Warn().Err(err).Msg("Cannot preserve modtime of database file")
		}

		dbs = append(dbs, f)
		return nil
	})

	return dbs, err
}
