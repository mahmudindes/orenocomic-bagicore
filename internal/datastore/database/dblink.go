package database

import (
	"context"
	"errors"
	"time"

	"github.com/mahmudindes/orenocomic-bagicore/internal/model"
)

const (
	NameErrLinkKey  = "link_website_id_relative_url_key"
	NameErrLinkFKey = "link_website_id_fkey"
)

func (db Database) AddLink(ctx context.Context, data model.AddLink, v *model.Link) error {
	var websiteID any
	switch {
	case data.WebsiteID != nil:
		websiteID = data.WebsiteID
	case data.WebsiteDomain != nil:
		websiteID = model.DBWebsiteDomainToID(*data.WebsiteDomain)
	}
	cols, vals, args := SetInsert(map[string]any{
		model.DBWebsiteGenericWebsiteID: websiteID,
		model.DBLinkRelativeURL:         data.RelativeURL,
		model.DBLinkMachineTL:           data.MachineTL,
	})
	sql := "INSERT INTO " + model.DBLink + " (" + cols + ") VALUES (" + vals + ")"
	if v != nil {
		sql += " RETURNING *"
		sql = "WITH data AS (" + sql + ")"
		sql += " SELECT w." + model.DBGenericID
		sql += ", w." + model.DBGenericCreatedAt + ", w." + model.DBGenericUpdatedAt
		sql += ", w." + model.DBWebsiteGenericWebsiteID + ", w." + model.DBLinkRelativeURL
		sql += ", w." + model.DBLinkMachineTL
		sql += ", l." + model.DBWebsiteDomain + " AS website_domain"
		sql += " FROM data w JOIN " + model.DBWebsite + " l"
		sql += " ON w." + model.DBWebsiteGenericWebsiteID + " = l." + model.DBGenericID
		if err := db.QueryOne(ctx, v, sql, args...); err != nil {
			return linkSetError(err)
		}
	} else {
		if err := db.Exec(ctx, sql, args...); err != nil {
			return linkSetError(err)
		}
	}
	return nil
}

func (db Database) GetLink(ctx context.Context, conds any) (*model.Link, error) {
	var result model.Link
	args := []any{}
	cond := SetWhere(conds, &args)
	sql := "SELECT * FROM (SELECT w." + model.DBGenericID
	sql += ", w." + model.DBGenericCreatedAt + ", w." + model.DBGenericUpdatedAt
	sql += ", w." + model.DBWebsiteGenericWebsiteID + ", w." + model.DBLinkRelativeURL
	sql += ", w." + model.DBLinkMachineTL
	sql += ", l." + model.DBWebsiteDomain + " AS website_domain"
	sql += " FROM " + model.DBLink + " w JOIN " + model.DBWebsite + " l"
	sql += " ON w." + model.DBWebsiteGenericWebsiteID + " = l." + model.DBGenericID
	sql += ")"
	sql += " WHERE " + cond
	if err := db.QueryOne(ctx, &result, sql, args...); err != nil {
		return nil, err
	}
	return &result, nil
}

func (db Database) UpdateLink(ctx context.Context, data model.SetLink, conds any, v *model.Link) error {
	data0 := map[string]any{}
	switch {
	case data.WebsiteID != nil:
		data0[model.DBWebsiteGenericWebsiteID] = data.WebsiteID
	case data.WebsiteDomain != nil:
		data0[model.DBWebsiteGenericWebsiteID] = model.DBWebsiteDomainToID(*data.WebsiteDomain)
	}
	if data.RelativeURL != nil {
		data0[model.DBLinkRelativeURL] = data.RelativeURL
	}
	if data.MachineTL != nil {
		data0[model.DBLinkMachineTL] = data.MachineTL
	}
	for _, null := range data.SetNull {
		data0[null] = nil
	}
	data0[model.DBGenericUpdatedAt] = time.Now().UTC()
	sets, args := SetUpdate(data0)
	cond := SetWhere([]any{conds, SetUpdateWhere(data0)}, &args)
	sql := "UPDATE " + model.DBLink + " SET " + sets + " WHERE " + cond
	if v != nil {
		sql += " RETURNING *"
		sql = "WITH data AS (" + sql + ")"
		sql += " SELECT w." + model.DBGenericID
		sql += ", w." + model.DBGenericCreatedAt + ", w." + model.DBGenericUpdatedAt
		sql += ", w." + model.DBWebsiteGenericWebsiteID + ", w." + model.DBLinkRelativeURL
		sql += ", w." + model.DBLinkMachineTL
		sql += ", l." + model.DBWebsiteDomain + " AS website_domain"
		sql += " FROM data w JOIN " + model.DBWebsite + " l"
		sql += " ON w." + model.DBWebsiteGenericWebsiteID + " = l." + model.DBGenericID
		if err := db.QueryOne(ctx, v, sql, args...); err != nil {
			return linkSetError(err)
		}
	} else {
		if err := db.Exec(ctx, sql, args...); err != nil {
			return linkSetError(err)
		}
	}
	return nil
}

