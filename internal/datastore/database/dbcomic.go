package database

import (
	"context"
	"errors"
	"time"

	"github.com/mahmudindes/orenocomic-bagicore/internal/model"
)

const (
	NameErrComicKey = "comic_code_key"
)

func (db Database) AddComic(ctx context.Context, data model.AddComic, v *model.Comic) error {
	if err := db.GenericAdd(ctx, model.DBComic, map[string]any{
		model.DBComicCode: data.Code,
	}, v); err != nil {
		return comicSetError(err)
	}
	return nil
}

func (db Database) GetComic(ctx context.Context, conds any) (*model.Comic, error) {
	var result model.Comic
	if err := db.GenericGet(ctx, model.DBComic, conds, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (db Database) UpdateComic(ctx context.Context, data model.SetComic, conds any, v *model.Comic) error {
	data0 := map[string]any{}
	if data.Code != nil {
		data0[model.DBComicCode] = data.Code
	}
	if err := db.GenericUpdate(ctx, model.DBComic, data0, conds, v); err != nil {
		return comicSetError(err)
	}
	return nil
}

func (db Database) DeleteComic(ctx context.Context, conds any, v *model.Comic) error {
	return db.GenericDelete(ctx, model.DBComic, conds, v)
}

func (db Database) ListComic(ctx context.Context, params model.ListParams) ([]*model.Comic, error) {
	result := []*model.Comic{}
	if len(params.OrderBys) < 1 {
		params.OrderBys = append(params.OrderBys, model.OrderBy{Field: model.DBComicCode})
	}
	if params.Pagination == nil {
		params.Pagination = &model.Pagination{Page: 1, Limit: model.ComicPaginationDef}
	}
	if err := db.GenericList(ctx, model.DBComic, params, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (db Database) CountComic(ctx context.Context, conds any) (int, error) {
	return db.GenericCount(ctx, model.DBComic, conds)
}

func (db Database) ExistsComic(ctx context.Context, conds any) (bool, error) {
	return db.GenericExists(ctx, model.DBComic, conds)
}

func comicSetError(err error) error {
	var errDatabase model.DatabaseError
	if errors.As(err, &errDatabase) {
		if errDatabase.Code == CodeErrExists && errDatabase.Name == NameErrComicKey {
			return model.GenericError("same code already exists")
		}
	}
	return err
}

const (
	NameErrComicLinkPKey  = "comic_link_pkey"
	NameErrComicLinkFKey0 = "comic_link_comic_id_fkey"
	NameErrComicLinkFKey1 = "comic_link_link_id_fkey"
)

func (db Database) AddComicLink(ctx context.Context, data model.AddComicLink, v *model.ComicLink) error {
	var comicID any
	switch {
	case data.ComicID != nil:
		comicID = data.ComicID
	case data.ComicCode != nil:
		comicID = model.DBComicCodeToID(*data.ComicCode)
	}
	var linkID any
	switch {
	case data.LinkID != nil:
		linkID = data.LinkID
	case data.LinkSID != nil:
		linkID = model.DBLinkSIDToID(*data.LinkSID)
	}
	cols, vals, args := SetInsert(map[string]any{
		model.DBComicGenericComicID: comicID,
		model.DBLinkGenericLinkID:   linkID,
	})
	sql := "INSERT INTO " + model.DBComicLink + " (" + cols + ") VALUES (" + vals + ")"
	if v != nil {
		sql += " RETURNING *"
		sql = "WITH data AS (" + sql + ")"
		sql += " SELECT a." + model.DBGenericCreatedAt + ", a." + model.DBGenericUpdatedAt
		sql += ", a." + model.DBComicGenericComicID + ", a." + model.DBLinkGenericLinkID
		sql += ", b." + model.DBLinkRelativeURL + " AS link_relative_url"
		sql += ", c." + model.DBWebsiteDomain + " AS link_website_domain"
		sql += " FROM data a JOIN " + model.DBLink + " b"
		sql += " ON a." + model.DBLinkGenericLinkID + " = b." + model.DBGenericID
		sql += " JOIN " + model.DBWebsite + " c"
		sql += " ON b." + model.DBWebsiteGenericWebsiteID + " = c." + model.DBGenericID
		if err := db.QueryOne(ctx, v, sql, args...); err != nil {
			return comicLinkSetError(err)
		}
	} else {
		if err := db.Exec(ctx, sql, args...); err != nil {
			return comicLinkSetError(err)
		}
	}
	return nil
}

func (db Database) GetComicLink(ctx context.Context, conds any) (*model.ComicLink, error) {
	var result model.ComicLink
	args := []any{}
	cond := SetWhere(conds, &args)
	sql := "SELECT * FROM (SELECT a." + model.DBGenericCreatedAt + ", a." + model.DBGenericUpdatedAt
	sql += ", a." + model.DBComicGenericComicID + ", a." + model.DBLinkGenericLinkID
	sql += ", b." + model.DBLinkRelativeURL + " AS link_relative_url"
	sql += ", c." + model.DBWebsiteDomain + " AS link_website_domain"
	sql += " FROM " + model.DBComicLink + " a JOIN " + model.DBLink + " b"
	sql += " ON a." + model.DBLinkGenericLinkID + " = b." + model.DBGenericID
	sql += " JOIN " + model.DBWebsite + " c"
	sql += " ON b." + model.DBWebsiteGenericWebsiteID + " = c." + model.DBGenericID
	sql += ")"
	sql += " WHERE " + cond
	if err := db.QueryOne(ctx, &result, sql, args...); err != nil {
		return nil, err
	}
	return &result, nil
}

func (db Database) UpdateComicLink(ctx context.Context, data model.SetComicLink, conds any, v *model.ComicLink) error {
	data0 := map[string]any{}
	switch {
	case data.ComicID != nil:
		data0[model.DBComicGenericComicID] = data.ComicID
	case data.ComicCode != nil:
		data0[model.DBComicGenericComicID] = model.DBComicCodeToID(*data.ComicCode)
	}
	switch {
	case data.LinkID != nil:
		data0[model.DBLinkGenericLinkID] = data.LinkID
	case data.LinkSID != nil:
		data0[model.DBLinkGenericLinkID] = model.DBLinkSIDToID(*data.LinkSID)
	}
	data0[model.DBGenericUpdatedAt] = time.Now().UTC()
	sets, args := SetUpdate(data0)
	cond := SetWhere([]any{conds, SetUpdateWhere(data0)}, &args)
	sql := "UPDATE " + model.DBComicLink + " SET " + sets + " WHERE " + cond
	if v != nil {
		sql += " RETURNING *"
		sql = "WITH data AS (" + sql + ")"
		sql += " SELECT a." + model.DBGenericCreatedAt + ", a." + model.DBGenericUpdatedAt
		sql += ", a." + model.DBComicGenericComicID + ", a." + model.DBLinkGenericLinkID
		sql += ", b." + model.DBLinkRelativeURL + " AS link_relative_url"
		sql += ", c." + model.DBWebsiteDomain + " AS link_website_domain"
		sql += " FROM data a JOIN " + model.DBLink + " b"
		sql += " ON a." + model.DBLinkGenericLinkID + " = b." + model.DBGenericID
		sql += " JOIN " + model.DBWebsite + " c"
		sql += " ON b." + model.DBWebsiteGenericWebsiteID + " = c." + model.DBGenericID
		if err := db.QueryOne(ctx, v, sql, args...); err != nil {
			return comicLinkSetError(err)
		}
	} else {
		if err := db.Exec(ctx, sql, args...); err != nil {
			return comicLinkSetError(err)
		}
	}
	return nil
}

func (db Database) DeleteComicLink(ctx context.Context, conds any, v *model.ComicLink) error {
	args := []any{}
	cond := SetWhere(conds, &args)
	sql := "DELETE FROM " + model.DBComicLink + " WHERE " + cond
	if v != nil {
		sql += " RETURNING *"
		sql = "WITH data AS (" + sql + ")"
		sql += " SELECT a." + model.DBGenericCreatedAt + ", a." + model.DBGenericUpdatedAt
		sql += ", a." + model.DBComicGenericComicID + ", a." + model.DBLinkGenericLinkID
		sql += ", b." + model.DBLinkRelativeURL + " AS link_relative_url"
		sql += ", c." + model.DBWebsiteDomain + " AS link_website_domain"
		sql += " FROM data a JOIN " + model.DBLink + " b"
		sql += " ON a." + model.DBLinkGenericLinkID + " = b." + model.DBGenericID
		sql += " JOIN " + model.DBWebsite + " c"
		sql += " ON b." + model.DBWebsiteGenericWebsiteID + " = c." + model.DBGenericID
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

func (db Database) ListComicLink(ctx context.Context, params model.ListParams) ([]*model.ComicLink, error) {
	result := []*model.ComicLink{}
	args := []any{}
	sql := "SELECT * FROM (SELECT a." + model.DBGenericCreatedAt + ", a." + model.DBGenericUpdatedAt
	sql += ", a." + model.DBComicGenericComicID + ", a." + model.DBLinkGenericLinkID
	sql += ", b." + model.DBLinkRelativeURL + " AS link_relative_url"
	sql += ", c." + model.DBWebsiteDomain + " AS link_website_domain"
	sql += " FROM " + model.DBComicLink + " a JOIN " + model.DBLink + " b"
	sql += " ON a." + model.DBLinkGenericLinkID + " = b." + model.DBGenericID
	sql += " JOIN " + model.DBWebsite + " c"
	sql += " ON b." + model.DBWebsiteGenericWebsiteID + " = c." + model.DBGenericID
	sql += ")"
	if cond := SetWhere(params.Conditions, &args); cond != "" {
		sql += " WHERE " + cond
	}
	if len(params.OrderBys) < 1 {
		params.OrderBys = append(params.OrderBys, model.OrderBy{Field: model.DBLinkGenericLinkID})
	}
	sql += " ORDER BY " + SetOrderBys(params.OrderBys, &args)
	if params.Pagination == nil {
		params.Pagination = &model.Pagination{Page: 1, Limit: model.ComicLinkPaginationDef}
	}
	if lmof := SetPagination(*params.Pagination, &args); lmof != "" {
		sql += lmof
	}
	if err := db.QueryAll(ctx, &result, sql, args...); err != nil {
		return nil, err
	}
	return result, nil
}

func (db Database) CountComicLink(ctx context.Context, conds any) (int, error) {
	return db.GenericCount(ctx, model.DBComicLink, conds)
}

func comicLinkSetError(err error) error {
	var errDatabase model.DatabaseError
	if errors.As(err, &errDatabase) {
		if errDatabase.Code == CodeErrForeign {
			switch errDatabase.Name {
			case NameErrComicLinkFKey0:
				return model.GenericError("comic does not exist")
			case NameErrComicLinkFKey1:
				return model.GenericError("link does not exist")
			}
		}
		if errDatabase.Code == CodeErrExists && errDatabase.Name == NameErrComicLinkPKey {
			return model.GenericError("same link id already exists")
		}
	}
	return err
}

const (
	NameErrComicChapterFKey = "comic_chapter_comic_id_fkey"
	NameErrComicChapterKey  = "comic_chapter_comic_id_chapter_version_key"
)

func (db Database) AddComicChapter(ctx context.Context, data model.AddComicChapter, v *model.ComicChapter) error {
	var comicID any
	switch {
	case data.ComicID != nil:
		comicID = data.ComicID
	case data.ComicCode != nil:
		comicID = model.DBComicCodeToID(*data.ComicCode)
	}
	cols, vals, args := SetInsert(map[string]any{
		model.DBComicGenericComicID:    comicID,
		model.DBComicChapterChapter:    data.Chapter,
		model.DBComicChapterVersion:    data.Version,
		model.DBComicChapterReleasedAt: data.ReleasedAt,
	})
	sql := "INSERT INTO " + model.DBComicChapter + " (" + cols + ") VALUES (" + vals + ")"
	if v != nil {
		sql += " RETURNING *"
		sql = "WITH data AS (" + sql + ")"
		sql += " SELECT w." + model.DBGenericID
		sql += ", w." + model.DBGenericCreatedAt + ", w." + model.DBGenericUpdatedAt
		sql += ", w." + model.DBComicGenericComicID + ", w." + model.DBComicChapterChapter
		sql += ", w." + model.DBComicChapterVersion + ", w." + model.DBComicChapterReleasedAt
		sql += ", l." + model.DBComicCode + " AS comic_code"
		sql += " FROM data w JOIN " + model.DBComic + " l"
		sql += " ON w." + model.DBComicGenericComicID + " = l." + model.DBGenericID
		if err := db.QueryOne(ctx, v, sql, args...); err != nil {
			return comicChapterSetError(err)
		}
	} else {
		if err := db.Exec(ctx, sql, args...); err != nil {
			return comicChapterSetError(err)
		}
	}
	return nil
}

func (db Database) GetComicChapter(ctx context.Context, conds any) (*model.ComicChapter, error) {
	var result model.ComicChapter
	args := []any{}
	cond := SetWhere(conds, &args)
	sql := "SELECT * FROM (SELECT w." + model.DBGenericID
	sql += ", w." + model.DBGenericCreatedAt + ", w." + model.DBGenericUpdatedAt
	sql += ", w." + model.DBComicGenericComicID + ", w." + model.DBComicChapterChapter
	sql += ", w." + model.DBComicChapterVersion + ", w." + model.DBComicChapterReleasedAt
	sql += ", l." + model.DBComicCode + " AS comic_code"
	sql += " FROM " + model.DBComicChapter + " w JOIN " + model.DBComic + " l"
	sql += " ON w." + model.DBComicGenericComicID + " = l." + model.DBGenericID
	sql += ")"
	sql += " WHERE " + cond
	if err := db.QueryOne(ctx, &result, sql, args...); err != nil {
		return nil, err
	}
	return &result, nil
}

func (db Database) UpdateComicChapter(ctx context.Context, data model.SetComicChapter, conds any, v *model.ComicChapter) error {
	data0 := map[string]any{}
	switch {
	case data.ComicID != nil:
		data0[model.DBComicGenericComicID] = data.ComicID
	case data.ComicCode != nil:
		data0[model.DBComicGenericComicID] = model.DBComicCodeToID(*data.ComicCode)
	}
	if data.Chapter != nil {
		data0[model.DBComicChapterChapter] = data.Chapter
	}
	if data.Version != nil {
		data0[model.DBComicChapterVersion] = data.Version
	}
	if data.ReleasedAt != nil {
		data0[model.DBComicChapterReleasedAt] = data.ReleasedAt
	}
	for _, null := range data.SetNull {
		data0[null] = nil
	}
	data0[model.DBGenericUpdatedAt] = time.Now().UTC()
	sets, args := SetUpdate(data0)
	cond := SetWhere([]any{conds, SetUpdateWhere(data0)}, &args)
	sql := "UPDATE " + model.DBComicChapter + " SET " + sets + " WHERE " + cond
	if v != nil {
		sql += " RETURNING *"
		sql = "WITH data AS (" + sql + ")"
		sql += " SELECT w." + model.DBGenericID
		sql += ", w." + model.DBGenericCreatedAt + ", w." + model.DBGenericUpdatedAt
		sql += ", w." + model.DBComicGenericComicID + ", w." + model.DBComicChapterChapter
		sql += ", w." + model.DBComicChapterVersion + ", w." + model.DBComicChapterReleasedAt
		sql += ", l." + model.DBComicCode + " AS comic_code"
		sql += " FROM data w JOIN " + model.DBComic + " l"
		sql += " ON w." + model.DBComicGenericComicID + " = l." + model.DBGenericID
		if err := db.QueryOne(ctx, v, sql, args...); err != nil {
			return comicChapterSetError(err)
		}
	} else {
		if err := db.Exec(ctx, sql, args...); err != nil {
			return comicChapterSetError(err)
		}
	}
	return nil
}

func (db Database) DeleteComicChapter(ctx context.Context, conds any, v *model.ComicChapter) error {
	args := []any{}
	cond := SetWhere(conds, &args)
	sql := "DELETE FROM " + model.DBComicChapter + " WHERE " + cond
	if v != nil {
		sql += " RETURNING *"
		sql = "WITH data AS (" + sql + ")"
		sql += " SELECT w." + model.DBGenericID
		sql += ", w." + model.DBGenericCreatedAt + ", w." + model.DBGenericUpdatedAt
		sql += ", w." + model.DBComicGenericComicID + ", w." + model.DBComicChapterChapter
		sql += ", w." + model.DBComicChapterVersion + ", w." + model.DBComicChapterReleasedAt
		sql += ", l." + model.DBComicCode + " AS comic_code"
		sql += " FROM data w JOIN " + model.DBComic + " l"
		sql += " ON w." + model.DBComicGenericComicID + " = l." + model.DBGenericID
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

func (db Database) ListComicChapter(ctx context.Context, params model.ListParams) ([]*model.ComicChapter, error) {
	result := []*model.ComicChapter{}
	args := []any{}
	sql := "SELECT * FROM (SELECT w." + model.DBGenericID
	sql += ", w." + model.DBGenericCreatedAt + ", w." + model.DBGenericUpdatedAt
	sql += ", w." + model.DBComicGenericComicID + ", w." + model.DBComicChapterChapter
	sql += ", w." + model.DBComicChapterVersion + ", w." + model.DBComicChapterReleasedAt
	sql += ", l." + model.DBComicCode + " AS comic_code"
	sql += " FROM " + model.DBComicChapter + " w JOIN " + model.DBComic + " l"
	sql += " ON w." + model.DBComicGenericComicID + " = l." + model.DBGenericID
	sql += ")"
	if cond := SetWhere(params.Conditions, &args); cond != "" {
		sql += " WHERE " + cond
	}
	if len(params.OrderBys) < 1 {
		params.OrderBys = append(params.OrderBys, model.OrderBy{Field: model.DBComicChapterReleasedAt})
	}
	sql += " ORDER BY " + SetOrderBys(params.OrderBys, &args)
	if params.Pagination == nil {
		params.Pagination = &model.Pagination{Page: 1, Limit: model.ComicChapterPaginationDef}
	}
	if lmof := SetPagination(*params.Pagination, &args); lmof != "" {
		sql += lmof
	}
	if err := db.QueryAll(ctx, &result, sql, args...); err != nil {
		return nil, err
	}
	return result, nil
}

func (db Database) CountComicChapter(ctx context.Context, conds any) (int, error) {
	return db.GenericCount(ctx, model.DBComicChapter, conds)
}

func (db Database) ExistsComicChapter(ctx context.Context, conds any) (bool, error) {
	return db.GenericExists(ctx, model.DBComicChapter, conds)
}

func comicChapterSetError(err error) error {
	var errDatabase model.DatabaseError
	if errors.As(err, &errDatabase) {
		if errDatabase.Code == CodeErrForeign && errDatabase.Name == NameErrComicChapterFKey {
			return model.GenericError("comic does not exist")
		}
		if errDatabase.Code == CodeErrExists && errDatabase.Name == NameErrComicChapterKey {
			return model.GenericError("same comic id + chapter + version already exists")
		}
	}
	return err
}

const (
	NameErrComicChapterLinkPKey  = "comic_chapter_link_pkey"
	NameErrComicChapterLinkFKey0 = "comic_chapter_link_chapter_id_fkey"
	NameErrComicChapterLinkFKey1 = "comic_chapter_link_link_id_fkey"
)

func (db Database) AddComicChapterLink(ctx context.Context, data model.AddComicChapterLink, v *model.ComicChapterLink) error {
	var chapterID any
	switch {
	case data.ChapterID != nil:
		chapterID = data.ChapterID
	case data.ChapterSID != nil:
		chapterID = model.DBComicChapterSIDToID(*data.ChapterSID)
	}
	var linkID any
	switch {
	case data.LinkID != nil:
		linkID = data.LinkID
	case data.LinkSID != nil:
		linkID = model.DBLinkSIDToID(*data.LinkSID)
	}
	cols, vals, args := SetInsert(map[string]any{
		model.DBComicChapterGenericChapterID: chapterID,
		model.DBLinkGenericLinkID:            linkID,
	})
	sql := "INSERT INTO " + model.DBComicChapterLink + " (" + cols + ") VALUES (" + vals + ")"
	if v != nil {
		sql += " RETURNING *"
		sql = "WITH data AS (" + sql + ")"
		sql += " SELECT a." + model.DBGenericCreatedAt + ", a." + model.DBGenericUpdatedAt
		sql += ", a." + model.DBComicChapterGenericChapterID + ", a." + model.DBLinkGenericLinkID
		sql += ", b." + model.DBLinkRelativeURL + " AS link_relative_url"
		sql += ", c." + model.DBWebsiteDomain + " AS link_website_domain"
		sql += " FROM data a JOIN " + model.DBLink + " b"
		sql += " ON a." + model.DBLinkGenericLinkID + " = b." + model.DBGenericID
		sql += " JOIN " + model.DBWebsite + " c"
		sql += " ON b." + model.DBWebsiteGenericWebsiteID + " = c." + model.DBGenericID
		if err := db.QueryOne(ctx, v, sql, args...); err != nil {
			return comicChapterLinkSetError(err)
		}
	} else {
		if err := db.Exec(ctx, sql, args...); err != nil {
			return comicChapterLinkSetError(err)
		}
	}
	return nil
}

func (db Database) GetComicChapterLink(ctx context.Context, conds any) (*model.ComicChapterLink, error) {
	var result model.ComicChapterLink
	args := []any{}
	cond := SetWhere(conds, &args)
	sql := "SELECT * FROM (SELECT a." + model.DBGenericCreatedAt + ", a." + model.DBGenericUpdatedAt
	sql += ", a." + model.DBComicChapterGenericChapterID + ", a." + model.DBLinkGenericLinkID
	sql += ", b." + model.DBLinkRelativeURL + " AS link_relative_url"
	sql += ", c." + model.DBWebsiteDomain + " AS link_website_domain"
	sql += " FROM " + model.DBComicChapterLink + " a JOIN " + model.DBLink + " b"
	sql += " ON a." + model.DBLinkGenericLinkID + " = b." + model.DBGenericID
	sql += " JOIN " + model.DBWebsite + " c"
	sql += " ON b." + model.DBWebsiteGenericWebsiteID + " = c." + model.DBGenericID
	sql += ")"
	sql += " WHERE " + cond
	if err := db.QueryOne(ctx, &result, sql, args...); err != nil {
		return nil, err
	}
	return &result, nil
}

func (db Database) UpdateComicChapterLink(ctx context.Context, data model.SetComicChapterLink, conds any, v *model.ComicChapterLink) error {
	data0 := map[string]any{}
	switch {
	case data.ChapterID != nil:
		data0[model.DBComicChapterGenericChapterID] = data.ChapterID
	case data.ChapterSID != nil:
		data0[model.DBComicChapterGenericChapterID] = model.DBComicChapterSIDToID(*data.ChapterSID)
	}
	switch {
	case data.LinkID != nil:
		data0[model.DBLinkGenericLinkID] = data.LinkID
	case data.LinkSID != nil:
		data0[model.DBLinkGenericLinkID] = model.DBLinkSIDToID(*data.LinkSID)
	}
	data0[model.DBGenericUpdatedAt] = time.Now().UTC()
	sets, args := SetUpdate(data0)
	cond := SetWhere([]any{conds, SetUpdateWhere(data0)}, &args)
	sql := "UPDATE " + model.DBComicChapterLink + " SET " + sets + " WHERE " + cond
	if v != nil {
		sql += " RETURNING *"
		sql = "WITH data AS (" + sql + ")"
		sql += " SELECT a." + model.DBGenericCreatedAt + ", a." + model.DBGenericUpdatedAt
		sql += ", a." + model.DBComicChapterGenericChapterID + ", a." + model.DBLinkGenericLinkID
		sql += ", b." + model.DBLinkRelativeURL + " AS link_relative_url"
		sql += ", c." + model.DBWebsiteDomain + " AS link_website_domain"
		sql += " FROM data a JOIN " + model.DBLink + " b"
		sql += " ON a." + model.DBLinkGenericLinkID + " = b." + model.DBGenericID
		sql += " JOIN " + model.DBWebsite + " c"
		sql += " ON b." + model.DBWebsiteGenericWebsiteID + " = c." + model.DBGenericID
		if err := db.QueryOne(ctx, v, sql, args...); err != nil {
			return comicChapterLinkSetError(err)
		}
	} else {
		if err := db.Exec(ctx, sql, args...); err != nil {
			return comicChapterLinkSetError(err)
		}
	}
	return nil
}

func (db Database) DeleteComicChapterLink(ctx context.Context, conds any, v *model.ComicChapterLink) error {
	args := []any{}
	cond := SetWhere(conds, &args)
	sql := "DELETE FROM " + model.DBComicChapterLink + " WHERE " + cond
	if v != nil {
		sql += " RETURNING *"
		sql = "WITH data AS (" + sql + ")"
		sql += " SELECT a." + model.DBGenericCreatedAt + ", a." + model.DBGenericUpdatedAt
		sql += ", a." + model.DBComicChapterGenericChapterID + ", a." + model.DBLinkGenericLinkID
		sql += ", b." + model.DBLinkRelativeURL + " AS link_relative_url"
		sql += ", c." + model.DBWebsiteDomain + " AS link_website_domain"
		sql += " FROM data a JOIN " + model.DBLink + " b"
		sql += " ON a." + model.DBLinkGenericLinkID + " = b." + model.DBGenericID
		sql += " JOIN " + model.DBWebsite + " c"
		sql += " ON b." + model.DBWebsiteGenericWebsiteID + " = c." + model.DBGenericID
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

func (db Database) ListComicChapterLink(ctx context.Context, params model.ListParams) ([]*model.ComicChapterLink, error) {
	result := []*model.ComicChapterLink{}
	args := []any{}
	sql := "SELECT * FROM (SELECT a." + model.DBGenericCreatedAt + ", a." + model.DBGenericUpdatedAt
	sql += ", a." + model.DBComicChapterGenericChapterID + ", a." + model.DBLinkGenericLinkID
	sql += ", b." + model.DBLinkRelativeURL + " AS link_relative_url"
	sql += ", c." + model.DBWebsiteDomain + " AS link_website_domain"
	sql += " FROM " + model.DBComicChapterLink + " a JOIN " + model.DBLink + " b"
	sql += " ON a." + model.DBLinkGenericLinkID + " = b." + model.DBGenericID
	sql += " JOIN " + model.DBWebsite + " c"
	sql += " ON b." + model.DBWebsiteGenericWebsiteID + " = c." + model.DBGenericID
	sql += ")"
	if cond := SetWhere(params.Conditions, &args); cond != "" {
		sql += " WHERE " + cond
	}
	if len(params.OrderBys) < 1 {
		params.OrderBys = append(params.OrderBys, model.OrderBy{Field: model.DBLinkGenericLinkID})
	}
	sql += " ORDER BY " + SetOrderBys(params.OrderBys, &args)
	if params.Pagination == nil {
		params.Pagination = &model.Pagination{Page: 1, Limit: model.ComicChapterLinkPaginationDef}
	}
	if lmof := SetPagination(*params.Pagination, &args); lmof != "" {
		sql += lmof
	}
	if err := db.QueryAll(ctx, &result, sql, args...); err != nil {
		return nil, err
	}
	return result, nil
}

func (db Database) CountComicChapterLink(ctx context.Context, conds any) (int, error) {
	return db.GenericCount(ctx, model.DBComicChapterLink, conds)
}

func comicChapterLinkSetError(err error) error {
	var errDatabase model.DatabaseError
	if errors.As(err, &errDatabase) {
		if errDatabase.Code == CodeErrForeign {
			switch errDatabase.Name {
			case NameErrComicChapterLinkFKey0:
				return model.GenericError("chapter does not exist")
			case NameErrComicChapterLinkFKey1:
				return model.GenericError("website does not exist")
			}
		}
		if errDatabase.Code == CodeErrExists && errDatabase.Name == NameErrComicChapterLinkPKey {
			return model.GenericError("same link id already exists")
		}
	}
	return err
}
