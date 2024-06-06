package infrastructure

import (
	"context"
	"metrix/internal/interfaces"
	"metrix/internal/logger"
	"metrix/internal/middlewares"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func Dispatch(
	ctx context.Context,
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
		)

	router.HandleFunc(
		"/value/",
		widgetController.ShowSingleEndpoint,
	).
		Methods(http.MethodPost).
		Headers(
			"Content-Type", "application/json",
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

	server := http.Server{
		Addr:    addr,
		Handler: router,
	}

	serverErr := make(chan error, 1)
	go func() {
		serverErr <- server.ListenAndServe()
	}()

	var err error
	select {
	case <-ctx.Done():
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		err = server.Shutdown(ctx)
	case err = <-serverErr:
		logger.Log.Errorw(err.Error(), "address", addr)
	}
}
