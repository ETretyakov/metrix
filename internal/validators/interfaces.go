package validators

import (
	"io"

	"metrix/internal/model"
)

// MetricsValidator - the interfeace to describe validator for metrics.
type MetricsValidator interface {
	FromVars(vars map[string]string) (*model.Metric, error)
	FromBody(body io.ReadCloser) (*model.Metric, error)
	ManyFromBody(body io.ReadCloser) ([]*model.Metric, error)
}
