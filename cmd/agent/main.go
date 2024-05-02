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

	config, err := config.LoadConfig()
	if err != nil {
		log.Error().Caller().Str("Stage", "loading-config").Err(err).Msg("failed to load config")
	}

	log.Info().Caller().Msg("building watcher")
	watcher := watcher.NewWatcher(*config)

	log.Info().Caller().Msg("starting watcher")
	go watcher.Start(ctx, config.PollInterval)

	log.Info().Caller().Msg("starting to report")
	watcher.Report(ctx, config.ServerUrl, config.ReportInterval)
}
