package app

import (
	"fmt"
	"runtime"
	"strings"
	"sync/atomic"
	"time"

	"github.com/crazy-max/geoip-updater/internal/config"
	"github.com/crazy-max/geoip-updater/pkg/maxmind"
	"github.com/docker/go-units"
	"github.com/hako/durafmt"
	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog/log"
)

// Client represents an active geoip-updater object
type Client struct {
	cfg    *config.Configuration
	loc    *time.Location
	cron   *cron.Cron
	mm     *maxmind.Client
	eids   []maxmind.EditionID
	jobID  cron.EntryID
	locker uint32
}

// New creates new geoip-updater instance
func New(cfg *config.Configuration, loc *time.Location) (*Client, error) {
	// Check edition IDs
	var eids []maxmind.EditionID
	eidList := strings.Split(cfg.Flags.EditionIDs, ",")
	for _, eidStr := range eidList {
		eid, err := maxmind.GetEditionID(eidStr)
		if err != nil {
			return nil, err
		}
		eids = append(eids, eid)
	}

	// MaxMind client
	mmcli, err := maxmind.New(maxmind.Config{
		Logger:     log.Logger,
		LicenseKey: cfg.Flags.LicenseKey,
		UserAgent:  fmt.Sprintf("geoip-updater/%s go/%s %s", cfg.App.Version, runtime.Version()[2:], strings.Title(runtime.GOOS)),
	})
	if err != nil {
		return nil, err
	}

	return &Client{
		cfg: cfg,
		loc: loc,
		cron: cron.New(cron.WithLocation(loc), cron.WithParser(cron.NewParser(
			cron.SecondOptional|cron.Minute|cron.Hour|cron.Dom|cron.Month|cron.Dow|cron.Descriptor),
		)),
		mm:   mmcli,
		eids: eids,
	}, nil
}

// Start starts geoip-updater
func (c *Client) Start() error {
	var err error

	// Run on startup
	c.Run()

	// Check scheduler enabled
	if c.cfg.Flags.Schedule == "" {
		return nil
	}

	// Init scheduler
	c.jobID, err = c.cron.AddJob(c.cfg.Flags.Schedule, c)
	if err != nil {
		return err
	}
	log.Info().Msgf("Cron initialized with schedule %s", c.cfg.Flags.Schedule)

	// Start scheduler
	c.cron.Start()
	log.Info().Msgf("Next run in %s (%s)",
		durafmt.ParseShort(c.cron.Entry(c.jobID).Next.Sub(time.Now())).String(),
		c.cron.Entry(c.jobID).Next)

	select {}
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
			durafmt.ParseShort(c.cron.Entry(c.jobID).Next.Sub(time.Now())).String(),
			c.cron.Entry(c.jobID).Next)
	}

	for _, eid := range c.eids {
		sublog := log.With().Str("edition_id", eid.String()).Logger()

		dbs, err := c.mm.DownloadDB(eid, c.cfg.Flags.DownloadPath)
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

// Close closes geoip-updater
func (c *Client) Close() {
	if c.cron != nil {
		c.cron.Stop()
	}
}
