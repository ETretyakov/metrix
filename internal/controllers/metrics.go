package controllers

import (
	"context"
	"fmt"
	"metrix/internal/model"
	"metrix/internal/repository"
	"metrix/internal/validators"
	"metrix/pkg/logger"
)

var (
	metricValidator = validators.NewMetricValidator()
)

type MetricControllerImpl struct {
	repoGroup *repository.Group
}

func NewMetricController(repoGroup *repository.Group) *MetricControllerImpl {
	return &MetricControllerImpl{repoGroup: repoGroup}
}

func (m *MetricControllerImpl) Set(
	ctx context.Context,
	vars map[string]string,
) (*model.Metric, error) {
	metricIn, err := metricValidator.FromVars(vars)
	if err != nil {
		logger.Debug(ctx, fmt.Sprintf("failed to parse structure: %s", err))
		return nil, fmt.Errorf("failed to parse metric: %w", err)
	}

	metric, err := m.repoGroup.MetricRepo.Read(ctx, metricIn.ID)
	if err != nil {
		logger.Debug(ctx, fmt.Sprintf("failed to retrieve structure: %s", err))
		return nil, fmt.Errorf("failed to retrieve metric: %w", err)
	}

	if metric != nil {
		metric.SetValue(metricIn.Delta, metricIn.Value)
		metric, err = m.repoGroup.MetricRepo.Update(ctx, metric)
		if err != nil {
			logger.Debug(ctx, fmt.Sprintf("failed to update metric: %s", err))
			return nil, fmt.Errorf("failed to update metric: %w", err)
		}
	} else {
		metric, err = m.repoGroup.MetricRepo.Create(ctx, metricIn)
		if err != nil {
			logger.Debug(ctx, fmt.Sprintf("failed to insert metric: %s", err))
			return nil, fmt.Errorf("failed to insert metric: %w", err)
		}
	}

	return metric, nil
}

func (m *MetricControllerImpl) Get(
	ctx context.Context,
	metricID string,
) (*model.Metric, error) {
	metric, err := m.repoGroup.MetricRepo.Read(ctx, metricID)
	if err != nil {
		logger.Debug(ctx, fmt.Sprintf("failed to retrieve structure: %s", err))
		return nil, fmt.Errorf("failed to retrieve metric: %w", err)
	}

	return metric, nil
}
