package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"metrix/pkg/agent/config"
	"metrix/pkg/agent/monitoring"
	"metrix/pkg/logger"
)

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer cancel()

	cfg, err := config.NewConfig()
	if err != nil {
		logger.Error(ctx, "failed to read config", err)
	}

	logger.InitDefault(cfg.LogLevel)

	watcher := monitoring.NewWatcher(cfg.Metrics)
	watcher.Run(ctx, cfg)
}
