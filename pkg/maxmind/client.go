package maxmind

import (
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

// Client represents an active maxmind object
type Client struct {
	log        zerolog.Logger
	http       *http.Client
	tmpdir     string
	licenseKey string
	editionIDs []EditionID
	dlPath     string
	baseURL    string
	userAgent  string
}

// Config defines the config for maxmind
type Config struct {
	Logger       zerolog.Logger
	LicenseKey   string
	EditionIDs   []EditionID
	DownloadPath string
	BaseURL      string
	UserAgent    string
}

// New returns a mmdbdl client
func New(config Config) (*Client, error) {
	var err error
	if config.LicenseKey == "" {
		return nil, errors.New("License key required")
	}

	if len(config.EditionIDs) == 0 {
		return nil, errors.New("At least one edition ID is required")
	}

	if config.DownloadPath == "" {
		config.DownloadPath = filepath.Dir(os.Args[0])
	}
	config.DownloadPath, err = filepath.Abs(path.Clean(config.DownloadPath))
	if err != nil {
		return nil, errors.Wrap(err, "Cannot get absolute download path")
	}
	if err := os.MkdirAll(config.DownloadPath, 0755); err != nil {
		return nil, errors.Wrap(err, "Cannot create download folder")
	}
	if err := isDirWriteable(config.DownloadPath); err != nil {
		return nil, errors.Wrap(err, "Download folder is not writable")
	}
	config.DownloadPath = formatPath(config.DownloadPath)
	config.Logger.Debug().Msgf("Download path: %s", config.DownloadPath)

	if config.BaseURL == "" {
		config.BaseURL = "https://download.maxmind.com"
	}
	_, err = url.ParseRequestURI(config.BaseURL)
	if err != nil {
		return nil, errors.Wrap(err, "Invalid base URL")
	}

	tmpdir := formatPath(path.Join(os.TempDir(), "geoip-updater"))
	config.Logger.Debug().Msgf("Temp path: %s", tmpdir)
	if err := os.MkdirAll(tmpdir, 0755); err != nil {
		return nil, errors.Wrap(err, "Cannot create temp folder")
	}
	if err := isDirWriteable(config.DownloadPath); err != nil {
		return nil, errors.Wrap(err, "Temp folder is not writable")
	}

	return &Client{
		log:        config.Logger,
		http:       http.DefaultClient,
		tmpdir:     tmpdir,
		licenseKey: config.LicenseKey,
		editionIDs: config.EditionIDs,
		dlPath:     config.DownloadPath,
		baseURL:    config.BaseURL,
		userAgent:  config.UserAgent,
	}, nil
}
