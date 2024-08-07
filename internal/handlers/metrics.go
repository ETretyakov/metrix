package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"metrix/internal/controllers"
	"metrix/internal/repository"
	"metrix/internal/validators"
	"metrix/pkg/logger"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

var (
	metricValidator = validators.NewMetricValidator()
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

	metricIn, err := metricValidator.FromVars(vars)
	if err != nil {
		var parsingValueError validators.ParsingValueError
		if errors.As(err, &parsingValueError) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		logger.Warn(ctx, fmt.Sprintf("failed to parse payload: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	metric, err := h.controller.Set(ctx, metricIn)
	if err != nil {
		logger.Warn(ctx, fmt.Sprintf("failed to trigger controller: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(metric.GetValue()))
}

func (h *MetricsHandlers) SetWithModel(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctx := r.Context()
	metricIn, err := metricValidator.FromBody(r.Body)
	if err != nil {
		var parsingValueError validators.ParsingValueError
		if errors.As(err, &parsingValueError) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		logger.Warn(ctx, fmt.Sprintf("failed to parse payload: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	metric, err := h.controller.Set(ctx, metricIn)
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

	err = json.NewEncoder(w).Encode(metric)
	if err != nil {
		logger.Error(
			ctx,
			"failed to encode response json",
			err,
			"address", r.RemoteAddr,
			"method", r.Method,
			"url", r.URL,
		)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *MetricsHandlers) SetMany(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctx := r.Context()
	metricsIn, err := metricValidator.ManyFromBody(r.Body)
	if err != nil {
		var parsingValueError validators.ParsingValueError
		if errors.As(err, &parsingValueError) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		logger.Warn(ctx, fmt.Sprintf("failed to parse payload: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(metricsIn) == 0 {
		logger.Info(ctx, "payload is an empty array")
		w.WriteHeader(http.StatusOK)
		return
	}

	_, err = h.controller.SetMany(ctx, metricsIn)
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
}

func (h *MetricsHandlers) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	metricID, ok := vars["id"]
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

func (h *MetricsHandlers) GetWithModel(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctx := r.Context()

	metricIn, err := metricValidator.FromBody(r.Body)
	if err != nil {
		var parsingValueError validators.ParsingValueError
		if errors.As(err, &parsingValueError) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		logger.Warn(ctx, fmt.Sprintf("failed to parse payload: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	metric, err := h.controller.Get(ctx, metricIn.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if metric == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err = json.NewEncoder(w).Encode(metric)
	if err != nil {
		logger.Error(
			ctx,
			"failed to encode response json",
			err,
			"address", r.RemoteAddr,
			"method", r.Method,
			"url", r.URL,
		)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *MetricsHandlers) GetIDs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	ctx := r.Context()

	ids, err := h.controller.GetIDs(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	body := ""
	if ids != nil {
		body = strings.Join(*ids, "\n")
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(body))
}
