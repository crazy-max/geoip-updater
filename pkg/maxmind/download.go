package maxmind

import (
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/docker/go-units"
	"github.com/mholt/archives"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// Downloader represents an active downloader object
type Downloader struct {
	*Client
	eid   EditionID
	dlDir string
}

// NewDownloader returns a new downloader instance
func (c *Client) NewDownloader(eid EditionID, dlDir string) (*Downloader, error) {
	var err error

	// Check download directory
	if dlDir == "" {
		dlDir = filepath.Dir(os.Args[0])
	}
	dlDir, err = filepath.Abs(path.Clean(dlDir))
	if err != nil {
		return nil, errors.Wrap(err, "Cannot get absolute path of download directory")
	}
	if err := os.MkdirAll(dlDir, 0o755); err != nil {
		return nil, errors.Wrap(err, "Cannot create download directory")
	}
	if err := isDirWriteable(dlDir); err != nil {
		return nil, errors.Wrap(err, "Download directory is not writable")
	}

	return &Downloader{
		Client: c,
		eid:    eid,
		dlDir:  dlDir,
	}, nil
}

// Download downloads a database
func (d *Downloader) Download() ([]os.FileInfo, error) {
	// Retrieve expected hash
	expHash, err := d.expectedHash()
	if err != nil {
		return nil, errors.Wrap(err, "Cannot get archive checksum")
	}

	// Download DB archive
	archive := path.Join(d.workDir, d.eid.Filename())
	if err := d.downloadArchive(expHash, archive); err != nil {
		return nil, err
	}

	// Create checksum file
	checksumFile := path.Join(d.workDir, fmt.Sprintf(".%s.%s", d.eid.Filename(), "sha256"))
	if err := createFile(checksumFile, expHash); err != nil {
		return nil, errors.Errorf("Cannot create checksum file %s", checksumFile)
	}

	// Extract DB from archive
	dbs, err := d.extractArchive(archive)
	if err != nil {
		return nil, err
	}

	return dbs, nil
}

func (d *Downloader) expectedHash() (string, error) {
	req, err := http.NewRequestWithContext(d.ctx, "GET", fmt.Sprintf("%s/app/geoip_download", d.baseURL), nil)
	if err != nil {
		return "", errors.Wrap(err, "Request failed")
	}

	q := req.URL.Query()
	q.Add("license_key", d.licenseKey)
	q.Add("edition_id", d.eid.String())
	q.Add("suffix", fmt.Sprintf("%s.sha256", d.eid.Suffix().String()))
	req.URL.RawQuery = q.Encode()

	if d.userAgent != "" {
		req.Header.Add("User-Agent", d.userAgent)
	}

	res, err := d.http.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", errors.Errorf("Received invalid status code %d", res.StatusCode)
	}

	checksum, err := io.ReadAll(res.Body)
	if err != nil {
		return "", errors.Wrap(err, "Cannot download checksum file")
	}

	checksumAr := strings.SplitN(strings.TrimSpace(string(checksum)), " ", 2)
	if len(checksumAr[0]) != 64 {
		return "", errors.Errorf("Invalid checksum: %s", checksum)
	}

	return checksumAr[0], nil
}

func (d *Downloader) downloadArchive(expHash string, archive string) error {
	if _, err := os.Stat(archive); err == nil {
		curHash, err := checksumFromFile(archive)
		if err != nil {
			return errors.Wrap(err, "Cannot get archive checksum")
		}
		if expHash == curHash {
			d.log.Debug().
				Str("edition_id", d.eid.String()).
				Str("hash", expHash).
				Msgf("Archive already downloaded and valid. Skipping download")
			return nil
		}
	}

	d.log.Info().
		Str("edition_id", d.eid.String()).
		Msgf("Downloading %s archive...", filepath.Base(archive))

	req, err := http.NewRequestWithContext(d.ctx, "GET", fmt.Sprintf("%s/app/geoip_download", d.baseURL), nil)
	if err != nil {
		return errors.Wrap(err, "Request failed")
	}

	q := req.URL.Query()
	q.Add("license_key", d.licenseKey)
	q.Add("edition_id", d.eid.String())
	q.Add("suffix", d.eid.Suffix().String())
	req.URL.RawQuery = q.Encode()

	if d.userAgent != "" {
		req.Header.Add("User-Agent", d.userAgent)
	}

	res, err := d.http.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return errors.Errorf("Received invalid status code %d", res.StatusCode)
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

	curHash, err := checksumFromFile(archive)
	if err != nil {
		return errors.Wrap(err, "Cannot get archive checksum")
	}

	if expHash != curHash {
		return errors.Errorf("Checksum of downloaded archive (%s) does not match the expected one (%s)", curHash, expHash)
	}

	return nil
}

func (d *Downloader) extractArchive(archive string) ([]os.FileInfo, error) {
	var dbs []os.FileInfo

	fsys, err := archives.FileSystem(d.ctx, archive, nil)
	if err != nil {
		return nil, err
	}

	err = fs.WalkDir(fsys, ".", func(path string, de fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if de.IsDir() {
			return nil
		}

		f, err := fsys.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		if filepath.Ext(de.Name()) != ".csv" && filepath.Ext(de.Name()) != ".mmdb" {
			return nil
		}

		expHash, reader, err := checksumFromReader(f)
		if err != nil {
			return err
		}

		fi, err := de.Info()
		if err != nil {
			return err
		}

		sublog := log.With().
			Str("edition_id", d.eid.String()).
			Str("db_name", fi.Name()).
			Str("db_size", units.HumanSize(float64(fi.Size()))).
			Time("db_modtime", fi.ModTime()).
			Str("db_hash", expHash).
			Logger()

		dbpath := filepath.Join(d.dlDir, fi.Name())
		if fileExists(dbpath) {
			curHash, err := checksumFromFile(dbpath)
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
			return errors.Wrap(err, fmt.Sprintf("Cannot create database file %s", fi.Name()))
		}
		defer dbfile.Close()

		_, err = io.Copy(dbfile, reader)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("Cannot extract database file %s", fi.Name()))
		}

		if err = os.Chtimes(dbpath, fi.ModTime(), fi.ModTime()); err != nil {
			sublog.Warn().Err(err).Msg("Cannot preserve modtime of database file")
		}

		dbs = append(dbs, fi)
		return nil
	})

	return dbs, err
}