func (db Database) DeleteLink(ctx context.Context, conds any, v *model.Link) error {
	args := []any{}
	cond := SetWhere(conds, &args)
	sql := "DELETE FROM " + model.DBLink + " WHERE " + cond
	if v != nil {
		sql += " RETURNING *"
		sql = "WITH data AS (" + sql + ")"
		sql += " SELECT w." + model.DBGenericID
		sql += ", w." + model.DBGenericCreatedAt + ", w." + model.DBGenericUpdatedAt
		sql += ", w." + model.DBWebsiteGenericWebsiteID + ", w." + model.DBLinkRelativeURL
		sql += ", w." + model.DBLinkMachineTL
		sql += ", l." + model.DBWebsiteDomain + " AS website_domain"
		sql += " FROM data w JOIN " + model.DBWebsite + " l"
		sql += " ON w." + model.DBWebsiteGenericWebsiteID + " = l." + model.DBGenericID
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

func (db Database) ListLink(ctx context.Context, params model.ListParams) ([]*model.Link, error) {
	result := []*model.Link{}
	args := []any{}
	sql := "SELECT * FROM (SELECT w." + model.DBGenericID
	sql += ", w." + model.DBGenericCreatedAt + ", w." + model.DBGenericUpdatedAt
	sql += ", w." + model.DBWebsiteGenericWebsiteID + ", w." + model.DBLinkRelativeURL
	sql += ", w." + model.DBLinkMachineTL
	sql += ", l." + model.DBWebsiteDomain + " AS website_domain"
	sql += " FROM " + model.DBLink + " w JOIN " + model.DBWebsite + " l"
	sql += " ON w." + model.DBWebsiteGenericWebsiteID + " = l." + model.DBGenericID
	sql += ")"
	if cond := SetWhere(params.Conditions, &args); cond != "" {
		sql += " WHERE " + cond
	}
	if len(params.OrderBys) < 1 {
		params.OrderBys = append(params.OrderBys, model.OrderBy{Field: model.DBGenericID})
	}
	sql += " ORDER BY " + SetOrderBys(params.OrderBys, &args)
	if params.Pagination == nil {
		params.Pagination = &model.Pagination{Page: 1, Limit: model.LinkPaginationDef}
	}
	if lmof := SetPagination(*params.Pagination, &args); lmof != "" {
		sql += lmof
	}
	if err := db.QueryAll(ctx, &result, sql, args...); err != nil {
		return nil, err
	}
	return result, nil
}

func (db Database) CountLink(ctx context.Context, conds any) (int, error) {
	return db.GenericCount(ctx, model.DBLink, conds)
}

func (db Database) ExistsLink(ctx context.Context, conds any) (bool, error) {
	return db.GenericExists(ctx, model.DBLink, conds)
}

func linkSetError(err error) error {
	var errDatabase model.DatabaseError
	if errors.As(err, &errDatabase) {
		if errDatabase.Code == CodeErrExists && errDatabase.Name == NameErrLinkKey {
			return model.GenericError("same website id + path already exists")
		}
	}
	return err
}

const (
	NameErrLinkTLLanguagePKey  = "link_tllanguage_pkey"
	NameErrLinkTLLanguageFKey0 = "link_tllanguage_link_id_fkey"
	NameErrLinkTLLanguageFKey1 = "link_tllanguage_language_id_fkey"
)

func (db Database) AddLinkTLLanguage(ctx context.Context, data model.AddLinkTLLanguage, v *model.LinkTLLanguage) error {
	var linkID any
	switch {
	case data.LinkID != nil:
		linkID = data.LinkID
	case data.LinkSID != nil:
		linkID = model.DBLinkSIDToID(*data.LinkSID)
	}
	var languageID any
	switch {
	case data.LanguageID != nil:
		languageID = data.LanguageID
	case data.LanguageIETF != nil:
		languageID = model.DBLanguageIETFToID(*data.LanguageIETF)
	}
	cols, vals, args := SetInsert(map[string]any{
		model.DBLinkGenericLinkID:         linkID,
		model.DBLanguageGenericLanguageID: languageID,
	})
	sql := "INSERT INTO " + model.DBLinkTLLanguage + " (" + cols + ") VALUES (" + vals + ")"
	if v != nil {
		sql += " RETURNING *"
		sql = "WITH data AS (" + sql + ")"
		sql += " SELECT w." + model.DBGenericCreatedAt + ", w." + model.DBGenericUpdatedAt
		sql += ", w." + model.DBLinkGenericLinkID + ", w." + model.DBLanguageGenericLanguageID
		sql += ", l." + model.DBLanguageIETF + " AS language_ietf"
		sql += " FROM data w JOIN " + model.DBLanguage + " l"
		sql += " ON w." + model.DBLanguageGenericLanguageID + " = l." + model.DBGenericID
		if err := db.QueryOne(ctx, v, sql, args...); err != nil {
			return linkTLLanguageSetError(err)
		}
	} else {
		if err := db.Exec(ctx, sql, args...); err != nil {
			return linkTLLanguageSetError(err)
		}
	}
	return nil
}

func (db Database) GetLinkTLLanguage(ctx context.Context, conds any) (*model.LinkTLLanguage, error) {
	var result model.LinkTLLanguage
	args := []any{}
	cond := SetWhere(conds, &args)
	sql := "SELECT * FROM (SELECT w." + model.DBGenericCreatedAt + ", w." + model.DBGenericUpdatedAt
	sql += ", w." + model.DBLinkGenericLinkID + ", w." + model.DBLanguageGenericLanguageID
	sql += ", l." + model.DBLanguageIETF + " AS language_ietf"
	sql += " FROM " + model.DBLinkTLLanguage + " w JOIN " + model.DBLanguage + " l"
	sql += " ON w." + model.DBLanguageGenericLanguageID + " = l." + model.DBGenericID
	sql += ")"
	sql += " WHERE " + cond
	if err := db.QueryOne(ctx, &result, sql, args...); err != nil {
		return nil, err
	}
	return &result, nil
}

func (db Database) UpdateLinkTLLanguage(ctx context.Context, data model.SetLinkTLLanguage, conds any, v *model.LinkTLLanguage) error {
	data0 := map[string]any{}
	switch {
	case data.LinkID != nil:
		data0[model.DBLinkGenericLinkID] = data.LinkID
	case data.LinkSID != nil:
		data0[model.DBLinkGenericLinkID] = model.DBLinkSIDToID(*data.LinkSID)
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
	sql := "UPDATE " + model.DBLinkTLLanguage + " SET " + sets + " WHERE " + cond
	if v != nil {
		sql += " RETURNING *"
		sql = "WITH data AS (" + sql + ")"
		sql += " SELECT w." + model.DBGenericCreatedAt + ", w." + model.DBGenericUpdatedAt
		sql += ", w." + model.DBLinkGenericLinkID + ", w." + model.DBLanguageGenericLanguageID
		sql += ", l." + model.DBLanguageIETF + " AS language_ietf"
		sql += " FROM data w JOIN " + model.DBLanguage + " l"
		sql += " ON w." + model.DBLanguageGenericLanguageID + " = l." + model.DBGenericID
		if err := db.QueryOne(ctx, v, sql, args...); err != nil {
			return linkTLLanguageSetError(err)
		}
	} else {
		if err := db.Exec(ctx, sql, args...); err != nil {
			return linkTLLanguageSetError(err)
		}
	}
	return nil
}

func (db Database) DeleteLinkTLLanguage(ctx context.Context, conds any, v *model.LinkTLLanguage) error {
	args := []any{}
	cond := SetWhere(conds, &args)
	sql := "DELETE FROM " + model.DBLinkTLLanguage + " WHERE " + cond
	if v != nil {
		sql += " RETURNING *"
		sql = "WITH data AS (" + sql + ")"
		sql += " SELECT w." + model.DBGenericCreatedAt + ", w." + model.DBGenericUpdatedAt
		sql += ", w." + model.DBLinkGenericLinkID + ", w." + model.DBLanguageGenericLanguageID
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

func (db Database) ListLinkTLLanguage(ctx context.Context, params model.ListParams) ([]*model.LinkTLLanguage, error) {
	result := []*model.LinkTLLanguage{}
	args := []any{}
	sql := "SELECT * FROM (SELECT w." + model.DBGenericCreatedAt + ", w." + model.DBGenericUpdatedAt
	sql += ", w." + model.DBLinkGenericLinkID + ", w." + model.DBLanguageGenericLanguageID
	sql += ", l." + model.DBLanguageIETF + " AS language_ietf"
	sql += " FROM " + model.DBLinkTLLanguage + " w JOIN " + model.DBLanguage + " l"
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
		params.Pagination = &model.Pagination{Page: 1, Limit: model.LinkTLLanguagePaginationDef}
	}
	if lmof := SetPagination(*params.Pagination, &args); lmof != "" {
		sql += lmof
	}
	if err := db.QueryAll(ctx, &result, sql, args...); err != nil {
		return nil, err
	}
	return result, nil
}

func (db Database) CountLinkTLLanguage(ctx context.Context, conds any) (int, error) {
	return db.GenericCount(ctx, model.DBLinkTLLanguage, conds)
}

func linkTLLanguageSetError(err error) error {
	var errDatabase model.DatabaseError
	if errors.As(err, &errDatabase) {
		if errDatabase.Code == CodeErrForeign {
			switch errDatabase.Name {
			case NameErrLinkTLLanguageFKey0:
				return model.GenericError("link does not exist")
			case NameErrLinkTLLanguageFKey1:
				return model.GenericError("language does not exist")
			}
		}
		if errDatabase.Code == CodeErrExists && errDatabase.Name == NameErrLinkTLLanguagePKey {
			return model.GenericError("same language id already exists")
		}
	}
	return err
}
