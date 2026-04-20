package maxmind

import (
	"archive/zip"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExtractArchiveUsesClientLogger(t *testing.T) {
	homeDir := t.TempDir()
	t.Setenv("HOME", homeDir)
	t.Setenv("USERPROFILE", homeDir)

	archivePath := filepath.Join(t.TempDir(), "archive.zip")
	dbName := "GeoLite2-City.mmdb"
	dbContent := []byte("existing database content")
	require.NoError(t, writeZipArchive(archivePath, dbName, dbContent))

	dlDir := t.TempDir()
	dbPath := filepath.Join(dlDir, dbName)
	require.NoError(t, os.WriteFile(dbPath, dbContent, 0o600))

	var logBuf bytes.Buffer
	logger := zerolog.New(&logBuf).Level(zerolog.DebugLevel)
	client, err := New(context.Background(), Config{
		Logger:     logger,
		LicenseKey: "0123456789",
	})
	require.NoError(t, err)

	downloader, err := client.NewDownloader(EIDGeoLite2City, dlDir)
	require.NoError(t, err)

	dbs, err := downloader.extractArchive(archivePath)
	require.NoError(t, err)

	assert.Empty(t, dbs)
	assert.Contains(t, logBuf.String(), "Database is already up to date")
}

func TestDownloadDoesNotWriteChecksumSidecar(t *testing.T) {
	homeDir := t.TempDir()
	t.Setenv("HOME", homeDir)
	t.Setenv("USERPROFILE", homeDir)

	eid := EIDGeoLite2ASNCSV
	dbName := "GeoLite2-ASN-Blocks-IPv4.csv"
	dbContent := []byte("network,start_ip,end_ip\n1,1.1.1.0,1.1.1.255\n")

	var archive bytes.Buffer
	zipWriter := zip.NewWriter(&archive)
	entry, err := zipWriter.Create(dbName)
	require.NoError(t, err)
	_, err = entry.Write(dbContent)
	require.NoError(t, err)
	require.NoError(t, zipWriter.Close())

	archiveChecksum := sha256.Sum256(archive.Bytes())

	srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, "/app/geoip_download", req.URL.Path)
		assert.Equal(t, "0123456789", req.URL.Query().Get("license_key"))
		assert.Equal(t, eid.String(), req.URL.Query().Get("edition_id"))

		switch req.URL.Query().Get("suffix") {
		case fmt.Sprintf("%s.sha256", eid.Suffix().String()):
			rw.WriteHeader(http.StatusOK)
			_, _ = rw.Write([]byte(hex.EncodeToString(archiveChecksum[:]) + " " + eid.Filename()))
		case eid.Suffix().String():
			rw.WriteHeader(http.StatusOK)
			_, _ = rw.Write(archive.Bytes())
		default:
			rw.WriteHeader(http.StatusBadRequest)
		}
	}))
	defer srv.Close()

	client, err := New(context.Background(), Config{
		HTTPClient: srv.Client(),
		LicenseKey: "0123456789",
		BaseURL:    srv.URL,
	})
	require.NoError(t, err)

	dlDir := t.TempDir()
	downloader, err := client.NewDownloader(eid, dlDir)
	require.NoError(t, err)

	dbs, err := downloader.Download()
	require.NoError(t, err)

	if assert.Len(t, dbs, 1) {
		assert.Equal(t, dbName, dbs[0].Name())
	}

	extractedContent, err := os.ReadFile(filepath.Join(dlDir, dbName))
	require.NoError(t, err)
	assert.Equal(t, string(dbContent), string(extractedContent))

	checksumSidecar := filepath.Join(client.workDir, fmt.Sprintf(".%s.sha256", eid.Filename()))
	_, err = os.Stat(checksumSidecar)
	require.ErrorIs(t, err, os.ErrNotExist)
}

func TestExpectedHashRejectsMalformedHexChecksum(t *testing.T) {
	client := newTestClient(t, Config{
		HTTPClient: newStaticResponseClient(http.StatusOK, stringsRepeat("z", sha256.Size*2)),
		LicenseKey: "0123456789",
		BaseURL:    "https://example.com",
	})

	downloader, err := client.NewDownloader(EIDGeoLite2ASNCSV, t.TempDir())
	require.NoError(t, err)

	_, err = downloader.expectedHash()
	assert.EqualError(t, err, "Invalid checksum: "+stringsRepeat("z", sha256.Size*2))
}

func TestDownloadArchiveSkipsNetworkWhenCachedArchiveMatches(t *testing.T) {
	var calls atomic.Int32

	client := newTestClient(t, Config{
		HTTPClient: &http.Client{
			Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
				calls.Add(1)
				return nil, assert.AnError
			}),
		},
		LicenseKey: "0123456789",
	})

	archive := filepath.Join(client.workDir, EIDGeoLite2ASNCSV.Filename())
	content := []byte("already cached archive")
	require.NoError(t, os.WriteFile(archive, content, 0o600))

	expHash := sha256.Sum256(content)
	downloader, err := client.NewDownloader(EIDGeoLite2ASNCSV, t.TempDir())
	require.NoError(t, err)

	err = downloader.downloadArchive(hex.EncodeToString(expHash[:]), archive)
	require.NoError(t, err)

	assert.Zero(t, calls.Load())
}

func TestDownloadArchiveChecksumMismatchLeavesNoArchive(t *testing.T) {
	client := newTestClient(t, Config{
		HTTPClient: newStaticResponseClient(http.StatusOK, "wrong archive bytes"),
		LicenseKey: "0123456789",
		BaseURL:    "https://example.com",
	})

	archive := filepath.Join(client.workDir, EIDGeoLite2ASNCSV.Filename())
	downloader, err := client.NewDownloader(EIDGeoLite2ASNCSV, t.TempDir())
	require.NoError(t, err)

	err = downloader.downloadArchive(stringsRepeat("0", sha256.Size*2), archive)
	require.Error(t, err)

	_, statErr := os.Stat(archive)
	require.ErrorIs(t, statErr, os.ErrNotExist)
}

