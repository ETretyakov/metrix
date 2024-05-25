package infrastructure

import (
	"metrix/internal/interfaces"
	"metrix/internal/logger"
	"metrix/internal/middlewares"
	"net/http"

	"github.com/gorilla/mux"
)

func Dispatch(
	addr string,
	storageHandler interfaces.StorageHandler,
) {
	widgetController := interfaces.NewWidgetController(storageHandler)

	router := mux.NewRouter()

	router.HandleFunc(
		"/update/",
		widgetController.UpdateSingleEndpoint,
	).
		Methods(http.MethodPost).
		Headers(
			"Content-Type", "application/json",
			"Accept-Encoding", "gzip",
		)

	router.HandleFunc(
		"/value/",
		widgetController.ShowSingleEndpoint,
	).
		Methods(http.MethodPost).
		Headers(
			"Content-Type", "application/json",
			"Accept-Encoding", "gzip",
		)

	router.HandleFunc(
		"/update/{widgetType}/{name}/{value}",
		widgetController.Update,
	).Methods(http.MethodPost)

	router.HandleFunc(
		"/value/{widgetType}/{name}",
		widgetController.Show,
	).Methods(http.MethodGet)

	router.HandleFunc(
		"/",
		widgetController.Keys,
	).Methods(http.MethodGet)

	router.Use(logger.LoggingMiddleware)
	router.Use(middlewares.GzipMiddleware)

	logger.Log.Infof("starting server at: %s", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		logger.Log.Errorw(err.Error(), "address", addr)
	}
}
