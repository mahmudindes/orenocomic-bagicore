package database

import (
	"context"
	"errors"
	"time"

	"github.com/mahmudindes/orenocomic-bagicore/internal/model"
)

const (
	NameErrWebsiteKey             = "website_domain_key"
	NameErrWebsiteTLLanguagePKey  = "website_tllanguage_pkey"
	NameErrWebsiteTLLanguageFKey0 = "website_tllanguage_website_id_fkey"
	NameErrWebsiteTLLanguageFKey1 = "website_tllanguage_language_id_fkey"
)

func (db Database) AddWebsite(ctx context.Context, data model.AddWebsite, v *model.Website) error {
	if err := db.GenericAdd(ctx, model.DBWebsite, map[string]any{
		model.DBWebsiteDomain:    data.Domain,
		model.DBWebsiteName:      data.Name,
		model.DBWebsiteMachineTL: data.MachineTL,
	}, v); err != nil {
		return websiteSetError(err)
	}
	return nil
}

func (db Database) GetWebsite(ctx context.Context, conds any) (*model.Website, error) {
	var result model.Website
	if err := db.GenericGet(ctx, model.DBWebsite, conds, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (db Database) UpdateWebsite(ctx context.Context, data model.SetWebsite, conds any, v *model.Website) error {
	data0 := map[string]any{}
	if data.Domain != nil {
		data0[model.DBWebsiteDomain] = data.Domain
	}
	if data.Name != nil {
		data0[model.DBWebsiteName] = data.Name
	}
	if data.MachineTL != nil {
		data0[model.DBWebsiteMachineTL] = data.MachineTL
	}
	for _, null := range data.SetNull {
		data0[null] = nil
	}
	if err := db.GenericUpdate(ctx, model.DBWebsite, data0, conds, v); err != nil {
		return websiteSetError(err)
	}
	return nil
}

func (db Database) DeleteWebsite(ctx context.Context, conds any, v *model.Website) error {
	return db.GenericDelete(ctx, model.DBWebsite, conds, v)
}

func (db Database) ListWebsite(ctx context.Context, params model.ListParams) ([]*model.Website, error) {
	result := []*model.Website{}
	if len(params.OrderBys) < 1 {
		params.OrderBys = append(params.OrderBys, model.OrderBy{Field: model.DBWebsiteDomain})
	}
	if params.Pagination == nil {
		params.Pagination = &model.Pagination{Page: 1, Limit: model.WebsitePaginationDef}
	}
	if err := db.GenericList(ctx, model.DBWebsite, params, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (db Database) CountWebsite(ctx context.Context, conds any) (int, error) {
	return db.GenericCount(ctx, model.DBWebsite, conds)
}

func (db Database) ExistsWebsite(ctx context.Context, conds any) (bool, error) {
	return db.GenericExists(ctx, model.DBWebsite, conds)
}

func websiteSetError(err error) error {
	var errDatabase model.DatabaseError
	if errors.As(err, &errDatabase) {
		if errDatabase.Code == CodeErrExists && errDatabase.Name == NameErrWebsiteKey {
			return model.GenericError("same domain already exists")
		}
	}
	return err
}

func (db Database) AddWebsiteTLLanguage(ctx context.Context, data model.AddWebsiteTLLanguage, v *model.WebsiteTLLanguage) error {
	var websiteID any
	switch {
	case data.WebsiteID != nil:
		websiteID = data.WebsiteID
	case data.WebsiteDomain != nil:
		websiteID = model.DBWebsiteDomainToID(*data.WebsiteDomain)
	}
	var languageID any
	switch {
	case data.LanguageID != nil:
		languageID = data.LanguageID
	case data.LanguageIETF != nil:
		languageID = model.DBLanguageIETFToID(*data.LanguageIETF)
	}
	cols, vals, args := SetInsert(map[string]any{
		model.DBWebsiteGenericWebsiteID:   websiteID,
		model.DBLanguageGenericLanguageID: languageID,
	})
	sql := "INSERT INTO " + model.DBWebsiteTLLanguage + " (" + cols + ") VALUES (" + vals + ")"
	if v != nil {
		sql += " RETURNING *"
		sql = "WITH data AS (" + sql + ")"
		sql += " SELECT w." + model.DBGenericCreatedAt + ", w." + model.DBGenericUpdatedAt
		sql += ", w." + model.DBWebsiteGenericWebsiteID + ", w." + model.DBLanguageGenericLanguageID
		sql += ", l." + model.DBLanguageIETF + " AS language_ietf"
		sql += " FROM data w JOIN " + model.DBLanguage + " l"
		sql += " ON w." + model.DBLanguageGenericLanguageID + " = l." + model.DBGenericID
		if err := db.QueryOne(ctx, v, sql, args...); err != nil {
			return websiteTLLanguageSetError(err)
		}
	} else {
		if err := db.Exec(ctx, sql, args...); err != nil {
			return websiteTLLanguageSetError(err)
		}
	}
	return nil
}

func (db Database) GetWebsiteTLLanguage(ctx context.Context, conds any) (*model.WebsiteTLLanguage, error) {
	var result model.WebsiteTLLanguage
	args := []any{}
	cond := SetWhere(conds, &args)
	sql := "SELECT * FROM (SELECT w." + model.DBGenericCreatedAt + ", w." + model.DBGenericUpdatedAt
	sql += ", w." + model.DBWebsiteGenericWebsiteID + ", w." + model.DBLanguageGenericLanguageID
	sql += ", l." + model.DBLanguageIETF + " AS language_ietf"
	sql += " FROM " + model.DBWebsiteTLLanguage + " w JOIN " + model.DBLanguage + " l"
	sql += " ON w." + model.DBLanguageGenericLanguageID + " = l." + model.DBGenericID
	sql += ")"
	sql += " WHERE " + cond
	if err := db.QueryOne(ctx, &result, sql, args...); err != nil {
		return nil, err
	}
	return &result, nil
}

func (db Database) UpdateWebsiteTLLanguage(ctx context.Context, data model.SetWebsiteTLLanguage, conds any, v *model.WebsiteTLLanguage) error {
	data0 := map[string]any{}
	switch {
	case data.WebsiteID != nil:
		data0[model.DBWebsiteGenericWebsiteID] = data.WebsiteID
	case data.WebsiteDomain != nil:
		data0[model.DBWebsiteGenericWebsiteID] = model.DBWebsiteDomainToID(*data.WebsiteDomain)
	}
	switch {
	case data.LanguageID != nil:
		data0[model.DBLanguageGenericLanguageID] = data.LanguageID
	case data.LanguageIETF != nil:
		data0[model.DBLanguageGenericLanguageID] = model.DBLanguageIETFToID(*data.LanguageIETF)
	}
	data0[model.DBGenericUpdatedAt] = time.Now().UTC()
	sets, args := SetUpdate(data0)
	cond := SetWhere([]any{conds, SetUpdateWhere(data0)}, &args)
	sql := "UPDATE " + model.DBWebsiteTLLanguage + " SET " + sets + " WHERE " + cond
	if v != nil {
		sql += " RETURNING *"
		sql = "WITH data AS (" + sql + ")"
		sql += " SELECT w." + model.DBGenericCreatedAt + ", w." + model.DBGenericUpdatedAt
		sql += ", w." + model.DBWebsiteGenericWebsiteID + ", w." + model.DBLanguageGenericLanguageID
		sql += ", l." + model.DBLanguageIETF + " AS language_ietf"
		sql += " FROM data w JOIN " + model.DBLanguage + " l"
		sql += " ON w." + model.DBLanguageGenericLanguageID + " = l." + model.DBGenericID
		if err := db.QueryOne(ctx, v, sql, args...); err != nil {
			return websiteTLLanguageSetError(err)
		}
	} else {
		if err := db.Exec(ctx, sql, args...); err != nil {
			return websiteTLLanguageSetError(err)
		}
	}
	return nil
}

func (db Database) DeleteWebsiteTLLanguage(ctx context.Context, conds any, v *model.WebsiteTLLanguage) error {
	args := []any{}
	cond := SetWhere(conds, &args)
	sql := "DELETE FROM " + model.DBWebsiteTLLanguage + " WHERE " + cond
	if v != nil {
		sql += " RETURNING *"
		sql = "WITH data AS (" + sql + ")"
		sql += " SELECT w." + model.DBGenericCreatedAt + ", w." + model.DBGenericUpdatedAt
		sql += ", w." + model.DBWebsiteGenericWebsiteID + ", w." + model.DBLanguageGenericLanguageID
		sql += ", l." + model.DBLanguageIETF + " AS language_ietf"
		sql += " FROM data w JOIN " + model.DBLanguage + " l"
		sql += " ON w." + model.DBLanguageGenericLanguageID + " = l." + model.DBGenericID
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

func (db Database) ListWebsiteTLLanguage(ctx context.Context, params model.ListParams) ([]*model.WebsiteTLLanguage, error) {
	result := []*model.WebsiteTLLanguage{}
	args := []any{}
	sql := "SELECT * FROM (SELECT w." + model.DBGenericCreatedAt + ", w." + model.DBGenericUpdatedAt
	sql += ", w." + model.DBWebsiteGenericWebsiteID + ", w." + model.DBLanguageGenericLanguageID
	sql += ", l." + model.DBLanguageIETF + " AS language_ietf"
	sql += " FROM " + model.DBWebsiteTLLanguage + " w JOIN " + model.DBLanguage + " l"
	sql += " ON w." + model.DBLanguageGenericLanguageID + " = l." + model.DBGenericID
	sql += ")"
	if cond := SetWhere(params.Conditions, &args); cond != "" {
		sql += " WHERE " + cond
	}
	if len(params.OrderBys) < 1 {
		params.OrderBys = append(params.OrderBys, model.OrderBy{Field: model.DBLanguageGenericLanguageID})
	}
	sql += " ORDER BY " + SetOrderBys(params.OrderBys, &args)
	if params.Pagination == nil {
		params.Pagination = &model.Pagination{Page: 1, Limit: model.WebsiteTLLanguagePaginationDef}
	}
	if lmof := SetPagination(*params.Pagination, &args); lmof != "" {
		sql += lmof
	}
	if err := db.QueryAll(ctx, &result, sql, args...); err != nil {
		return nil, err
	}
	return result, nil
}

func (db Database) CountWebsiteTLLanguage(ctx context.Context, conds any) (int, error) {
	return db.GenericCount(ctx, model.DBWebsiteTLLanguage, conds)
}

func websiteTLLanguageSetError(err error) error {
	var errDatabase model.DatabaseError
	if errors.As(err, &errDatabase) {
		if errDatabase.Code == CodeErrForeign {
			switch errDatabase.Name {
			case NameErrWebsiteTLLanguageFKey0:
				return model.GenericError("website does not exist")
			case NameErrWebsiteTLLanguageFKey1:
				return model.GenericError("language does not exist")
			}
		}
		if errDatabase.Code == CodeErrExists && errDatabase.Name == NameErrWebsiteTLLanguagePKey {
			return model.GenericError("same language id already exists")
		}
	}
	return err
}
