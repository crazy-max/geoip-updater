package config

import "github.com/alecthomas/kong"

// Configuration holds configuration details
type Configuration struct {
	Cli Cli
	App App
}

// Cli holds command line args, flags and cmds
type Cli struct {
	Version      kong.VersionFlag
	EditionIDs   []string `kong:"required,name='edition-ids',env='EDITION_IDS',sep=',',placeholder='GeoLite2-City,GeoLite2-Country',help='MaxMind Edition ID dbs to download (comma separated).'"`
	LicenseKey   string   `kong:"required,name='license-key',env='LICENSE_KEY',placeholder='0123456789',help='MaxMind License Key'"`
	DownloadPath string   `kong:"name='download-path',env='DOWNLOAD_PATH',placeholder='./',help='Directory where databases will be stored.'"`
	Schedule     string   `kong:"name='schedule',env='SCHEDULE',placeholder='0 0 * * 0',help='CRON expression format.'"`
	Timezone     string   `kong:"name='timezone',env='TZ',default='UTC',help='Timezone assigned to geoip-updater.'"`
	LogLevel     string   `kong:"name='log-level',env='LOG_LEVEL',default='info',help='Set log level.'"`
	LogJSON      bool     `kong:"name='log-json',env='LOG_JSON',default='false',help='Enable JSON logging output.'"`
	LogCaller    bool     `kong:"name='log-caller',env='LOG_CALLER',default='false',help='Add file:line of the caller to log output.'"`
}

// App holds application details
type App struct {
	Name    string
	Desc    string
	URL     string
	Author  string
	Version string
}

// Load returns Configuration struct
func Load(cli Cli, version string) (*Configuration, error) {
	cfg := &Configuration{
		Cli: cli,
		App: App{
			Name:    "geoip-updater",
			Desc:    "Download MaxMind's GeoIP2 databases on a time-based schedule",
			URL:     "https://github.com/crazy-max/geoip-updater",
			Author:  "CrazyMax",
			Version: version,
		},
	}
	return cfg, nil
}
