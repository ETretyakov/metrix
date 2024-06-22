package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"metrix/pkg/agent/config"
	"metrix/pkg/agent/watcher"

	"github.com/rs/zerolog/log"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Error().Caller().Str("Stage", "loading-config").Err(err).Msg("failed to load config")
	}

	log.Info().Caller().Msg("building watcher")
	w := watcher.NewWatcher(*cfg)

	log.Info().Caller().Msg("starting watcher")
	go w.Start(ctx, cfg.PollInterval)

	log.Info().Caller().Msg("starting to report")
	w.Report(ctx, cfg.ServerURL, cfg.ReportInterval)
}
