package maxmind

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChecksumFromFile(t *testing.T) {
	const content = "geoip-updater checksum file content"

	filename := filepath.Join(t.TempDir(), "sample.txt")
	require.NoError(t, os.WriteFile(filename, []byte(content), 0o600))

	got, err := checksumFromFile(filename)
	require.NoError(t, err)

	want := sha256.Sum256([]byte(content))
	assert.Equal(t, hex.EncodeToString(want[:]), got)
}

func TestChecksumFromReader(t *testing.T) {
	const content = "geoip-updater checksum reader content"

	got, replay, err := checksumFromReader(bytes.NewBufferString(content))
	require.NoError(t, err)

	want := sha256.Sum256([]byte(content))
	assert.Equal(t, hex.EncodeToString(want[:]), got)

	replayed, err := io.ReadAll(replay)
	require.NoError(t, err)
	assert.Equal(t, content, string(replayed))
}

func TestWriteFileAtomicallyCreatesFile(t *testing.T) {
	filename := filepath.Join(t.TempDir(), "atomic.txt")
	var seenTempName string

	err := writeFileAtomically(filename, 0o644, func(tempName string, f *os.File) error {
		seenTempName = tempName
		if !assert.NotEqual(t, filename, tempName) {
			return nil
		}
		_, statErr := os.Stat(tempName)
		require.NoError(t, statErr)
		_, err := f.WriteString("atomic content")
		return err
	})
	require.NoError(t, err)

	assert.NotEmpty(t, seenTempName)

	content, err := os.ReadFile(filename)
	require.NoError(t, err)
	assert.Equal(t, "atomic content", string(content))

	_, err = os.Stat(seenTempName)
	require.ErrorIs(t, err, os.ErrNotExist)
}

func TestWriteFileAtomicallyReplacesExistingFile(t *testing.T) {
	filename := filepath.Join(t.TempDir(), "atomic.txt")
	require.NoError(t, os.WriteFile(filename, []byte("old content"), 0o600))

	err := writeFileAtomically(filename, 0o644, func(_ string, f *os.File) error {
		_, err := f.WriteString("new content")
		return err
	})
	require.NoError(t, err)

	content, err := os.ReadFile(filename)
	require.NoError(t, err)
	assert.Equal(t, "new content", string(content))
}

func TestWriteFileAtomicallyCleansUpOnError(t *testing.T) {
	dir := t.TempDir()
	filename := filepath.Join(dir, "atomic.txt")
	require.NoError(t, os.WriteFile(filename, []byte("old content"), 0o600))

	err := writeFileAtomically(filename, 0o644, func(_ string, f *os.File) error {
		if _, err := f.WriteString("partial content"); err != nil {
			return err
		}
		return assert.AnError
	})
	require.ErrorIs(t, err, assert.AnError)

	content, readErr := os.ReadFile(filename)
	require.NoError(t, readErr)
	assert.Equal(t, "old content", string(content))

	tempFiles, globErr := filepath.Glob(filepath.Join(dir, filepath.Base(filename)+".tmp-*"))
	require.NoError(t, globErr)
	assert.Empty(t, tempFiles)
}
