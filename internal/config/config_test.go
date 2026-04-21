package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadPreservesCLIAndSetsApplicationMetadata(t *testing.T) {
	cli := Cli{
		EditionIDs:   []string{"GeoLite2-City", "GeoLite2-ASN"},
		LicenseKey:   "0123456789",
		DownloadPath: "/data",
		Schedule:     "@daily",
		LogLevel:     "debug",
		LogJSON:      true,
		LogCaller:    true,
	}

	cfg, err := Load(cli, "1.2.3")
	require.NoError(t, err)

	require.Equal(t, cli, cfg.Cli)
	require.Equal(t, App{
		Name:    "geoip-updater",
		Desc:    "Download and update MaxMind's GeoIP2 databases on a time-based schedule",
		URL:     "https://github.com/crazy-max/geoip-updater",
		Author:  "CrazyMax",
		Version: "1.2.3",
	}, cfg.App)
}
