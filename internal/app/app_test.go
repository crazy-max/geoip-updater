package app

import (
	"archive/zip"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/crazy-max/geoip-updater/internal/config"
	"github.com/crazy-max/geoip-updater/pkg/maxmind"
	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
)

func TestNewRejectsInvalidEditionID(t *testing.T) {
	_, err := New(&config.Configuration{
		Cli: config.Cli{
			EditionIDs: []string{"not-a-real-edition"},
		},
	})

	require.EqualError(t, err, "invalid edition ID: not-a-real-edition")
}

func TestNewPreservesEditionIDOrder(t *testing.T) {
	client, err := New(&config.Configuration{
		Cli: config.Cli{
			EditionIDs: []string{
				maxmind.EIDGeoLite2Country.String(),
				maxmind.EIDGeoLite2ASN.String(),
				maxmind.EIDGeoLite2City.String(),
			},
		},
	})
	require.NoError(t, err)
	require.Equal(t, []maxmind.EditionID{
		maxmind.EIDGeoLite2Country,
		maxmind.EIDGeoLite2ASN,
		maxmind.EIDGeoLite2City,
	}, client.eids)
}

func TestStartReturnsErrorWhenMaxMindClientInitFails(t *testing.T) {
	homeDir := t.TempDir()
	t.Setenv("HOME", homeDir)
	t.Setenv("USERPROFILE", homeDir)

	client, err := New(&config.Configuration{
		Cli: config.Cli{
			EditionIDs: []string{maxmind.EIDGeoLite2City.String()},
		},
	})
	require.NoError(t, err)

	err = client.Start(context.Background())
	require.EqualError(t, err, "License key required")
}

func TestStartReturnsImmediatelyWithoutSchedule(t *testing.T) {
	homeDir := t.TempDir()
	t.Setenv("HOME", homeDir)
	t.Setenv("USERPROFILE", homeDir)

	client, err := New(&config.Configuration{
		Cli: config.Cli{
			LicenseKey: "0123456789",
		},
	})
	require.NoError(t, err)

	errCh := make(chan error, 1)
	go func() {
		errCh <- client.Start(context.Background())
	}()

	select {
	case err := <-errCh:
		require.NoError(t, err)
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for Start to return without a schedule")
	}
}

func TestStartReturnsWhenContextCanceled(t *testing.T) {
	homeDir := t.TempDir()
	t.Setenv("HOME", homeDir)
	t.Setenv("USERPROFILE", homeDir)

	ctx, cancel := context.WithCancelCause(context.Background())
	t.Cleanup(func() { cancel(nil) })

	client := &Client{
		cfg: &config.Configuration{
			Cli: config.Cli{
				LicenseKey: "0123456789",
				Schedule:   "@every 1m",
			},
		},
		cron: cron.New(cron.WithParser(cron.NewParser(
			cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor),
		)),
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- client.Start(ctx)
	}()

	require.Never(t, func() bool {
		select {
		case err := <-errCh:
			require.NoError(t, err)
			return true
		default:
			return false
		}
	}, 100*time.Millisecond, 10*time.Millisecond)

	cancel(nil)

	select {
	case err := <-errCh:
		require.NoError(t, err)
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for Start to return")
	}
}

func TestStartReturnsWhenContextAlreadyCanceled(t *testing.T) {
	homeDir := t.TempDir()
	t.Setenv("HOME", homeDir)
	t.Setenv("USERPROFILE", homeDir)

	ctx, cancel := context.WithCancelCause(context.Background())
	cancel(nil)

	client, err := New(&config.Configuration{
		Cli: config.Cli{
			EditionIDs: []string{maxmind.EIDGeoLite2City.String()},
			LicenseKey: "0123456789",
		},
	})
	require.NoError(t, err)

	require.NoError(t, client.Start(ctx))
}

func TestRunSkipsWhenAlreadyRunning(t *testing.T) {
	logBuf := captureGlobalLog(t)

	client := &Client{
		locker: 1,
	}

	client.Run()

	require.Contains(t, logBuf.String(), "Already running")
}

