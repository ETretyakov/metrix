// Module "bootstrap" aggragtes all the operations related to service bootstrap.
package bootstrap

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"

	"metrix/internal/closer"
	"metrix/internal/config"
	"metrix/pkg/logger"

	_ "github.com/jackc/pgx/stdlib"
)

// Migrate - function to perform database migration operations.
func Migrate(ctx context.Context, cfg *config.Postgres) {
	db, err := goose.OpenDBWithDriver("pgx", cfg.DSN)
	if err != nil {
		logger.Fatal(ctx, "failed to open DB for migration", err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			logger.Fatal(ctx, "failed to close DB after migration", err)
		}
	}()

	if err := goose.RunContext(
		ctx,
		"up",
		db,
		cfg.MigrationFolder,
	); err != nil {
		logger.Fatal(ctx, "failed to run up DB migration", err)
	}
}

// InitDB - function that performs initiation of the service database.
func InitDB(ctx context.Context, cfg *config.Postgres) (*sqlx.DB, error) {
	logger.Info(ctx, "migrating db")
	Migrate(ctx, cfg)

	logger.Info(ctx, "initializing db")
	db, err := sqlx.Open("pgx", cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("make sql db instance error: %w", err)
	}
	logger.Info(ctx, "db initialized")

	closer.Add(db.Close)

	db.SetConnMaxLifetime(0)
	db.SetConnMaxIdleTime(0)
	db.SetMaxOpenConns(cfg.MaxOpenConn)
	db.SetMaxIdleConns(cfg.IdleConn)

	go func() {
		ticker := time.NewTicker(cfg.PingInterval)
		for {
			select {
			case <-ticker.C:
				if err := db.Ping(); err != nil {
					logger.Error(ctx, "error ping db", err)
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return db, nil
}
