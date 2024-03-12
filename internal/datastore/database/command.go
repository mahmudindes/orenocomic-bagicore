package database

import (
	"context"
	"errors"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/mahmudindes/orenocomic-bagicore/internal/model"
)

const (
	CodeErrForeign    = "23503"
	CodeErrExists     = "23505"
	CodeErrValidation = "23514"
)

func (db Database) Exec(ctx context.Context, sql string, args ...any) error {
	_, err := db.client.Exec(ctx, sql, args...)

	return databaseError(err)
}

func (db Database) QueryAll(ctx context.Context, dst any, sql string, args ...any) error {
	rows, err := db.client.Query(ctx, sql, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	if err := pgxscan.ScanAll(dst, rows); err != nil {
		if pgxscan.NotFound(err) {
			return model.NotFoundError(err)
		}

		return databaseError(err)
	}

	return nil
}

func (db Database) QueryOne(ctx context.Context, dst any, sql string, args ...any) error {
	rows, err := db.client.Query(ctx, sql, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	if err := pgxscan.ScanOne(dst, rows); err != nil {
		if pgxscan.NotFound(err) {
			return model.NotFoundError(err)
		}

		return databaseError(err)
	}

	return nil
}

func databaseError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case CodeErrValidation:
			return model.WrappedError(model.DatabaseError{
				Name: pgErr.ConstraintName,
				Code: pgErr.Code,
				Err:  err,
			}, "database validation failed")
		case CodeErrForeign, CodeErrExists:
			return model.DatabaseError{Name: pgErr.ConstraintName, Code: pgErr.Code, Err: err}
		default:
			return model.DatabaseError{Code: pgErr.Code, Err: err}
		}
	}
	return err
}
