package main

import (
	"net/http"
	"os"
	"time"

	"github.com/ETretyakov/metrix/internal/handlers"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	router := mux.NewRouter()

	router.HandleFunc("/health", handlers.HealthHandler).
		Methods(http.MethodGet)

	router.HandleFunc("/update/counter/{name}/{value}", handlers.CounterWidgetUpdateHandler).
		Methods(http.MethodPost)

	router.HandleFunc("/update/gauge/{name}/{value}", handlers.GaugeWidgetUpdateHandler).
		Methods(http.MethodPost)

	router.HandleFunc("/update/{unkown}/{name}/{value}", handlers.UnknownMetricHandler).
		Methods(http.MethodPost)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Info().Msg("[main] starting server")
	server.ListenAndServe()
}
