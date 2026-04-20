package maxmind

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

func checksumFromFile(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

func checksumFromReader(reader io.Reader) (string, io.Reader, error) {
	b, err := io.ReadAll(reader)
	if err != nil {
		return "", nil, err
	}

	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:]), bytes.NewReader(b), nil
}

func checkDirWritable(dir string) error {
	f, err := os.CreateTemp(dir, ".geoip-updater-write-*")
	if err != nil {
		return err
	}
	name := f.Name()
	if err := f.Close(); err != nil {
		_ = os.Remove(name)
		return err
	}
	return os.Remove(name)
}

func writeFileAtomically(filename string, perm os.FileMode, write func(tempName string, f *os.File) error) error {
	f, err := os.CreateTemp(filepath.Dir(filename), filepath.Base(filename)+".tmp-*")
	if err != nil {
		return err
	}

	tempName := f.Name()
	committed := false
	defer func() {
		_ = f.Close()
		if !committed {
			_ = os.Remove(tempName)
		}
	}()

	if err := f.Chmod(perm); err != nil {
		return err
	}
	if err := write(tempName, f); err != nil {
		return err
	}
	if err := f.Sync(); err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	if err := os.Remove(filename); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	if err := os.Rename(tempName, filename); err != nil {
		return err
	}

	committed = true
	return nil
}
