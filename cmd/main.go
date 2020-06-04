package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alecthomas/kong"
	"github.com/crazy-max/geoip-updater/internal/app"
	"github.com/crazy-max/geoip-updater/internal/config"
	"github.com/crazy-max/geoip-updater/internal/logging"
	"github.com/rs/zerolog/log"
)

var (
	geoipupd *app.Client
	cli      config.Cli
	version  = "dev"
)

func main() {
	// Parse command line
	_ = kong.Parse(&cli,
		kong.Name("geoip-updater"),
		kong.Description(`Download and update MaxMind's GeoIP2 databases on a time-based schedule. More info: https://github.com/crazy-max/geoip-updater`),
		kong.UsageOnError(),
		kong.Vars{
			"version": fmt.Sprintf("%s", version),
		},
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
			Summary: true,
		}))

	// Load timezone location
	location, err := time.LoadLocation(cli.Timezone)
	if err != nil {
		log.Panic().Err(err).Msgf("Cannot load timezone %s", cli.Timezone)
	}

	// Init
	logging.Configure(&cli, location)
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
	cfg, err := config.Load(cli, version)
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
