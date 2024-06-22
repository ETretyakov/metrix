package infrastructure

import (
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
		"/show/{widgetType}/{name}",
		widgetController.Show,
	).Methods(http.MethodGet)

	if err := http.ListenAndServe(":8080", router); err != nil {
		logger.LogError("%s", err)
	}
}
