package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"metrix/pkg/agent/config"
	"metrix/pkg/agent/watcher"
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

	logger.Info(ctx, "building watcher")
	w := watcher.NewWatcher(*cfg)

	logger.Info(ctx, "starting watcher")
	go w.Start(ctx, time.Second*time.Duration(cfg.PollInterval))

	logger.Info(ctx, "starting to report")
	w.Report(
		ctx,
		"http://"+cfg.Address,
		time.Second*time.Duration(cfg.ReportInterval),
	)
}
