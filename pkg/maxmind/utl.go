package maxmind

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
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
	var b bytes.Buffer

	hash := sha256.New()
	if _, err := io.Copy(&b, io.TeeReader(reader, hash)); err != nil {
		return "", nil, err
	}

	return hex.EncodeToString(hash.Sum(nil)), &b, nil
}

func createFile(path string, content string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err = file.WriteString(content); err != nil {
		return err
	}
	if err = file.Sync(); err != nil {
		return err
	}
	return nil
}

func isDirWriteable(dir string) error {
	f := filepath.Join(dir, ".touch")
	if err := os.WriteFile(f, []byte(""), 0o600); err != nil {
		return err
	}
	return os.Remove(f)
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
