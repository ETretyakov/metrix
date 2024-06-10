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
	qu, _, err := goqu.
		Insert(metricTName).
		Rows(metric).
		Returning("id").
		ToSQL()
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

func (r *MetricRepositoryImpl) ReadIDs(ctx context.Context) (*[]string, error) {
	qu, _, err := goqu.
		Select("id").
		From(metricTName).
		ToSQL()
	if err != nil {
		return nil, fmt.Errorf("read ids metric error during query building: %w", err)
	}

	rows, err := r.gr.DB.QueryxContext(ctx, qu)
	if err != nil {
		return nil, fmt.Errorf("read metrics error during querying: %w", err)
	}

	var ids []string
	for rows.Next() {
		id := ""
		err := rows.Scan(&id)
		if err != nil {
			return nil, fmt.Errorf("read metrics error during scan rows: %w", err)
		}
		ids = append(ids, id)
	}

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &ids, nil
}

func (r *MetricRepositoryImpl) ReadMany(
	ctx context.Context,
	metricIDs []string,
) (*[]model.Metric, error) {
	inIDs := []any{}
	for _, id := range metricIDs {
		inIDs = append(inIDs, id)
	}

	qu, _, err := goqu.
		Select(&model.Metric{}).
		From(metricTName).
		Where(goqu.C("id").In(inIDs...)).
		ToSQL()
	if err != nil {
		return nil, fmt.Errorf("read metrics error during query building: %w", err)
	}

	rows, err := r.gr.DB.QueryxContext(ctx, qu)
	if err != nil {
		return nil, fmt.Errorf("read metrics error during querying: %w", err)
	}

	newMetrics := []model.Metric{}
	for rows.Next() {
		metric := model.Metric{}
		err := rows.StructScan(&metric)
		if err != nil {
			return nil, fmt.Errorf("read metrics error during scan rows: %w", err)
		}
		newMetrics = append(newMetrics, metric)
	}

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &newMetrics, nil
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

	tx, err := r.gr.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Commit()

	if _, err := tx.ExecContext(ctx, qu); err != nil {
		return nil, fmt.Errorf("update metric error during execute query: %w", err)
	}

	metricOut, err := r.Read(ctx, metric.ID)
	if err != nil {
		return nil, fmt.Errorf("refresh metric error: %w", err)
	}

	return metricOut, nil
}

func (r *MetricRepositoryImpl) UpsertMany(
	ctx context.Context,
	metrics []model.Metric,
) (bool, error) {
	tx, err := r.gr.BeginTxx(ctx, nil)
	if err != nil {
		return false, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Commit()

	for _, m := range metrics {
		metric, err := r.Read(ctx, m.ID)
		if err != nil {
			return false, fmt.Errorf("failed to get metric: %w", err)
		}

		qu := ""
		if metric == nil {
			qu, _, err = goqu.
				Insert(metricTName).
				Rows(m).
				Returning("id").
				ToSQL()
			if err != nil {
				return false, fmt.Errorf("create metric error during query building: %w", err)
			}
			if _, err := tx.ExecContext(ctx, qu); err != nil {
				return false, fmt.Errorf("failed to create metric: %w", err)
			}
		} else {
			qu, _, err = goqu.
				Update(metricTName).
				Set(m).
				Where(goqu.Ex{"id": metric.ID}).
				ToSQL()
			if err != nil {
				return false, fmt.Errorf("create metric error during query building: %w", err)
			}
			if _, err := tx.ExecContext(ctx, qu); err != nil {
				return false, fmt.Errorf("failed to update metric: %w", err)
			}
		}
	}

	return true, nil
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
	return err == nil
}
