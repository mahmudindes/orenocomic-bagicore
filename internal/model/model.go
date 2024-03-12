package model

import (
	"strings"

	"github.com/mahmudindes/orenocomic-bagicore/internal/utila"
)

type (
	DBLogicalAND        struct{}
	DBLogicalOR         struct{}
	DBIsDistinctFrom    struct{ Value any }
	DBIsNotDistinctFrom struct{ Value any }
	DBIsNull            struct{}
	DBIsNotNull         struct{}
	DBBooleanIs         bool
	DBBooleanIsNot      bool
	DBInsensitiveLike   string

	DBConditionalKV struct {
		Key   string
		Value any
	}
	DBQueryValue struct {
		Expression, Table     string
		ZeroValue, Conditions any
	}
	DBCrossConditional struct {
		Table      string
		Conditions any
	}
)

const (
	SecretString       = "*****"
	DBGenericID        = "id"
	DBGenericCreatedAt = "created_at"
	DBGenericUpdatedAt = "updated_at"
)

var (
	EmptyString         = ""
	GenericOrderByAllow = []string{DBGenericID, DBGenericCreatedAt, DBGenericUpdatedAt}
)

type ListParams struct {
	Conditions any
	OrderBys   OrderBys
	Pagination *Pagination
}

func (m ListParams) Validate() error {
	if err := m.OrderBys.Validate(); err != nil {
		return err
	}

	if m.Pagination != nil {
		if err := m.Pagination.Validate(); err != nil {
			return err
		}
	}

	return nil
}

type OrderBy struct {
	Field any
	Sort  string
	Null  string
}

func (ob OrderBy) Validate() error {
	if ob.Field == nil || ob.Field == "" {
		return GenericError("order by field must exist and cannot be empty")
	}

	if ob.Sort != "" {
		switch strings.ToLower(ob.Sort) {
		case "a", "asc", "ascend", "ascending":
			// Noop
		case "d", "desc", "descend", "descending":
			// Noop
		default:
			return GenericError("order by sort must be ascending or descending")
		}
	}

	if ob.Null != "" {
		switch strings.ToLower(ob.Null) {
		case "f", "first":
			// Noop
		case "l", "last":
			// Noop
		default:
			return GenericError("order by empty must be first or last")
		}
	}

	return nil
}

type OrderBys []OrderBy

func (obs OrderBys) Validate() error {
	for i, ob := range obs {
		if err := ob.Validate(); err != nil {
			return GenericError(utila.OrdinalNumber(i) + " " + err.Error())
		}
	}

	return nil
}

type Pagination struct {
	Page  int
	Limit int
}

func (p Pagination) Validate() error {
	if p.Page < 1 {
		return GenericError("pagination page must be at least 1")
	}

	if p.Limit < 1 {
		return GenericError("pagination limit must be at least 1")
	}

	return nil
}
