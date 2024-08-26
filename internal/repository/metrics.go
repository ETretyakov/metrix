// Module "repositry" that holds the functionality that is related to database communication.
package repository

import (
	"context"
	"database/sql"
	"fmt"

	"metrix/internal/model"

	"github.com/doug-martin/goqu/v9"
	"github.com/pkg/errors"
)

const metricTName = "mtr_metrics"

// MetricRepositoryImpl - the structure for implementation of the MetricRepository concept.
type MetricRepositoryImpl struct {
	gr *Group
}

// NewMetricRepository - the builder function for MetricRepositoryImpl.
func NewMetricRepository(db *Group) *MetricRepositoryImpl {
	return &MetricRepositoryImpl{
		gr: db,
	}
}

// Create - the method to insert metric into database table.
func (r *MetricRepositoryImpl) Create(
	ctx context.Context,
	metric *model.Metric,
) (*model.Metric, error) {
	qu, _, err := goqu.
		Insert(metricTName).
		Rows(metric).
		Returning("id", "mtype", "delta", "value").
		ToSQL()
	if err != nil {
		return nil, errors.Wrapf(err, "create metric error during query building")
	}

	err = r.gr.DB.QueryRowxContext(ctx, qu).StructScan(metric)
	if err != nil {
		return nil, fmt.Errorf("failed to run insert: %w", err)
	}

	return metric, nil
}

// Read - the method to read metric record from database.
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
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("read metric error during scan row: %w", err)
	}

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	return &newMetric, nil
}

// ReadIDs - the method that retrieves metrics ids.
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

	if rows.Err() != nil {
		return nil, fmt.Errorf("read metrics error during querying: %w", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			return
		}
	}()

	return &ids, nil
}

// ReadMany - the method to read metrics in batch.
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

	if rows.Err() != nil {
		return nil, fmt.Errorf("read metrics error during querying: %w", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			return
		}
	}()

	return &newMetrics, nil
}

// Update - the method to update metric record.
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

	tx, err := r.gr.DB.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			return
		}
	}()

	if _, err := tx.ExecContext(ctx, qu); err != nil {
		return nil, fmt.Errorf("update metric error during execute query: %w", err)
	}

	metricOut, err := r.Read(ctx, metric.ID)
	if err != nil {
		return nil, errors.Wrapf(err, "refresh metric error")
	}

	if err := tx.Commit(); err != nil {
		return nil, errors.Wrapf(err, "failed to commit")
	}

	return metricOut, nil
}

// UpsertMany - the method to insert/update metric record in batch.
func (r *MetricRepositoryImpl) UpsertMany(
	ctx context.Context,
	metrics []model.Metric,
) (bool, error) {
	tx, err := r.gr.DB.BeginTxx(ctx, nil)
	if err != nil {
		return false, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if err := tx.Rollback(); err != nil {
			return
		}
	}()

	rowsIn := []any{}
	for _, m := range metrics {
		rowsIn = append(rowsIn, m)
	}
	qu, _, err := goqu.Insert(metricTName).Rows(rowsIn...).ToSQL()
	if err != nil {
		return false, fmt.Errorf("upsert metric error during query building: %w", err)
	}
	qu += " ON CONFLICT ON CONSTRAINT mtr_metrics_pk DO UPDATE SET delta = excluded.delta, value = excluded.value"

	if _, err := tx.ExecContext(ctx, qu); err != nil {
		return false, fmt.Errorf("failed to upsert metric: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return false, errors.Wrapf(err, "failed to commit")
	}

	return true, nil
}

// Delete - the method to remove records from the database.
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

	if _, err := r.gr.DB.ExecContext(ctx, qu); err != nil {
		return fmt.Errorf("delete metric error during execute query: %w", err)
	}

	return nil
}

// PingDB - the method to ping database connection.
func (r *MetricRepositoryImpl) PingDB(ctx context.Context) bool {
	err := r.gr.DB.PingDB(ctx)
	return err == nil
}