func TestRunContinuesAfterDownloadFailure(t *testing.T) {
	homeDir := t.TempDir()
	t.Setenv("HOME", homeDir)
	t.Setenv("USERPROFILE", homeDir)

	asnEntryName := "GeoLite2-ASN-Blocks-IPv4.csv"
	asnContent := []byte("network,start_ip,end_ip\n1,1.1.1.0,1.1.1.255\n")
	asnArchive := buildZipArchive(t, asnEntryName, asnContent)
	asnChecksum := sha256.Sum256(asnArchive)
	cityChecksum := hex.EncodeToString(bytes.Repeat([]byte("0"), sha256.Size))

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/app/geoip_download" {
			t.Errorf("unexpected request path: %s", req.URL.Path)
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		editionID := req.URL.Query().Get("edition_id")
		suffix := req.URL.Query().Get("suffix")
		switch editionID {
		case maxmind.EIDGeoLite2CityCSV.String():
			if suffix == maxmind.EIDGeoLite2CityCSV.Suffix().String()+".sha256" {
				rw.WriteHeader(http.StatusOK)
				_, _ = rw.Write([]byte(cityChecksum + " " + maxmind.EIDGeoLite2CityCSV.Filename()))
				return
			}
			rw.WriteHeader(http.StatusInternalServerError)
		case maxmind.EIDGeoLite2ASNCSV.String():
			switch suffix {
			case maxmind.EIDGeoLite2ASNCSV.Suffix().String() + ".sha256":
				rw.WriteHeader(http.StatusOK)
				_, _ = rw.Write([]byte(hex.EncodeToString(asnChecksum[:]) + " " + maxmind.EIDGeoLite2ASNCSV.Filename()))
			case maxmind.EIDGeoLite2ASNCSV.Suffix().String():
				rw.WriteHeader(http.StatusOK)
				_, _ = rw.Write(asnArchive)
			default:
				rw.WriteHeader(http.StatusBadRequest)
			}
		default:
			rw.WriteHeader(http.StatusBadRequest)
		}
	}))
	defer server.Close()

	mm, err := maxmind.New(context.Background(), maxmind.Config{
		BaseURL:    server.URL,
		HTTPClient: server.Client(),
		LicenseKey: "0123456789",
		Logger:     zerolog.New(io.Discard),
	})
	require.NoError(t, err)

	dlDir := t.TempDir()
	client := &Client{
		cfg: &config.Configuration{
			Cli: config.Cli{
				DownloadPath: dlDir,
			},
		},
		mm: mm,
		eids: []maxmind.EditionID{
			maxmind.EIDGeoLite2CityCSV,
			maxmind.EIDGeoLite2ASNCSV,
		},
	}

	logBuf := captureGlobalLog(t)
	client.Run()

	require.Contains(t, logBuf.String(), "Cannot download database")

	asnPath := filepath.Join(dlDir, asnEntryName)
	asnData, err := os.ReadFile(asnPath)
	require.NoError(t, err)
	require.Equal(t, asnContent, asnData)

	_, err = os.Stat(filepath.Join(dlDir, "GeoLite2-City-Locations-en.csv"))
	require.ErrorIs(t, err, os.ErrNotExist)
}

func captureGlobalLog(t *testing.T) *bytes.Buffer {
	t.Helper()

	var logBuf bytes.Buffer
	oldLogger := log.Logger
	oldLevel := zerolog.GlobalLevel()

	log.Logger = zerolog.New(&logBuf).Level(zerolog.DebugLevel)
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	t.Cleanup(func() {
		log.Logger = oldLogger
		zerolog.SetGlobalLevel(oldLevel)
	})

	return &logBuf
}

func buildZipArchive(t *testing.T, name string, content []byte) []byte {
	t.Helper()

	var archive bytes.Buffer
	zipWriter := zip.NewWriter(&archive)

	entry, err := zipWriter.Create(name)
	require.NoError(t, err)

	_, err = entry.Write(content)
	require.NoError(t, err)

	require.NoError(t, zipWriter.Close())
	return archive.Bytes()
}
