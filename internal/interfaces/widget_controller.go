package interfaces

import (
	"encoding/json"
	"errors"
	"metrix/internal/domain"
	"metrix/internal/exceptions"
	"metrix/internal/usecases"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

const (
	errorMsg = "[ERROR] %s %s %s: %s\n"
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
		wc.Logger.LogError(errorMsg, r.RemoteAddr, r.Method, r.URL, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	name := vars["name"]

	widget, err := wc.WidgetInteractor.Show(namespace, widgetType, name)
	if err != nil {
		wc.Logger.LogError(errorMsg, r.RemoteAddr, r.Method, r.URL, err)
		var recordNotFound exceptions.RecordNotFoundError
		if errors.As(err, &recordNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(strconv.FormatFloat(widget.Value, 'f', -1, 64)))

	// w.WriteHeader(http.StatusOK)
	// err = json.NewEncoder(w).Encode(widget)
	// if err != nil {
	// 	wc.Logger.LogError(errorMsg, r.RemoteAddr, r.Method, r.URL, err)
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }
}

func (wc *WidgetController) Update(w http.ResponseWriter, r *http.Request) {
	wc.Logger.LogAccess("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	vars := mux.Vars(r)

	namespace := "default"
	widgetType, err := domain.ParseWidgetType(vars["widgetType"])
	if err != nil {
		wc.Logger.LogError(errorMsg, r.RemoteAddr, r.Method, r.URL, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	name := vars["name"]
	value := vars["value"]

	val, err := strconv.ParseFloat(value, 64)
	if err != nil {
		wc.Logger.LogError(errorMsg, r.RemoteAddr, r.Method, r.URL, err)
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
		wc.Logger.LogError(errorMsg, r.RemoteAddr, r.Method, r.URL, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(strconv.FormatFloat(widget.Value, 'f', -1, 64)))

	// err = json.NewEncoder(w).Encode(widget)
	// if err != nil {
	// 	wc.Logger.LogError(errorMsg, r.RemoteAddr, r.Method, r.URL, err)
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }
}

func (wc *WidgetController) Keys(w http.ResponseWriter, r *http.Request) {
	wc.Logger.LogAccess("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	namespace := "default"

	keys, err := wc.WidgetInteractor.Keys(namespace)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	count := 0
	items := make([]domain.WidgetListItem, 0)
	for i, k := range keys {
		splitKey := strings.Split(k, ":")
		items = append(
			items,
			domain.WidgetListItem{
				Namespace: splitKey[0],
				Type:      domain.WidgetType(splitKey[1]),
				Name:      splitKey[2],
			},
		)
		count = i + 1
	}

	widgetList := domain.WidgetList{Count: count, Items: items}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(widgetList)
	if err != nil {
		wc.Logger.LogError(errorMsg, r.RemoteAddr, r.Method, r.URL, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}