package maxmind

import (
	"context"
	"net/http"
	"net/url"
	"os"
	"path"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

// Client represents an active maxmind object
type Client struct {
	ctx        context.Context
	log        zerolog.Logger
	http       *http.Client
	workDir    string
	licenseKey string
	baseURL    string
	userAgent  string
}

// Config defines the config for maxmind
type Config struct {
	Logger     zerolog.Logger
	LicenseKey string
	BaseURL    string
	UserAgent  string
}

// New returns a maxmind client
func New(config Config) (*Client, error) {
	if config.LicenseKey == "" {
		return nil, errors.New("License key required")
	}

	if config.BaseURL == "" {
		config.BaseURL = "https://download.maxmind.com"
	}
	_, err := url.ParseRequestURI(config.BaseURL)
	if err != nil {
		return nil, errors.Wrap(err, "Invalid base URL")
	}

	workBaseDir, err := os.UserHomeDir()
	if err != nil {
		workBaseDir = os.TempDir()
	}
	workDir := path.Join(workBaseDir, ".geoip-updater")
	if err := os.MkdirAll(workDir, 0o755); err != nil {
		return nil, errors.Wrap(err, "Cannot create work directory")
	}
	if err := isDirWriteable(workDir); err != nil {
		return nil, errors.Wrap(err, "Work directory is not writable")
	}
	config.Logger.Debug().Msgf("Work directory is %s", workDir)

	return &Client{
		ctx:        context.Background(),
		log:        config.Logger,
		http:       http.DefaultClient,
		workDir:    workDir,
		licenseKey: config.LicenseKey,
		baseURL:    config.BaseURL,
		userAgent:  config.UserAgent,
	}, nil
}
