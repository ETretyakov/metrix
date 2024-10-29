// Module "http" aggregates functionality to launch and manage HTTP-server.
package http

import (
	"context"
	"net/http"

	"metrix/internal/closer"
	"metrix/internal/config"
	"metrix/internal/handlers"
	"metrix/pkg/logger"

	"github.com/pkg/errors"
)

// Server - the structure that holds HTTP-server config and related handlers.
type Server struct {
	cfg     *config.Config
	srv     *http.Server
	health  *handlers.HealthHandlers
	metrics *handlers.MetricsHandlers
}

// New - the builder function for server entity.
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

// Start - the method to start the http-server.
func (s *Server) Start(ctx context.Context) {
	s.srv.Handler = s.setupRoutes()

	go func() {
		logger.Info(ctx, "starting listening http srv at "+s.cfg.HTTPAddress)
		if err := s.srv.ListenAndServe(); err != http.ErrServerClosed {
			logger.Fatal(ctx, "error start http srv, err: %+v", err)
		}
	}()

	closer.Add(s.Close)
}

// Close - the method to close the http-server.
func (s *Server) Close() error {
	ctx := context.TODO()
	if err := s.srv.Shutdown(ctx); err != nil {
		logger.Error(ctx, "error stop http srv, err", err)
		return errors.Wrap(err, "close server error")
	}

	logger.Info(ctx, "http server shutdown done")

	return nil
}
