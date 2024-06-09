package controllers

import (
	"context"
	"metrix/internal/model"
)

type MetricsController interface {
	Set(ctx context.Context, metricIn *model.Metric) (*model.Metric, error)
	Get(ctx context.Context, metricID string) (*model.Metric, error)
	GetIDs(ctx context.Context) (*[]string, error)
}

type HealthController interface {
	SetReadiness(state bool)
	SetLiveness(state bool)
	ReadinessState() bool
	LivenessState() bool
	PingDB() bool
}
