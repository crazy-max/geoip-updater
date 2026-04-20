package maxmind

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/go-units"
	"github.com/mholt/archives"
	"github.com/pkg/errors"
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

	if dlDir == "" {
		execPath, err := os.Executable()
		if err != nil {
			return nil, errors.Wrap(err, "Cannot determine executable path")
		}
		dlDir = filepath.Dir(execPath)
	}

	dlDir, err = filepath.Abs(filepath.Clean(dlDir))
	if err != nil {
		return nil, errors.Wrap(err, "Cannot get absolute path of download directory")
	}

	if err := os.MkdirAll(dlDir, 0o755); err != nil {
		return nil, errors.Wrap(err, "Cannot create download directory")
	}

	if err := checkDirWritable(dlDir); err != nil {
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
	expHash, err := d.expectedHash()
	if err != nil {
		return nil, errors.Wrap(err, "Cannot get archive checksum")
	}

	archive := filepath.Join(d.workDir, d.eid.Filename())
	if err := d.downloadArchive(expHash, archive); err != nil {
		return nil, err
	}

	dbs, err := d.extractArchive(archive)
	if err != nil {
		return nil, err
	}

	return dbs, nil
}

func (d *Downloader) expectedHash() (string, error) {
	downloadURL, err := url.JoinPath(d.baseURL, "app", "geoip_download")
	if err != nil {
		return "", errors.Wrap(err, "Cannot create download URL")
	}

	req, err := http.NewRequestWithContext(d.ctx, http.MethodGet, downloadURL, nil)
	if err != nil {
		return "", errors.Wrap(err, "Request failed")
	}

	q := req.URL.Query()
	q.Add("license_key", d.licenseKey)
	q.Add("edition_id", d.eid.String())
	q.Add("suffix", fmt.Sprintf("%s.sha256", d.eid.Suffix().String()))
	req.URL.RawQuery = q.Encode()

	if d.userAgent != "" {
		req.Header.Set("User-Agent", d.userAgent)
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

	checksumFields := strings.Fields(string(checksum))
	if len(checksumFields) == 0 || len(checksumFields[0]) != sha256.Size*2 {
		return "", errors.Errorf("Invalid checksum: %s", checksum)
	}
	if _, err := hex.DecodeString(checksumFields[0]); err != nil {
		return "", errors.Errorf("Invalid checksum: %s", checksum)
	}

	return checksumFields[0], nil
}

func (d *Downloader) downloadArchive(expHash string, archive string) error {
	archivePerm := os.FileMode(0o644)
	if info, err := os.Stat(archive); err == nil {
		archivePerm = info.Mode().Perm()
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

	downloadURL, err := url.JoinPath(d.baseURL, "app", "geoip_download")
	if err != nil {
		return errors.Wrap(err, "Cannot create download URL")
	}

	req, err := http.NewRequestWithContext(d.ctx, http.MethodGet, downloadURL, nil)
	if err != nil {
		return errors.Wrap(err, "Request failed")
	}

	q := req.URL.Query()
	q.Add("license_key", d.licenseKey)
	q.Add("edition_id", d.eid.String())
	q.Add("suffix", d.eid.Suffix().String())
	req.URL.RawQuery = q.Encode()

	if d.userAgent != "" {
		req.Header.Set("User-Agent", d.userAgent)
	}

	res, err := d.http.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return errors.Errorf("Received invalid status code %d", res.StatusCode)
	}

	err = writeFileAtomically(archive, archivePerm, func(tempArchive string, out *os.File) error {
		if _, err := io.Copy(out, res.Body); err != nil {
			return errors.Wrap(err, "Cannot download archive")
		}
		curHash, err := checksumFromFile(tempArchive)
		if err != nil {
			return errors.Wrap(err, "Cannot get archive checksum")
		}
		if expHash != curHash {
			return errors.Errorf("Checksum of downloaded archive (%s) does not match the expected one (%s)", curHash, expHash)
		}
		return nil
	})
	if err != nil {
		return errors.Wrap(err, "Cannot replace archive file")
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

		switch filepath.Ext(de.Name()) {
		case ".csv", ".mmdb":
		default:
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

		sublog := d.log.With().
			Str("edition_id", d.eid.String()).
			Str("db_name", fi.Name()).
			Str("db_size", units.HumanSize(float64(fi.Size()))).
			Time("db_modtime", fi.ModTime()).
			Str("db_hash", expHash).
			Logger()

		dbpath := filepath.Join(d.dlDir, fi.Name())
		dbPerm := fi.Mode().Perm()
		if existingInfo, err := os.Stat(dbpath); err == nil && !existingInfo.IsDir() {
			dbPerm = existingInfo.Mode().Perm()
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
		if err := writeFileAtomically(dbpath, dbPerm, func(tempDBPath string, dbfile *os.File) error {
			if _, err := io.Copy(dbfile, reader); err != nil {
				return errors.Wrapf(err, "Cannot extract database file %s", fi.Name())
			}
			if err := os.Chtimes(tempDBPath, fi.ModTime(), fi.ModTime()); err != nil {
				sublog.Warn().Err(err).Msg("Cannot preserve modtime of database file")
			}
			return nil
		}); err != nil {
			return errors.Wrapf(err, "Cannot replace database file %s", fi.Name())
		}

		dbs = append(dbs, fi)
		return nil
	})

	return dbs, err
}
