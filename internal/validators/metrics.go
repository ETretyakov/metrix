package validators

import (
	"encoding/json"
	"io"
	"strconv"

	"metrix/internal/model"

	"github.com/go-playground/validator/v10"
)

// MetricsValidatorImpl - the implementation structure for validations.
type MetricsValidatorImpl struct {
	validate *validator.Validate
}

// NewMetricsValidator - the builder function for MetricsValidatorImpl.
func NewMetricsValidator() *MetricsValidatorImpl {
	return &MetricsValidatorImpl{
		validate: validator.New(validator.WithRequiredStructEnabled()),
	}
}

// FromVars - the function that parses metric structure from map object.
func (v *MetricsValidatorImpl) FromVars(vars map[string]string) (*model.Metric, error) {
	metric := &model.Metric{}

	// Retrieving variables
	metricID, ok := vars["id"]
	if !ok {
		return nil, NewParsingValueError("failed to retrieve metricID path param")
	}

	mtype, ok := vars["type"]
	if !ok {
		return nil, NewParsingValueError("failed to retrieve mtype path param")
	}

	// Assigning values
	metric.ID = metricID

	if mtype == string(model.CounterType) {
		metric.MType = model.CounterType

		val, ok := vars["delta"]
		if !ok {
			val, ok = vars["value"]
			if !ok {
				return nil, NewParsingValueError("failed to retrieve value path param")
			}
		}

		delta, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return nil, NewParsingValueError("failed to parse value: %s", err)
		}

		metric.Delta = &delta
		metric.Value = nil
	}

	if mtype == string(model.GaugeType) {
		metric.MType = model.GaugeType

		val, ok := vars["value"]
		if !ok {
			return nil, NewParsingValueError("failed to retrieve value path param")
		}

		value, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return nil, NewParsingValueError("failed to parse value: %s", err)
		}

		metric.Value = &value
		metric.Delta = nil
	}

	// Validate structure
	err := v.validate.Struct(metric)
	if err != nil {
		return nil, NewParsingValueError("failed to validate metric: %s", err)
	}

	return metric, nil
}

// FromBody - the function that parses metric structure from reader.
func (v *MetricsValidatorImpl) FromBody(body io.ReadCloser) (*model.Metric, error) {
	metric := &model.Metric{}

	err := json.NewDecoder(body).Decode(metric)
	if err != nil {
		return nil, NewParsingValueError("failed to parse metric json: %s", err)
	}

	if metric.MType != model.CounterType && metric.MType != model.GaugeType {
		return nil, NewParsingValueError("failed to validate metric type: %s", err)
	}

	if metric.MType == model.CounterType {
		metric.Value = nil
	}

	if metric.MType == model.GaugeType {
		metric.Delta = nil
	}

	return metric, nil
}

// ManyFromBody - the function that parses many metric structures from reader.
func (v *MetricsValidatorImpl) ManyFromBody(body io.ReadCloser) ([]*model.Metric, error) {
	metrics := []*model.Metric{}

	err := json.NewDecoder(body).Decode(&metrics)
	if err != nil {
		return nil, NewParsingValueError("failed to parse metric json: %s", err)
	}

	for _, m := range metrics {
		if m.MType != model.CounterType && m.MType != model.GaugeType {
			return nil, NewParsingValueError("failed to validate metric type: %s", err)
		}

		if m.MType == model.CounterType {
			m.Value = nil
		}

		if m.MType == model.GaugeType {
			m.Delta = nil
		}
	}

	return metrics, nil
}
