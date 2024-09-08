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

var (
	buildVersion = "N/A"
	buildData    = "N/A"
	buildCommit  = "N/A"
)

// @Title MetrixAPI
// @Description The backend service for metrics aggregation
// @Version 1.0.0
// @Contact.email etretyakov@kaf65.ru
// @BasePath api/v1
// @Host localhost:8080.

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

	logger.Info(ctx, "Build version: "+buildVersion)
	logger.Info(ctx, "Build data: "+buildData)
	logger.Info(ctx, "Build commit: "+buildCommit)

	middlewares.SetSignKey(cfg.SignKey)

	if err := app.Run(ctx, cfg); err != nil {
		logger.Error(ctx, "error running http server", err)
	}
}