func TestExtractArchiveSkipsNonDatabaseFilesAndPreservesModTime(t *testing.T) {
	archivePath := filepath.Join(t.TempDir(), "archive.zip")
	modTime := time.Date(2026, time.April, 20, 12, 34, 56, 0, time.UTC).Truncate(2 * time.Second)
	dbName := "GeoLite2-City.mmdb"
	dbContent := []byte("db-content")

	err := writeZipArchiveWithEntries(archivePath, []zipEntry{
		{name: "README.txt", content: []byte("ignore me"), modTime: modTime},
		{name: dbName, content: dbContent, modTime: modTime},
	})
	require.NoError(t, err)

	client := newTestClient(t, Config{LicenseKey: "0123456789"})
	dlDir := t.TempDir()
	downloader, err := client.NewDownloader(EIDGeoLite2City, dlDir)
	require.NoError(t, err)

	dbs, err := downloader.extractArchive(archivePath)
	require.NoError(t, err)

	if assert.Len(t, dbs, 1) {
		assert.Equal(t, dbName, dbs[0].Name())
		assert.Equal(t, modTime, dbs[0].ModTime())
	}

	_, err = os.Stat(filepath.Join(dlDir, "README.txt"))
	require.ErrorIs(t, err, os.ErrNotExist)

	info, err := os.Stat(filepath.Join(dlDir, dbName))
	require.NoError(t, err)
	assert.Equal(t, modTime, info.ModTime().UTC())

	content, err := os.ReadFile(filepath.Join(dlDir, dbName))
	require.NoError(t, err)
	assert.Equal(t, dbContent, content)
}

func TestDownloadReturnsEmptyWhenArchiveAndDatabaseAreCurrent(t *testing.T) {
	homeDir := t.TempDir()
	t.Setenv("HOME", homeDir)
	t.Setenv("USERPROFILE", homeDir)

	eid := EIDGeoLite2ASNCSV
	dbName := "GeoLite2-ASN-Blocks-IPv4.csv"
	dbContent := []byte("network,start_ip,end_ip\n1,1.1.1.0,1.1.1.255\n")

	var archive bytes.Buffer
	zipWriter := zip.NewWriter(&archive)
	entry, err := zipWriter.Create(dbName)
	require.NoError(t, err)
	_, err = entry.Write(dbContent)
	require.NoError(t, err)
	require.NoError(t, zipWriter.Close())

	archiveChecksum := sha256.Sum256(archive.Bytes())
	var archiveRequests atomic.Int32

	srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		switch req.URL.Query().Get("suffix") {
		case fmt.Sprintf("%s.sha256", eid.Suffix().String()):
			rw.WriteHeader(http.StatusOK)
			_, _ = rw.Write([]byte(hex.EncodeToString(archiveChecksum[:]) + " " + eid.Filename()))
		case eid.Suffix().String():
			archiveRequests.Add(1)
			rw.WriteHeader(http.StatusOK)
			_, _ = rw.Write(archive.Bytes())
		default:
			rw.WriteHeader(http.StatusBadRequest)
		}
	}))
	defer srv.Close()

	client := newTestClient(t, Config{
		HTTPClient: srv.Client(),
		LicenseKey: "0123456789",
		BaseURL:    srv.URL,
	})

	dlDir := t.TempDir()
	downloader, err := client.NewDownloader(eid, dlDir)
	require.NoError(t, err)

	firstRun, err := downloader.Download()
	require.NoError(t, err)
	assert.Len(t, firstRun, 1)

	secondRun, err := downloader.Download()
	require.NoError(t, err)
	assert.Empty(t, secondRun)
	assert.Equal(t, int32(1), archiveRequests.Load())
}

func writeZipArchive(filename string, entryName string, content []byte) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	zipWriter := zip.NewWriter(file)
	entry, err := zipWriter.Create(entryName)
	if err != nil {
		_ = zipWriter.Close()
		return err
	}
	if _, err := entry.Write(content); err != nil {
		_ = zipWriter.Close()
		return err
	}
	return zipWriter.Close()
}

type zipEntry struct {
	name    string
	content []byte
	modTime time.Time
}

func writeZipArchiveWithEntries(filename string, entries []zipEntry) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	zipWriter := zip.NewWriter(file)
	for _, entry := range entries {
		header := &zip.FileHeader{
			Name:     entry.name,
			Method:   zip.Deflate,
			Modified: entry.modTime,
		}
		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			_ = zipWriter.Close()
			return err
		}
		if _, err := writer.Write(entry.content); err != nil {
			_ = zipWriter.Close()
			return err
		}
	}

	return zipWriter.Close()
}

func newTestClient(t *testing.T, cfg Config) *Client {
	t.Helper()

	homeDir := t.TempDir()
	t.Setenv("HOME", homeDir)
	t.Setenv("USERPROFILE", homeDir)

	client, err := New(context.Background(), cfg)
	require.NoError(t, err)
	return client
}

func newStaticResponseClient(statusCode int, body string) *http.Client {
	return &http.Client{
		Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: statusCode,
				Body:       io.NopCloser(bytes.NewBufferString(body)),
				Header:     make(http.Header),
			}, nil
		}),
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (fn roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}

func stringsRepeat(s string, count int) string {
	var b bytes.Buffer
	for range count {
		b.WriteString(s)
	}
	return b.String()
}
