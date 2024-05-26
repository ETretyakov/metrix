package interfaces

import (
	"encoding/json"
	"errors"
	"metrix/internal/domain"
	"metrix/internal/exceptions"
	"metrix/internal/logger"
	"metrix/internal/usecases"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

type WidgetController struct {
	WidgetInteractor usecases.WidgetInteractor
}

func NewWidgetController(storageHandler StorageHandler) *WidgetController {
	return &WidgetController{
		WidgetInteractor: usecases.WidgetInteractor{
			WidgetRepository: &WidgetRepository{
				StorageHandler: storageHandler,
			},
		},
	}
}

func (wc *WidgetController) Show(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	vars := mux.Vars(r)

	namespace := "default"
	widgetType, err := domain.ParseWidgetType(vars["widgetType"])
	if err != nil {
		logger.Log.Errorw(
			err.Error(),
			"address", r.RemoteAddr,
			"method", r.Method,
			"url", r.URL,
		)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	name := vars["name"]

	widget, err := wc.WidgetInteractor.Show(namespace, widgetType, name)
	if err != nil {
		logger.Log.Errorw(
			err.Error(),
			"address", r.RemoteAddr,
			"method", r.Method,
			"url", r.URL,
		)
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
}

func (wc *WidgetController) ShowSingleEndpoint(
	w http.ResponseWriter,
	r *http.Request,
) {
	w.Header().Set("Content-Type", "application/json")

	var schemeIn domain.Metrics

	err := json.NewDecoder(r.Body).Decode(&schemeIn)
	if err != nil {
		logger.Log.Errorw(
			err.Error(),
			"address", r.RemoteAddr,
			"method", r.Method,
			"url", r.URL,
		)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	namespace := "default"
	widgetType, err := domain.ParseWidgetType(schemeIn.MType)
	if err != nil {
		logger.Log.Errorw(
			err.Error(),
			"address", r.RemoteAddr,
			"method", r.Method,
			"url", r.URL,
		)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	widget, err := wc.WidgetInteractor.Show(namespace, widgetType, schemeIn.ID)
	if err != nil {
		logger.Log.Errorw(
			err.Error(),
			"address", r.RemoteAddr,
			"method", r.Method,
			"url", r.URL,
		)
		var recordNotFound exceptions.RecordNotFoundError
		if errors.As(err, &recordNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	schemeOut := domain.Metrics{
		ID:    widget.Name,
		MType: string(widget.Type),
	}

	switch widgetType {
	case domain.CounterWidget:
		val := int64(widget.Value)
		schemeOut.Delta = &val
	default:
		schemeOut.Value = &widget.Value
	}

	err = json.NewEncoder(w).Encode(schemeOut)
	if err != nil {
		logger.Log.Errorw(
			err.Error(),
			"address", r.RemoteAddr,
			"method", r.Method,
			"url", r.URL,
		)
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

func (wc *WidgetController) Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	vars := mux.Vars(r)

	namespace := "default"
	widgetType, err := domain.ParseWidgetType(vars["widgetType"])
	if err != nil {
		logger.Log.Errorw(
			err.Error(),
			"address", r.RemoteAddr,
			"method", r.Method,
			"url", r.URL,
		)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	name := vars["name"]
	value := vars["value"]

	val, err := strconv.ParseFloat(value, 64)
	if err != nil {
		logger.Log.Errorw(
			err.Error(),
			"address", r.RemoteAddr,
			"method", r.Method,
			"url", r.URL,
		)
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
		logger.Log.Errorw(
			err.Error(),
			"address", r.RemoteAddr,
			"method", r.Method,
			"url", r.URL,
		)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(strconv.FormatFloat(widget.Value, 'f', -1, 64)))
}

func (wc *WidgetController) UpdateSingleEndpoint(
	w http.ResponseWriter,
	r *http.Request,
) {
	w.Header().Set("Content-Type", "application/json")

	var schemeIn domain.Metrics

	err := json.NewDecoder(r.Body).Decode(&schemeIn)
	if err != nil {
		logger.Log.Errorw(
			err.Error(),
			"address", r.RemoteAddr,
			"method", r.Method,
			"url", r.URL,
		)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	namespace := "default"
	widgetType, err := domain.ParseWidgetType(schemeIn.MType)
	if err != nil {
		logger.Log.Errorw(
			err.Error(),
			"address", r.RemoteAddr,
			"method", r.Method,
			"url", r.URL,
		)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var widget domain.Widget

	switch widgetType {
	case domain.CounterWidget:
		if schemeIn.Delta == nil {
			err = errors.New("empty field for delta")
		} else {
			widget, err = wc.WidgetInteractor.Increment(
				namespace,
				widgetType,
				schemeIn.ID,
				float64(*schemeIn.Delta),
			)
		}
	default:
		if schemeIn.Value == nil {
			err = errors.New("empty field for value")
		} else {
			widget, err = wc.WidgetInteractor.Update(
				namespace,
				widgetType,
				schemeIn.ID,
				*schemeIn.Value,
			)
		}
	}
	if err != nil {
		logger.Log.Errorw(
			err.Error(),
			"address", r.RemoteAddr,
			"method", r.Method,
			"url", r.URL,
		)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	schemeOut := domain.Metrics{
		ID:    widget.Name,
		MType: string(widget.Type),
	}

	switch widgetType {
	case domain.CounterWidget:
		val := int64(widget.Value)
		schemeOut.Delta = &val
	default:
		schemeOut.Value = &widget.Value
	}

	err = json.NewEncoder(w).Encode(schemeOut)
	if err != nil {
		logger.Log.Errorw(
			err.Error(),
			"address", r.RemoteAddr,
			"method", r.Method,
			"url", r.URL,
		)
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

func (wc *WidgetController) Keys(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
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
		logger.Log.Errorw(
			err.Error(),
			"address", r.RemoteAddr,
			"method", r.Method,
			"url", r.URL,
		)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
