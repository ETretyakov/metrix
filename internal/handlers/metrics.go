package handlers

import (
	"errors"
	"fmt"
	"metrix/internal/controllers"
	"metrix/internal/logger"
	"metrix/internal/repository"
	"metrix/internal/validators"
	"net/http"

	"github.com/gorilla/mux"
)

type MetricsHandlers struct {
	controller controllers.MetricsController
}

func NewMetricsHandlers(repoGroup *repository.Group) *MetricsHandlers {
	controller := controllers.NewMetricController(repoGroup)
	return &MetricsHandlers{controller: controller}
}

func (h *MetricsHandlers) Set(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	metric, err := h.controller.Set(ctx, vars)
	if err != nil {
		var parsingValueError validators.ParsingValueError
		if errors.As(err, &parsingValueError) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		logger.Warn(ctx, fmt.Sprintf("failed to trigger controller: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(metric.GetValue()))
}

func (h *MetricsHandlers) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	metricID, ok := vars["metricID"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	metric, err := h.controller.Get(ctx, metricID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if metric == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(metric.GetValue()))
}
