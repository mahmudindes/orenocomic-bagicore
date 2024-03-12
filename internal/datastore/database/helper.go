package database

import (
	"context"
	"time"

	"github.com/mahmudindes/orenocomic-bagicore/internal/model"
)

func (db Database) GenericAdd(ctx context.Context, t string, data map[string]any, v any) error {
	cols, vals, args := SetInsert(data)
	sql := "INSERT INTO " + t + " (" + cols + ") VALUES (" + vals + ")"
	if v != nil {
		sql += ` RETURNING *`
		if err := db.QueryOne(ctx, v, sql, args...); err != nil {
			return err
		}
	} else {
		if err := db.Exec(ctx, sql, args...); err != nil {
			return err
		}
	}
	return nil
}

func (db Database) BatchAdd(ctx context.Context, t string, data []map[string]any, v any) error {
	cols, valx, args := SetBulkInsert(data)
	sql := "INSERT INTO " + t + " (" + cols + ") VALUES"
	for i, vals := range valx {
		if i > 0 {
			sql += ", "
		}
		sql += " (" + vals + ")"
	}
	if v != nil {
		sql += ` RETURNING *`
		if err := db.QueryAll(ctx, v, sql, args...); err != nil {
			return err
		}
	} else {
		if err := db.Exec(ctx, sql, args...); err != nil {
			return err
		}
	}
	return nil
}

func (db Database) GenericGet(ctx context.Context, t string, conds any, v any) error {
	args := []any{}
	cond := SetWhere(conds, &args)
	sql := "SELECT * FROM " + t + " WHERE " + cond
	return db.QueryOne(ctx, v, sql, args...)
}

func (db Database) GenericUpdate(ctx context.Context, t string, data map[string]any, conds any, v any) error {
	data[model.DBGenericUpdatedAt] = time.Now().UTC()
	sets, args := SetUpdate(data)
	cond := SetWhere([]any{conds, SetUpdateWhere(data)}, &args)
	sql := "UPDATE " + t + " SET " + sets + " WHERE " + cond
	if v != nil {
		sql += ` RETURNING *`
		if err := db.QueryOne(ctx, v, sql, args...); err != nil {
			return err
		}
	} else {
		if err := db.Exec(ctx, sql, args...); err != nil {
			return err
		}
	}
	return nil
}

func (db Database) GenericDelete(ctx context.Context, t string, conds any, v any) error {
	args := []any{}
	cond := SetWhere(conds, &args)
	sql := "DELETE FROM " + t + " WHERE " + cond
	if v != nil {
		sql += ` RETURNING *`
		if err := db.QueryOne(ctx, v, sql, args...); err != nil {
			return err
		}
	} else {
		if err := db.Exec(ctx, sql, args...); err != nil {
			return err
		}
	}
	return nil
}

func (db Database) GenericList(ctx context.Context, t string, params model.ListParams, v any) error {
	args := []any{}
	sql := "SELECT * FROM " + t
	if cond := SetWhere(params.Conditions, &args); cond != "" {
		sql += " WHERE " + cond
	}
	if odbs := SetOrderBys(params.OrderBys, &args); odbs != "" {
		sql += " ORDER BY " + odbs
	}
	if params.Pagination != nil {
		sql += SetPagination(*params.Pagination, &args)
	}
	return db.QueryAll(ctx, v, sql, args...)
}

func (db Database) GenericCount(ctx context.Context, t string, conds any) (int, error) {
	var dst int
	args := []any{}
	sql := "SELECT COUNT(*) FROM " + t
	if cond := SetWhere(conds, &args); cond != "" {
		sql += " WHERE " + cond
	}
	if err := db.QueryOne(ctx, &dst, sql, args...); err != nil {
		return -1, err
	}
	return dst, nil
}

func (db Database) GenericExists(ctx context.Context, t string, conds any) (bool, error) {
	var dst bool
	args := []any{}
	sql := "SELECT EXISTS(SELECT 1 FROM " + t
	if cond := SetWhere(conds, &args); cond != "" {
		sql += " WHERE " + cond
	}
	sql += ")"
	if err := db.QueryOne(ctx, &dst, sql, args...); err != nil {
		return false, err
	}
	return dst, nil
}
