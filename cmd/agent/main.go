package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"metrix/pkg/agent/config"
	"metrix/pkg/agent/monitoring"
	"metrix/pkg/crypto"
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
		return
	}

	logger.InitDefault(cfg.LogLevel)

	encryption, err := crypto.NewEncryption(cfg.CryptoKey)
	if err != nil {
		logger.Error(ctx, "failed to init encryption", err)
		return
	}

	watcher := monitoring.NewWatcher(cfg.Metrics, encryption)
	watcher.Run(ctx, cfg)
}
