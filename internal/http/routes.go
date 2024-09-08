package http

import (
	"net/http"

	"metrix/internal/middlewares"

	"github.com/gorilla/mux"
)

func (s *Server) setupRoutes() *mux.Router {
	m := mux.NewRouter()

	// Health handlers
	m.HandleFunc("/ping", s.health.PingDB)
	m.HandleFunc("/liveness", s.health.LivenessState)
	m.HandleFunc("/readiness", s.health.ReadinessState)

	// Metrics handlers
	m.HandleFunc("/", s.metrics.GetIDs)

	m.HandleFunc("/update/{type}/{id}/{value}", s.metrics.Set)
	m.HandleFunc("/value/{type}/{id}", s.metrics.Get)

	m.HandleFunc("/update/", s.metrics.SetWithModel).
		Methods(http.MethodPost).
		Headers("Content-Type", "application/json")
	m.HandleFunc("/value/", s.metrics.GetWithModel).
		Methods(http.MethodPost).
		Headers("Content-Type", "application/json")
	m.HandleFunc("/updates/", s.metrics.SetMany).
		Methods(http.MethodPost).
		Headers("Content-Type", "application/json")

	m.Use(middlewares.LoggingMiddleware)
	m.Use(middlewares.SignatureMiddleware)
	m.Use(middlewares.GzipMiddleware)

	// pprof
	m.PathPrefix("/debug/pprof/").Handler(http.DefaultServeMux)

	return m
}
