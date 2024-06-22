package http

import (
	"net/http"

	"metrix/internal/middlewares"

	"github.com/gorilla/mux"
)

func (s *Server) setupRoutes() *mux.Router {
	mux := mux.NewRouter()

	// Health handlers
	mux.HandleFunc("/ping", s.health.PingDB)
	mux.HandleFunc("/liveness", s.health.LivenessState)
	mux.HandleFunc("/readiness", s.health.ReadinessState)

	// Metrics handlers
	mux.HandleFunc("/", s.metrics.GetIDs)

	mux.HandleFunc("/update/{type}/{id}/{value}", s.metrics.Set)
	mux.HandleFunc("/value/{type}/{id}", s.metrics.Get)

	mux.HandleFunc("/update/", s.metrics.SetWithModel).
		Methods(http.MethodPost).
		Headers("Content-Type", "application/json")
	mux.HandleFunc("/value/", s.metrics.GetWithModel).
		Methods(http.MethodPost).
		Headers("Content-Type", "application/json")
	mux.HandleFunc("/updates/", s.metrics.SetMany).
		Methods(http.MethodPost).
		Headers("Content-Type", "application/json")

	mux.Use(middlewares.LoggingMiddleware)
	mux.Use(middlewares.SignatureMiddleware)
	mux.Use(middlewares.GzipMiddleware)

	return mux
}
