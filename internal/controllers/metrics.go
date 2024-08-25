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

func (m *MetricControllerImpl) SetMany(
	ctx context.Context,
	metricsIn []*model.Metric,
) (bool, error) {
	collapsedMapping := map[string]model.Metric{}
	for _, metric := range metricsIn {
		collapsed, ok := collapsedMapping[metric.ID]
		if ok {
			collapsed.SetValue(metric.Delta, metric.Value)
			collapsedMapping[metric.ID] = collapsed
		} else {
			collapsedMapping[metric.ID] = *metric
		}
	}

	metricIDs := []string{}
	for _, m := range collapsedMapping {
		metricIDs = append(metricIDs, m.ID)
	}

	curMetrics, err := m.repoGroup.MetricRepo.ReadMany(ctx, metricIDs)
	if err != nil {
		return false, fmt.Errorf("failed to get current metrics: %w", err)
	}
	mapCurMetrics := map[string]model.Metric{}
	if curMetrics != nil {
		for _, metric := range *curMetrics {
			mapCurMetrics[metric.ID] = metric
		}
	}

	newMetricsIn := []model.Metric{}
	for _, metric := range collapsedMapping {
		curMetric, ok := mapCurMetrics[metric.ID]
		if ok {
			curMetric.SetValue(metric.Delta, metric.Value)
			newMetricsIn = append(newMetricsIn, curMetric)
		} else {
			newMetricsIn = append(newMetricsIn, metric)
		}
	}

	status, err := m.repoGroup.MetricRepo.UpsertMany(ctx, newMetricsIn)
	if err != nil {
		logger.Debug(ctx, fmt.Sprintf("failed to upsert metrics: %s", err))
		return false, fmt.Errorf("failed to upsert metrics: %w", err)
	}

	return status, nil
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
