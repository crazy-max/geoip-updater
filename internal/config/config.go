package config

// Configuration holds configuration details
type Configuration struct {
	Flags Flags
	App   App
}

// Flags holds flags from command line
type Flags struct {
	LicenseKey   string
	EditionIDs   string
	DownloadPath string
	Schedule     string
	Timezone     string
	LogLevel     string
	LogJson      bool
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
func Load(fl Flags, version string) (*Configuration, error) {
	cfg := &Configuration{
		Flags: fl,
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
