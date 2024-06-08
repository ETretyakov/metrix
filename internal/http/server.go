package http

import (
	"context"
	"fmt"
	"metrix/internal/closer"
	"metrix/internal/config"
	"metrix/internal/handlers"
	"metrix/internal/logger"
	"net/http"
)

type Server struct {
	cfg     *config.Config
	srv     *http.Server
	health  *handlers.HealthHandlers
	metrics *handlers.MetricsHandlers
}

func New(
	cfg *config.Config,
	healthHandlers *handlers.HealthHandlers,
	metricsHandlers *handlers.MetricsHandlers,
) *Server {
	srv := &http.Server{
		Addr: cfg.HTTPAddress,
	}

	return &Server{
		cfg:     cfg,
		srv:     srv,
		health:  healthHandlers,
		metrics: metricsHandlers,
	}
}

func (s *Server) Start(ctx context.Context) {
	s.srv.Handler = s.setupRoutes()

	go func() {
		logger.Info(ctx, fmt.Sprintf("starting listening http srv at %s", s.cfg.HTTPAddress))
		if err := s.srv.ListenAndServe(); err != http.ErrServerClosed {
			logger.Fatal(ctx, "error start http srv, err: %+v", err)
		}
	}()

	closer.Add(s.Close)
}

func (s *Server) Close() error {
	ctx := context.TODO()
	if err := s.srv.Shutdown(ctx); err != nil {
		logger.Error(ctx, "error stop http srv, err", err)
		return err
	}

	logger.Info(ctx, "http server shutdown done")

	return nil
}
