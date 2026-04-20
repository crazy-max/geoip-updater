package logging

import (
	"io"
	"os"
	"testing"

	"github.com/crazy-max/geoip-updater/internal/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
)

func TestConfigureJSONLoggingAddsCallerAndAppliesLevel(t *testing.T) {
	restore := captureStdoutAndGlobals(t)

	Configure(&config.Cli{
		LogJSON:   true,
		LogCaller: true,
		LogLevel:  "debug",
	})

	log.Debug().Msg("json output")

	output := restore()
	require.Contains(t, output, `"level":"debug"`)
	require.Contains(t, output, `"message":"json output"`)
	require.Contains(t, output, `"caller":"`)
}

func TestConfigureConsoleLoggingUsesConsoleWriter(t *testing.T) {
	restore := captureStdoutAndGlobals(t)

	Configure(&config.Cli{
		LogJSON:  false,
		LogLevel: "info",
	})

	log.Debug().Msg("hidden debug")
	log.Info().Msg("console output")

	output := restore()
	require.Contains(t, output, "console output")
	require.NotContains(t, output, "hidden debug")
	require.NotContains(t, output, `"level":"info"`)
	require.NotContains(t, output, `"caller":"`)
}

func captureStdoutAndGlobals(t *testing.T) func() string {
	t.Helper()

	readPipe, writePipe, err := os.Pipe()
	require.NoError(t, err)

	oldStdout := os.Stdout
	oldLogger := log.Logger
	oldLevel := zerolog.GlobalLevel()
	os.Stdout = writePipe

	restored := false
	restore := func() string {
		if !restored {
			restored = true
			_ = writePipe.Close()
			os.Stdout = oldStdout
		}

		data, err := io.ReadAll(readPipe)
		require.NoError(t, err)
		require.NoError(t, readPipe.Close())

		log.Logger = oldLogger
		zerolog.SetGlobalLevel(oldLevel)
		return string(data)
	}

	t.Cleanup(func() {
		if !restored {
			_ = writePipe.Close()
			os.Stdout = oldStdout
		}
		log.Logger = oldLogger
		zerolog.SetGlobalLevel(oldLevel)
	})

	return restore
}
