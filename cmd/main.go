package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alecthomas/kingpin"
	"github.com/crazy-max/geoip-updater/internal/app"
	"github.com/crazy-max/geoip-updater/internal/config"
	"github.com/crazy-max/geoip-updater/internal/logging"
	"github.com/crazy-max/geoip-updater/pkg/maxmind"
	"github.com/rs/zerolog/log"
)

var (
	geoipupd *app.Client
	flags    config.Flags
	version  = "dev"
)

func main() {
	// Parse command line
	kingpin.Arg("edition-ids", "MaxMind Edition ID dbs to download (comma separated).").
		Required().
		Envar("EDITION_IDS").
		HintOptions(string(maxmind.EIDGeoLite2ASN), fmt.Sprintf("%s,%s", maxmind.EIDGeoLite2ASN, maxmind.EIDGeoLite2City)).
		StringVar(&flags.EditionIDs)
	kingpin.Flag("license-key", "MaxMind License Key.").
		Required().
		Envar("LICENSE_KEY").
		PlaceHolder("0123456789").
		StringVar(&flags.LicenseKey)
	kingpin.Flag("download-path", "Folder where databases will be stored.").
		Envar("DOWNLOAD_PATH").
		PlaceHolder("./").
		StringVar(&flags.DownloadPath)
	kingpin.Flag("schedule", "CRON expression format.").
		Envar("SCHEDULE").
		PlaceHolder("0 0 * * 0").
		StringVar(&flags.Schedule)
	kingpin.Flag("timezone", "Timezone assigned to geoip-updater.").
		Envar("TZ").
		Default("UTC").
		HintOptions("Europe/Paris").
		StringVar(&flags.Timezone)
	kingpin.Flag("log-level", "Set log level.").
		Envar("LOG_LEVEL").
		Default("info").
		HintOptions("info", "warn", "debug").
		StringVar(&flags.LogLevel)
	kingpin.Flag("log-json", "Enable JSON logging output.").
		Envar("LOG_JSON").
		Default("false").
		BoolVar(&flags.LogJson)
	kingpin.UsageTemplate(kingpin.CompactUsageTemplate).Version(version).Author("CrazyMax")
	kingpin.CommandLine.Name = "geoip-updater"
	kingpin.CommandLine.Help = `Download MaxMind's GeoIP2 databases on a time-based schedule. More info: https://github.com/crazy-max/geoip-updater`
	kingpin.Parse()

	// Load timezone location
	location, err := time.LoadLocation(flags.Timezone)
	if err != nil {
		log.Panic().Err(err).Msgf("Cannot load timezone %s", flags.Timezone)
	}

	// Init
	logging.Configure(&flags, location)
	log.Info().Msgf("Starting geoip-updater %s", version)

	// Handle os signals
	channel := make(chan os.Signal)
	signal.Notify(channel, os.Interrupt, syscall.SIGTERM)
	go func() {
		sig := <-channel
		geoipupd.Close()
		log.Warn().Msgf("Caught signal %v", sig)
		os.Exit(0)
	}()

	// Load and check configuration
	cfg, err := config.Load(flags, version)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot load configuration")
	}

	// Init
	if geoipupd, err = app.New(cfg, location); err != nil {
		log.Fatal().Err(err).Msg("Cannot initialize geoip-updater")
	}

	// Start
	if err = geoipupd.Start(); err != nil {
		log.Fatal().Err(err).Msg("Cannot start geoip-updater")
	}
}
