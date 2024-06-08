package controllers

import (
	"context"
	"metrix/internal/model"
)

type MetricsController interface {
	Set(ctx context.Context, vars map[string]string) (*model.Metric, error)
	Get(ctx context.Context, metricID string) (*model.Metric, error)
}

type HealthController interface {
	SetReadiness(state bool)
	SetLiveness(state bool)
	ReadinessState() bool
	LivenessState() bool
}
