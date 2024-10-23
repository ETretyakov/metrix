// Module "app" is the entrypoint for the HTTP serive.
package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"metrix/internal/bootstrap"
	"metrix/internal/closer"
	"metrix/internal/config"
	"metrix/internal/grpcapi/grpcservice"
	"metrix/internal/handlers"
	"metrix/internal/http"
	"metrix/internal/repository"
	"metrix/pkg/logger"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// Run - launches web-server for metrics aggregation.
func Run(ctx context.Context, cfg *config.Config) (err error) {
	ctx, cancel := context.WithCancel(ctx)

	// Database setup
	var db *sqlx.DB
	if cfg.Postgres.DSN != "" {
		db, err = bootstrap.InitDB(ctx, &cfg.Postgres)
		if err != nil {
			logger.Fatal(ctx, "failed to init db", err)
		}
	}

	repoGroup := repository.NewGroup(
		ctx,
		db,
		cfg.FileStoragePath,
		cfg.StoreInterval,
		cfg.Restore,
	)

	// HTTP server
	healthHandlers := handlers.NewHealthHandlers(repoGroup)
	metricsHandlers := handlers.NewMetricsHandlers(repoGroup)

	httpServer := http.New(
		cfg,
		healthHandlers,
		metricsHandlers,
	)

	httpServer.Start(ctx)

	healthHandlers.SetLiveness(true)
	healthHandlers.SetReadiness(true)

	// GRPC Server
	gs := grpcservice.NewGServiceServer(repoGroup.MetricRepo)
	if err := gs.Start(); err != nil {
		cancel()
		return errors.Wrap(err, "failed to start grpc")
	}

	gracefulShutDown(ctx, cancel)

	return nil
}

func gracefulShutDown(ctx context.Context, cancel context.CancelFunc) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(ch)

	sig := <-ch
	errorMessage := fmt.Sprintf("%s %v - %s", "Received shutdown signal:", sig, "Graceful shutdown done")
	logger.Info(ctx, errorMessage)
	cancel()
	closer.CloseAll()
}
