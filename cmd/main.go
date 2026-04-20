package main

import (
	"context"
	"os"
	"os/signal"
	"strings"
	"syscall"
	_ "time/tzdata"

	"github.com/alecthomas/kong"
	"github.com/crazy-max/geoip-updater/internal/app"
	"github.com/crazy-max/geoip-updater/internal/config"
	"github.com/crazy-max/geoip-updater/internal/logging"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

var version = "dev"

func main() {
	if err := run(); err != nil {
		log.Fatal().Err(err).Send()
	}
}

func run() error {
	cli := config.Cli{}
	_ = kong.Parse(&cli,
		kong.Name("geoip-updater"),
		kong.Description(`Download and update MaxMind's GeoIP2 databases on a time-based schedule. More info: https://github.com/crazy-max/geoip-updater`),
		kong.UsageOnError(),
		kong.Vars{
			"version": version,
		},
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
			Summary: true,
		}))

	logging.Configure(&cli)
	log.Info().Msgf("Starting geoip-updater %s", version)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Load(cli, version)
	if err != nil {
		return errors.Wrap(err, "cannot load configuration")
	}

	geoipupd, err := app.New(cfg)
	if err != nil {
		return errors.Wrap(err, "cannot initialize geoip-updater")
	}
	if err := geoipupd.Start(ctx); err != nil {
		return errors.Wrap(err, "cannot start geoip-updater")
	}

	if cause := context.Cause(ctx); cause != nil {
		log.Warn().Msg(strings.Title(cause.Error())) //nolint:staticcheck // ignoring "SA1019: strings.Title is deprecated", as for our use we don't need full unicode support
	}

	return nil
}
