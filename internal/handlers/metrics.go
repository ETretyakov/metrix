package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"metrix/internal/controllers"
	"metrix/internal/repository"
	"metrix/internal/validators"
	"metrix/pkg/logger"

	"github.com/gorilla/mux"
)

const (
	parseErrMsg string = "failed to parse payload: %v"
)

// MetricsHandlers - the implementation structure for the MetricsHandlers that manages
// access to the related controller.
type MetricsHandlers struct {
	controller controllers.MetricsController
	validator  validators.MetricsValidator
}

// NewMetricsHandlers - the builder function for the MetricsHandlers.
func NewMetricsHandlers(repoGroup *repository.Group) *MetricsHandlers {
	return &MetricsHandlers{
		controller: controllers.NewMetricController(repoGroup),
		validator:  validators.NewMetricsValidator(),
	}
}

// Set - the handler method that incapsulates validation logic for setting metrics
// values functionality.
func (h *MetricsHandlers) Set(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	metricIn, err := h.validator.FromVars(vars)
	if err != nil {
		var parsingValueError validators.ParsingValueError
		if errors.As(err, &parsingValueError) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		logger.Warn(ctx, fmt.Sprintf(parseErrMsg, err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	metric, err := h.controller.Set(ctx, metricIn)
	if err != nil {
		logger.Warn(ctx, fmt.Sprintf("failed to trigger controller: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = w.Write([]byte(metric.GetValue()))
	if err != nil {
		logger.Error(ctx, "failed to trigger controller", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// SetWithModel - the handler method that incapsulates validation logic for setting metrics
// values functionality in batching mode with dedicated model.
func (h *MetricsHandlers) SetWithModel(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctx := r.Context()
	metricIn, err := h.validator.FromBody(r.Body)
	if err != nil {
		var parsingValueError validators.ParsingValueError
		if errors.As(err, &parsingValueError) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		logger.Warn(ctx, fmt.Sprintf(parseErrMsg, err))
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

// SetMany - the handler method that incapsulates validation logic for setting metrics
// values functionality in batching mode.
func (h *MetricsHandlers) SetMany(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctx := r.Context()
	metricsIn, err := h.validator.ManyFromBody(r.Body)
	if err != nil {
		var parsingValueError validators.ParsingValueError
		if errors.As(err, &parsingValueError) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		logger.Warn(ctx, fmt.Sprintf(parseErrMsg, err))
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

// Get - the handler method that incapsulates validation logic for getting metrics.
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

	_, err = w.Write([]byte(metric.GetValue()))
	if err != nil {
		logger.Error(ctx, "failed to trigger controller", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// GetWithModel - the handler method that incapsulates validation logic for getting
// metrics with dedicated model.
func (h *MetricsHandlers) GetWithModel(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctx := r.Context()

	metricIn, err := h.validator.FromBody(r.Body)
	if err != nil {
		var parsingValueError validators.ParsingValueError
		if errors.As(err, &parsingValueError) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		logger.Warn(ctx, fmt.Sprintf(parseErrMsg, err))
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

// GetIDs - the handler method that incapsulates validation logic for getting metrics ids.
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

	_, err = w.Write([]byte(body))
	if err != nil {
		logger.Error(ctx, "failed to trigger controller", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
