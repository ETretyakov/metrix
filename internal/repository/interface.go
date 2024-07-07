package repository

import (
	"context"

	"metrix/internal/model"
	"metrix/internal/storages"
)

type FactoryExecutor interface {
	MetricRepository() MetricRepository
}

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

type Group struct {
	DB *storages.SQLDB

	MetricRepo MetricRepository
}

func NewGroup(
	ctx context.Context,
	db *storages.SQLDB,
	filePath string,
	storeInterval int64,
	restore bool,
) *Group {
	group := &Group{}

	if db != nil {
		group.DB = db
		group.MetricRepo = NewMetricRepository(group)
	} else {
		group.MetricRepo = storages.NewInMemmoryStorage(ctx, filePath, storeInterval, restore)
	}

	return group
}
