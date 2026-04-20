package app

import (
	"context"
	"fmt"
	"runtime"
	"strings"
	"sync/atomic"

	"github.com/crazy-max/geoip-updater/internal/config"
	"github.com/crazy-max/geoip-updater/pkg/maxmind"
	"github.com/docker/go-units"
	"github.com/dromara/carbon/v2"
	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog/log"
)

// Client represents an active geoip-updater object
type Client struct {
	cfg    *config.Configuration
	cron   *cron.Cron
	mm     *maxmind.Client
	eids   []maxmind.EditionID
	jobID  cron.EntryID
	locker uint32
}

// New creates new geoip-updater instance
func New(cfg *config.Configuration) (*Client, error) {
	var eids []maxmind.EditionID
	for _, eidStr := range cfg.Cli.EditionIDs {
		eid, err := maxmind.GetEditionID(eidStr)
		if err != nil {
			return nil, err
		}
		eids = append(eids, eid)
	}

	return &Client{
		cfg: cfg,
		cron: cron.New(cron.WithParser(cron.NewParser(
			cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor),
		)),
		eids: eids,
	}, nil
}

// Start starts geoip-updater
func (c *Client) Start(ctx context.Context) error {
	var err error
	c.mm, err = maxmind.New(ctx, maxmind.Config{
		Logger:     log.Logger,
		LicenseKey: c.cfg.Cli.LicenseKey,
		UserAgent:  fmt.Sprintf("geoip-updater/%s go/%s %s", c.cfg.App.Version, runtime.Version()[2:], strings.Title(runtime.GOOS)), //nolint:staticcheck // ignoring "SA1019: strings.Title is deprecated", as for our use we don't need full unicode support
	})
	if err != nil {
		return err
	}

	c.Run()

	if c.cfg.Cli.Schedule == "" {
		return nil
	}

	c.jobID, err = c.cron.AddJob(c.cfg.Cli.Schedule, c)
	if err != nil {
		return err
	}
	log.Info().Msgf("Cron initialized with schedule %s", c.cfg.Cli.Schedule)

	c.cron.Start()
	log.Info().Msgf("Next run in %s (%s)",
		carbon.CreateFromStdTime(c.cron.Entry(c.jobID).Next).DiffAbsInString(),
		c.cron.Entry(c.jobID).Next)

	<-ctx.Done()
	<-c.cron.Stop().Done()
	return nil
}

// Run runs geoip-updater process
func (c *Client) Run() {
	if !atomic.CompareAndSwapUint32(&c.locker, 0, 1) {
		log.Warn().Msg("Already running")
		return
	}
	defer atomic.StoreUint32(&c.locker, 0)

	if c.jobID > 0 {
		defer log.Info().Msgf("Next run in %s (%s)",
			carbon.CreateFromStdTime(c.cron.Entry(c.jobID).Next).DiffAbsInString(),
			c.cron.Entry(c.jobID).Next)
	}

	for _, eid := range c.eids {
		sublog := log.With().Str("edition_id", eid.String()).Logger()

		dcli, err := c.mm.NewDownloader(eid, c.cfg.Cli.DownloadPath)
		if err != nil {
			sublog.Error().Err(err).Msg("Cannot create downloader instance")
			continue
		}

		dbs, err := dcli.Download()
		if err != nil {
			sublog.Error().Err(err).Msg("Cannot download database")
			continue
		}

		if len(dbs) == 0 {
			sublog.Info().Msg("This edition is already up to date")
			continue
		}

		for _, db := range dbs {
			sublog.Info().
				Str("size", units.HumanSize(float64(db.Size()))).
				Time("modtime", db.ModTime()).
				Msgf("%s database successfully updated", db.Name())
		}
	}
}
