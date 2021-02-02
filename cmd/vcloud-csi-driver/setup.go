package main

import (
	"os"
	"strings"

	"github.com/proact-de/vcloud-csi-driver/pkg/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func setupLogger(cfg *config.Config) {
	switch strings.ToLower(cfg.Logs.Level) {
	case "panic":
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	case "fatal":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	if cfg.Logs.Pretty {
		log.Logger = log.Output(
			zerolog.ConsoleWriter{
				Out:     os.Stderr,
				NoColor: !cfg.Logs.Color,
			},
		)
	}
}

func checkEndpointDefined(cfg *config.Config) bool {
	if cfg.Driver.Endpoint == "" {
		log.Error().
			Str("endpoint", cfg.Driver.Endpoint).
			Msg("No endpoint have been defined")

		return false
	}

	return true
}

func ensureSocketRemoved(cfg *config.Config) bool {
	if err := os.Remove(
		cfg.Driver.Endpoint[7:],
	); err != nil && !os.IsNotExist(err) {
		log.Error().
			Err(err).
			Str("endpoint", cfg.Driver.Endpoint).
			Msg("Failed to delete old unix socket")

		return false
	}

	return true
}
