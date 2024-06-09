package controllers

import (
	"context"
	"fmt"
	"metrix/internal/model"
	"metrix/internal/repository"
	"metrix/pkg/logger"
)

type MetricControllerImpl struct {
	repoGroup *repository.Group
}

func NewMetricController(repoGroup *repository.Group) *MetricControllerImpl {
	return &MetricControllerImpl{repoGroup: repoGroup}
}

func (m *MetricControllerImpl) Set(
	ctx context.Context,
	metricIn *model.Metric,
) (*model.Metric, error) {
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

func (m *MetricControllerImpl) GetIDs(ctx context.Context) (*[]string, error) {
	ids, err := m.repoGroup.MetricRepo.ReadIDs(ctx)
	if err != nil {
		logger.Debug(ctx, fmt.Sprintf("failed to retrieve ids: %s", err))
		return nil, fmt.Errorf("failed to retrieve ids: %w", err)
	}

	return ids, nil
}
