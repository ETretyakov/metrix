package repository

import (
	"context"

	"metrix/internal/model"
	"metrix/internal/storages"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// MetricRepository - the interface that describes all metric repository methods.
type MetricRepository interface {
	Create(ctx context.Context, metric *model.Metric) (*model.Metric, error)
	Read(ctx context.Context, metricID string) (*model.Metric, error)
	ReadIDs(ctx context.Context) (*[]string, error)
	ReadMany(ctx context.Context, metricIDs []string) (*[]model.Metric, error)
	Update(ctx context.Context, metric *model.Metric) (*model.Metric, error)
	UpsertMany(ctx context.Context, metrics []model.Metric) (bool, error)
	Delete(ctx context.Context, metricID string) error
	PingDB(ctx context.Context) bool
}

// Group - the structure that stores all necessary repositories.
type Group struct {
	DB *sqlx.DB

	MetricRepo MetricRepository
}

// PingDB - the method to pind database.
func (r *Group) PingDB(ctx context.Context) error {
	if err := r.DB.PingContext(ctx); err != nil {
		return errors.Wrapf(err, "failed to ping db")
	}
	return nil
}

// NewGroup - the builder function for Group structure.
func NewGroup(
	ctx context.Context,
	db *sqlx.DB,
	filePath string,
	storeInterval int64,
	restore bool,
) *Group {
	group := &Group{}

	if db != nil {
		group.DB = db
		group.MetricRepo = NewMetricRepository(group)
	} else {
		group.MetricRepo = storages.NewInMemoryStorage(ctx, filePath, storeInterval, restore)
	}

	return group
}
