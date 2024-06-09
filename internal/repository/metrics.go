package repository

import (
	"context"
	"database/sql"
	"fmt"
	"metrix/internal/model"

	"github.com/doug-martin/goqu/v9"
)

const metricTName = "mtr_metrics"

type MetricRepositoryImpl struct {
	gr *Group
}

func NewMetricRepository(db *Group) *MetricRepositoryImpl {
	return &MetricRepositoryImpl{
		gr: db,
	}
}

func (r *MetricRepositoryImpl) Create(
	ctx context.Context,
	metric *model.Metric,
) (*model.Metric, error) {
	qu, _, err := goqu.Insert(metricTName).Rows(metric).Returning("id").ToSQL()
	if err != nil {
		return nil, fmt.Errorf("create metric error during query building: %w", err)
	}

	tx, err := r.gr.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Commit()

	if _, err := tx.ExecContext(ctx, qu); err != nil {
		return nil, fmt.Errorf("create metric error during execute query: %w", err)
	}

	return metric, nil
}

func (r *MetricRepositoryImpl) Read(
	ctx context.Context,
	metricID string,
) (*model.Metric, error) {
	qu, _, err := goqu.
		Select(&model.Metric{}).
		From(metricTName).
		Where(goqu.Ex{"id": metricID}).
		Limit(1).
		ToSQL()
	if err != nil {
		return nil, fmt.Errorf("read metric error during query building: %w", err)
	}

	var newMetric model.Metric
	row := r.gr.DB.QueryRowxContext(ctx, qu)
	err = row.StructScan(&newMetric)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("read metric error during scan row: %w", err)
	}

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &newMetric, nil
}

func (r *MetricRepositoryImpl) ReadIDs(
	ctx context.Context,
) (*[]string, error) {
	qu, _, err := goqu.
		Select("id").
		From(metricTName).
		ToSQL()
	if err != nil {
		return nil, fmt.Errorf("read ids metric error during query building: %w", err)
	}

	var ids []string
	row := r.gr.DB.QueryRowxContext(ctx, qu)
	err = row.StructScan(&ids)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("read metric error during scan row: %w", err)
	}

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &ids, nil
}

func (r *MetricRepositoryImpl) Update(
	ctx context.Context,
	metric *model.Metric,
) (*model.Metric, error) {
	qu, _, err := goqu.
		Update(metricTName).
		Set(metric).
		Where(goqu.Ex{"id": metric.ID}).
		ToSQL()
	if err != nil {
		return nil, fmt.Errorf("create metric error during query building: %w", err)
	}

	if _, err := r.gr.ExecContext(ctx, qu); err != nil {
		return nil, fmt.Errorf("update metric error during execute query: %w", err)
	}

	metricOut, err := r.Read(ctx, metric.ID)
	if err != nil {
		return nil, fmt.Errorf("refresh metric error: %w", err)
	}

	return metricOut, nil
}

func (r *MetricRepositoryImpl) Delete(
	ctx context.Context,
	metricID string,
) error {
	qu, _, err := goqu.
		Delete(metricTName).
		Where(goqu.Ex{"id": metricID}).
		ToSQL()
	if err != nil {
		return fmt.Errorf("delete metric error during query building: %w", err)
	}

	if _, err := r.gr.ExecContext(ctx, qu); err != nil {
		return fmt.Errorf("delete metric error during execute query: %w", err)
	}

	return nil
}

func (r *MetricRepositoryImpl) PingDB() bool {
	err := r.gr.DB.Ping()
	if err != nil {
		return false
	}

	return true
}
