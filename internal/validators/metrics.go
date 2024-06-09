package validators

import (
	"encoding/json"
	"io"
	"metrix/internal/model"
	"strconv"

	"github.com/go-playground/validator/v10"
)

type MetricValidator struct {
	validate *validator.Validate
}

func NewMetricValidator() *MetricValidator {
	return &MetricValidator{
		validate: validator.New(validator.WithRequiredStructEnabled()),
	}
}

func (v *MetricValidator) FromVars(vars map[string]string) (*model.Metric, error) {
	metric := &model.Metric{}

	// Retrieveing variables
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
	}

	// Validate structure
	err := v.validate.Struct(metric)
	if err != nil {
		return nil, NewParsingValueError("failed to validate metric: %s", err)
	}

	return metric, nil
}

func (v *MetricValidator) FromBody(body io.ReadCloser) (*model.Metric, error) {
	metric := &model.Metric{}

	err := json.NewDecoder(body).Decode(metric)
	if err != nil {
		return nil, NewParsingValueError("failed to parse metric json: %s", err)
	}

	if metric.MType != model.CounterType && metric.MType != model.GaugeType {
		return nil, NewParsingValueError("failed to validate metric type: %s", err)
	}

	return metric, nil
}
