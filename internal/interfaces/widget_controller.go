package interfaces

import (
	"encoding/json"
	"metrix/internal/domain"
	"metrix/internal/usecases"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type WidgetController struct {
	WidgetInteractor usecases.WidgetInteractor
	Logger           usecases.Logger
}

func NewWidgetController(storageHandler StorageHandler, logger usecases.Logger) *WidgetController {
	return &WidgetController{
		WidgetInteractor: usecases.WidgetInteractor{
			WidgetRepository: &WidgetRepository{
				StorageHandler: storageHandler,
			},
		},
		Logger: logger,
	}
}

func (wc *WidgetController) Show(w http.ResponseWriter, r *http.Request) {
	wc.Logger.LogAccess("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	vars := mux.Vars(r)

	namespace := "default"
	widgetType, err := domain.ParseWidgetType(vars["widgetType"])
	if err != nil {
		wc.Logger.LogError("[ERROR] %s %s %s: %s\n", r.RemoteAddr, r.Method, r.URL, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	name := vars["name"]

	widget, err := wc.WidgetInteractor.Show(namespace, widgetType, name)
	if err != nil {
		wc.Logger.LogError("[ERROR] %s %s %s: %s\n", r.RemoteAddr, r.Method, r.URL, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(widget)
	if err != nil {
		wc.Logger.LogError("[ERROR] %s %s %s: %s\n", r.RemoteAddr, r.Method, r.URL, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (wc *WidgetController) Update(w http.ResponseWriter, r *http.Request) {
	wc.Logger.LogAccess("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	vars := mux.Vars(r)

	namespace := "default"
	widgetType, err := domain.ParseWidgetType(vars["widgetType"])
	if err != nil {
		wc.Logger.LogError("[ERROR] %s %s %s: %s\n", r.RemoteAddr, r.Method, r.URL, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	name := vars["name"]
	value := vars["value"]

	val, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		wc.Logger.LogError("[ERROR] %s %s %s: %s\n", r.RemoteAddr, r.Method, r.URL, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var widget domain.Widget

	switch widgetType {
	case domain.CounterWidget:
		widget, err = wc.WidgetInteractor.Increment(namespace, widgetType, name, val)
	default:
		widget, err = wc.WidgetInteractor.Update(namespace, widgetType, name, val)
	}
	if err != nil {
		wc.Logger.LogError("[ERROR] %s %s %s: %s\n", r.RemoteAddr, r.Method, r.URL, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(widget)
	if err != nil {
		wc.Logger.LogError("[ERROR] %s %s %s: %s\n", r.RemoteAddr, r.Method, r.URL, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
