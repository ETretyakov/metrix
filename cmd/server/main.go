package main

import (
	"context"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	"metrix/internal/app"
	"metrix/internal/config"
	"metrix/internal/middlewares"
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

	middlewares.SetSignKey(cfg.SignKey)

	if err := app.Run(ctx, cfg); err != nil {
		logger.Error(ctx, "error running http server", err)
	}
}
