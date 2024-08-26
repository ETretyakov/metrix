package storages

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"io"
	"time"

	"metrix/pkg/logger"

	"github.com/jackc/pgx"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

var recoverableErrors = []error{
	pgx.ErrDeadConn,
	sql.ErrConnDone,
	pgx.ErrConnBusy,
	driver.ErrBadConn,
	io.EOF,
}

const (
	retriesNumber = 3
	retryDuration = 5 * time.Second
	retryFactor   = .5
)

// Retryer - the structure that holds neccesarry things for database retryer.
type Retryer struct {
	Strategy Strategy
	OnRetry  func(ctx context.Context, n int, err error)
}

// SQLDB - the structure that combines database and retryer functionality.
type SQLDB struct {
	db      *sqlx.DB
	retryer *Retryer
}

// NewSQLDB - the builder function for SQLDB.
func NewSQLDB(db *sqlx.DB) *SQLDB {
	if db == nil {
		return nil
	}

	r := &SQLDB{
		db: db,
	}

	r.retryer = &Retryer{
		Strategy: Backoff(retriesNumber, retryDuration, retryFactor, true),
		OnRetry: func(ctx context.Context, n int, err error) {
			logger.Info(ctx, fmt.Sprintf("reconnecting DB (%d): %s", n, err))
			if err = r.db.PingContext(ctx); err != nil {
				logger.Error(ctx, "failed to connect to DB", err)
			}
		},
	}

	return r
}

func (r *SQLDB) retry(ctx context.Context, fn func() error) error {
	return r.retryer.Do(ctx, fn, recoverableErrors...)
}

// BeginTxx - the method to get a transaction for SQLDB.
func (r *SQLDB) BeginTxx(ctx context.Context, opts *sql.TxOptions) (tx *sqlx.Tx, err error) {
	err = r.retry(ctx, func() error {
		var err error
		tx, err = r.db.BeginTxx(ctx, opts)
		return errors.Wrapf(err, "failed to begin transaction")
	})
	return tx, err
}

// QueryxContext - the method to query database with context.
func (r *SQLDB) QueryxContext(
	ctx context.Context,
	query string,
	args ...any,
) (*sqlx.Rows, error) {
	var rows *sqlx.Rows
	err := r.retry(ctx, func() error {
		var err error
		rows, err = r.db.QueryxContext(ctx, query, args...)
		if err != nil {
			return errors.Wrapf(err, "failed to query context")
		}
		defer func() {
			if err := rows.Close(); err != nil {
				return
			}
		}()
		err = rows.Err()
		if err != nil {
			return errors.Wrapf(err, "failed to query context error in rows")
		}
		return nil
	})
	return rows, err
}

// QueryRowxContext - the method to query rows with context.
func (r *SQLDB) QueryRowxContext(
	ctx context.Context,
	query string,
	args ...any,
) (row *sqlx.Row) {
	_ = r.retry(ctx, func() error {
		row = r.db.QueryRowxContext(ctx, query, args...)
		return row.Err()
	})
	return
}

// ExecContext - the method to query execute queries with context.
func (r *SQLDB) ExecContext(
	ctx context.Context,
	query string,
	args ...any,
) (sql.Result, error) {
	var res sql.Result
	err := r.retry(ctx, func() error {
		var err error
		res, err = r.db.ExecContext(ctx, query, args...)
		return errors.Wrapf(err, "failed to retry")
	})
	return res, err
}

// PingDB - the method to pind database.
func (r *SQLDB) PingDB(ctx context.Context) error {
	if err := r.db.PingContext(ctx); err != nil {
		return errors.Wrapf(err, "failed to ping db")
	}
	return nil
}
