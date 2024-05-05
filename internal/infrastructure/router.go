package infrastructure

import (
	"fmt"
	"metrix/internal/interfaces"
	"metrix/internal/usecases"
	"net/http"

	"github.com/gorilla/mux"
)

func Dispatch(logger usecases.Logger, storageHandler interfaces.StorageHandler) {
	widgetController := interfaces.NewWidgetController(storageHandler, logger)

	router := mux.NewRouter()

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

	addr := ":8080"
	logger.LogAccess(fmt.Sprintf("starting server at: %s", addr))
	if err := http.ListenAndServe(addr, router); err != nil {
		logger.LogError("%s", err)
	}
}
