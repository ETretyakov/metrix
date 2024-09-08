package controllers

import (
	"context"

	"metrix/internal/model"
)

// MetricsController - the interface that describes all the MetricsController methods.
type MetricsController interface {
	Set(ctx context.Context, metricIn *model.Metric) (*model.Metric, error)
	SetMany(ctx context.Context, metricsIn []*model.Metric) (bool, error)
	Get(ctx context.Context, metricID string) (*model.Metric, error)
	GetIDs(ctx context.Context) (*[]string, error)
}

// HealthController - the interface that describes all the HealthController methods.
type HealthController interface {
	SetReadiness(state bool)
	SetLiveness(state bool)
	ReadinessState() bool
	LivenessState() bool
	PingDB() bool
}
