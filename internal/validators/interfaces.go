package validators

import (
	"io"
	"metrix/internal/model"
)

type MetricsValidator interface {
	FromVars(vars map[string]string) (*model.Metric, error)
	FromBody(body io.ReadCloser) (*model.Metric, error)
	ManyFromBody(body io.ReadCloser) ([]*model.Metric, error)
}
