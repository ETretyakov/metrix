package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"metrix/pkg/agent/config"
	"metrix/pkg/agent/watcher"
	"metrix/pkg/client"
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

	useBatching, err := client.CheckBatching(
		ctx,
		"http://"+cfg.Address,
	)
	if err != nil {
		logger.Fatal(ctx, "server is not responding", err)
	}
	cfg.UseBatching = useBatching

	logger.InitDefault(cfg.LogLevel)

	logger.Info(ctx, fmt.Sprintf("starting agent with config: %+v", cfg))

	logger.Info(ctx, "building watcher")
	w := watcher.NewWatcher(*cfg)

	logger.Info(ctx, "starting watcher")
	go w.Start(ctx, time.Second*time.Duration(cfg.PollInterval))

	logger.Info(ctx, "starting to report")
	w.Report(
		ctx,
		"http://"+cfg.Address,
		time.Second*time.Duration(cfg.ReportInterval),
		useBatching,
	)
}
