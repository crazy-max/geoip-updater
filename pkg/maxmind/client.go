package maxmind

import (
	"context"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

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
	LicenseKey string
	BaseURL    string
	UserAgent  string
	HTTPClient *http.Client
	Logger     zerolog.Logger
}

// New returns a maxmind client
func New(ctx context.Context, config Config) (*Client, error) {
	if config.LicenseKey == "" {
		return nil, errors.New("License key required")
	}

	httpClient := config.HTTPClient
	if config.HTTPClient == nil {
		httpClient = http.DefaultClient
	}

	baseURL := config.BaseURL
	if config.BaseURL == "" {
		baseURL = "https://download.maxmind.com"
	}

	_, err := url.ParseRequestURI(baseURL)
	if err != nil {
		return nil, errors.Wrap(err, "Invalid base URL")
	}

	workBaseDir, err := os.UserHomeDir()
	if err != nil {
		workBaseDir = os.TempDir()
	}

	workDir := filepath.Join(workBaseDir, ".geoip-updater")
	if err := os.MkdirAll(workDir, 0o755); err != nil {
		return nil, errors.Wrap(err, "Cannot create work directory")
	}
	if err := checkDirWritable(workDir); err != nil {
		return nil, errors.Wrap(err, "Work directory is not writable")
	}
	config.Logger.Debug().Msgf("Work directory is %s", workDir)

	return &Client{
		ctx:        ctx,
		log:        config.Logger,
		http:       httpClient,
		workDir:    workDir,
		licenseKey: config.LicenseKey,
		baseURL:    baseURL,
		userAgent:  config.UserAgent,
	}, nil
}
