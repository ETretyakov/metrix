package http

import (
	"github.com/gorilla/mux"
)

func (s *Server) setupRoutes() *mux.Router {
	mux := mux.NewRouter()

	// Health handlers
	mux.HandleFunc("/liveness", s.health.LivenessState)
	mux.HandleFunc("/readiness", s.health.ReadinessState)

	// Metrics handlers
	mux.HandleFunc("/update/{mtype}/{metricID}/{value}", s.metrics.Set)
	mux.HandleFunc("/value/{mtype}/{metricID}", s.metrics.Get)

	return mux
}
