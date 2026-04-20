package app

import (
	"context"
	"testing"
	"time"

	"github.com/crazy-max/geoip-updater/internal/config"
	"github.com/crazy-max/geoip-updater/pkg/maxmind"
	"github.com/robfig/cron/v3"
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
