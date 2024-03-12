package database

import (
	"slices"
	"strconv"
	"strings"

	"github.com/mahmudindes/orenocomic-bagicore/internal/model"
	"github.com/mahmudindes/orenocomic-bagicore/internal/utila"
)

func SetValue(val any, args *[]any) string {
	switch v := val.(type) {
	case model.DBQueryValue:
		cond := SetWhere(v.Conditions, args)
		*args = append(*args, v.ZeroValue)
		subs := "SELECT " + v.Expression + " FROM " + v.Table + " WHERE " + cond
		return "(SELECT COALESCE((" + subs + "), $" + strconv.Itoa(len(*args)) + "))"
	default:
		*args = append(*args, v)
		return "$" + strconv.Itoa(len(*args))
	}
}

func SetInsert(data map[string]any) (cols string, vals string, args []any) {
	for key, val := range data {
		if utila.NilData(val) {
			continue
		}
		length := len(args)
		if length > 0 {
			cols += ", "
			vals += ", "
		}
		cols += key
		vals += SetValue(val, &args)
	}
	return
}

func SetBulkInsert(data []map[string]any) (cols string, valx []string, args []any) {
	colvals := map[string][]any{}
	for i, data := range data {
		keys := []string{}
		for key, val := range data {
			keys = append(keys, key)
			if utila.NilData(val) {
				colvals[key] = append(colvals[key], nil)
				continue
			}
			if n := len(colvals[key]); n != i {
				for n < i {
					colvals[key] = append(colvals[key], nil)
					n++
				}
			}
			colvals[key] = append(colvals[key], val)
		}
		for key := range colvals {
			if slices.Contains(keys, key) {
				continue
			}
			colvals[key] = append(colvals[key], nil)
		}
	}
	for key, cval := range colvals {
		if cols != "" {
			cols += ", "
		}
		cols += key
		for i, val := range cval {
			sarg := "DEFAULT"
			if val != nil {
				sarg = SetValue(val, &args)
			}
			if len(valx) < i+1 {
				valx = append(valx, sarg)
			} else {
				valx[i] += ", " + sarg
			}
		}
	}
	return
}

func SetUpdate(data map[string]any) (sets string, args []any) {
	for key, val := range data {
		if sets != "" {
			sets += ", "
		}
		if val == nil {
			sets += key + " = NULL"
			continue
		}
		sets += key + " = " + SetValue(val, &args)
	}
	return
}

func SetUpdateWhere(data map[string]any) (cond map[string]any) {
	cond = make(map[string]any)
	for key, val := range data {
		if utila.NilData(val) {
			cond[key] = model.DBIsNotNull{}
			continue
		}
		cond[key] = model.DBIsNotDistinctFrom{Value: val}
	}
	return
}

func SetWhere(conds any, args *[]any) (cond string) {
	switch conds := conds.(type) {
	case []any:
		lop := "OR"
		for _, conds := range conds {
			switch conds.(type) {
			case model.DBLogicalAND:
				lop = "AND"
			case model.DBLogicalOR:
				lop = "OR"
			default:
				conx := SetWhere(conds, args)
				if conx == "" {
					continue
				}
				if cond != "" {
					cond += " " + lop + " "
				}
				cond += conx
			}
		}
	case map[string]any:
		for key, val := range conds {
			if utila.NilData(val) {
				continue
			}
			conx := SetWhere(model.DBConditionalKV{Key: key, Value: val}, args)
			if conx == "" {
				continue
			}
			if cond != "" {
				cond += " AND "
			}
			cond += conx
		}
	case model.DBConditionalKV:
		switch val := conds.Value.(type) {
		case model.DBIsDistinctFrom:
			*args = append(*args, val.Value)
			cond += conds.Key + " IS DISTINCT FROM $" + strconv.Itoa(len(*args))
		case model.DBIsNotDistinctFrom:
			*args = append(*args, val.Value)
			cond += conds.Key + " IS NOT DISTINCT FROM $" + strconv.Itoa(len(*args))
		case model.DBIsNull:
			cond += conds.Key + " IS NULL"
		case model.DBIsNotNull:
			cond += conds.Key + " IS NOT NULL"
		case model.DBBooleanIs:
			cond += conds.Key + " IS " + strconv.FormatBool(bool(val))
		case model.DBBooleanIsNot:
			cond += conds.Key + " IS NOT " + strconv.FormatBool(bool(val))
		case model.DBInsensitiveLike:
			*args = append(*args, string(val))
			cond += conds.Key + " ILIKE $" + strconv.Itoa(len(*args))
		default:
			cond += conds.Key + " = " + SetValue(val, args)
		}
	}
	return
}

func SetOrderBy(m model.OrderBy, args *[]any) (ob string) {
	if m.Field == "" {
		return
	}

	switch field := m.Field.(type) {
	case string:
		ob += field
	default:
		return
	}

	if m.Sort != "" {
		switch strings.ToLower(m.Sort) {
		case "a", "asc", "ascend", "ascending":
			ob += " ASC"
		case "d", "desc", "descend", "descending":
			ob += " DESC"
		}
	}

	if m.Null != "" {
		switch strings.ToLower(m.Null) {
		case "f", "first":
			ob += " NULLS FIRST"
		case "l", "last":
			ob += " NULLS LAST"
		}
	}

	return
}

func SetOrderBys(m model.OrderBys, args *[]any) (obs string) {
	for _, ob := range m {
		if obs != "" {
			obs += ", "
		}

		obs += SetOrderBy(ob, args)
	}
	return
}

func SetPagination(m model.Pagination, args *[]any) (lo string) {
	if m.Limit < 1 {
		return
	}

	*args = append(*args, m.Limit)
	lo += " LIMIT $" + strconv.Itoa(len(*args))

	offset := m.Limit * (m.Page - 1)
	if offset > 0 {
		*args = append(*args, offset)
		lo += " OFFSET $" + strconv.Itoa(len(*args))
	}

	return
}
