package service

import (
	"context"
	"slices"

	"github.com/mahmudindes/orenocomic-bagicore/internal/model"
	"golang.org/x/sync/errgroup"
)

//
// Comic
//

func (svc Service) AddComic(ctx context.Context, data model.AddComic, v *model.Comic) error {
	if !svc.oauth.HasPermissionContext(ctx, svc.oauth.TokenPermissionKey("write")) {
		return model.GenericError("missing admin permission to add comic")
	}

	if err := data.Validate(); err != nil {
		return err
	}

	if v != nil {
		v.Links = []*model.Link{}
		v.Chapters = []*model.ComicChapter{}
	}

	return svc.database.AddComic(ctx, data, v)
}

func (svc Service) GetComicByCode(ctx context.Context, code string) (*model.Comic, error) {
	result, err := svc.database.GetComic(ctx, model.DBConditionalKV{
		Key:   model.DBComicCode,
		Value: code,
	})
	if err != nil {
		return nil, err
	}

	g, gctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		links0, err := svc.listComicLink(ctx, model.ListParams{
			Conditions: model.DBConditionalKV{Key: model.DBComicGenericComicID, Value: result.ID},
			Pagination: &model.Pagination{},
		})
		if err != nil {
			return err
		}

		if len(links0) == 0 {
			result.Links = []*model.Link{}
			return nil
		}

		conditions := make([]any, len(links0)+2)
		conditions = append(conditions, model.DBLogicalOR{})
		for _, link := range links0 {
			conditions = append(conditions, model.DBConditionalKV{
				Key:   model.DBGenericID,
				Value: link.LinkID,
			})
		}

		links1, err := svc.database.ListLink(ctx, model.ListParams{
			Conditions: conditions,
			Pagination: &model.Pagination{},
		})
		if err != nil {
			return err
		}

		result.Links = links1
		return nil
	})
	g.Go(func() error {
		chapters, err := svc.listComicChapter(gctx, model.ListParams{
			Conditions: model.DBConditionalKV{Key: model.DBComicGenericComicID, Value: result.ID},
			Pagination: &model.Pagination{},
		})
		if err != nil {
			return err
		}

		result.Chapters = chapters
		return nil
	})
	if err := g.Wait(); err != nil {
		return nil, err
	}

	return result, nil
}

func (svc Service) UpdateComicByCode(ctx context.Context, code string, data model.SetComic, v *model.Comic) error {
	if !svc.oauth.HasPermissionContext(ctx, svc.oauth.TokenPermissionKey("write")) {
		return model.GenericError("missing admin permission to update comic")
	}

	if err := data.Validate(); err != nil {
		return err
	}

	if err := svc.database.UpdateComic(ctx, data, model.DBConditionalKV{
		Key:   model.DBComicCode,
		Value: code,
	}, v); err != nil {
		return err
	}

	if v != nil {
		g, gctx := errgroup.WithContext(ctx)
		g.Go(func() error {
			links0, err := svc.listComicLink(ctx, model.ListParams{
				Conditions: model.DBConditionalKV{Key: model.DBComicGenericComicID, Value: v.ID},
				Pagination: &model.Pagination{},
			})
			if err != nil {
				return err
			}

			if len(links0) == 0 {
				v.Links = []*model.Link{}
				return nil
			}

			conditions := make([]any, len(links0)+2)
			conditions = append(conditions, model.DBLogicalOR{})
			for _, link := range links0 {
				conditions = append(conditions, model.DBConditionalKV{
					Key:   model.DBGenericID,
					Value: link.LinkID,
				})
			}

			links1, err := svc.database.ListLink(ctx, model.ListParams{
				Conditions: conditions,
				Pagination: &model.Pagination{},
			})
			if err != nil {
				return err
			}

			v.Links = links1
			return nil
		})
		g.Go(func() error {
			chapters, err := svc.listComicChapter(gctx, model.ListParams{
				Conditions: model.DBConditionalKV{Key: model.DBComicGenericComicID, Value: v.ID},
				Pagination: &model.Pagination{},
			})
			if err != nil {
				return err
			}

			v.Chapters = chapters
			return nil
		})
		if err := g.Wait(); err != nil {
			return err
		}
	}

	return nil
}

func (svc Service) DeleteComicByCode(ctx context.Context, code string) error {
	if !svc.oauth.HasPermissionContext(ctx, svc.oauth.TokenPermissionKey("write")) {
		return model.GenericError("missing admin permission to delete comic")
	}

	return svc.database.DeleteComic(ctx, model.DBConditionalKV{
		Key:   model.DBComicCode,
		Value: code,
	}, nil)
}

func (svc Service) ListComic(ctx context.Context, params model.ListParams) ([]*model.Comic, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	params.OrderBys = slices.DeleteFunc(params.OrderBys, func(ob model.OrderBy) bool {
		switch field := ob.Field.(type) {
		case string:
			return !slices.Contains(model.ComicOrderByAllow, field)
		}
		return true
	})
	if len(params.OrderBys) > model.ComicOrderBysMax {
		params.OrderBys = params.OrderBys[:model.ComicOrderBysMax]
	}
	if pagination := params.Pagination; pagination != nil {
		if pagination.Limit > model.ComicPaginationMax {
			pagination.Limit = model.ComicPaginationMax
		}
	}

	result, err := svc.database.ListComic(ctx, params)
	if err != nil {
		return nil, err
	}

	if len(result) > 0 {
		conds := make([]any, len(result)+1)
		conds = append(conds, model.DBLogicalOR{})
		for _, r := range result {
			conds = append(conds, model.DBConditionalKV{
				Key:   model.DBComicGenericComicID,
				Value: r.ID,
			})
		}
		g, gctx := errgroup.WithContext(ctx)
		g.Go(func() error {
			links0, err := svc.listComicLink(ctx, model.ListParams{
				Conditions: conds,
				Pagination: &model.Pagination{},
			})
			if err != nil {
				return err
			}
			links := map[uint]*model.Link{}
			for _, link := range links0 {
				links[link.LinkID] = nil
			}
			conditions := make([]any, len(links0)+2)
			conditions = append(conditions, model.DBLogicalOR{})
			for id := range links {
				conditions = append(conditions, model.DBConditionalKV{
					Key:   model.DBGenericID,
					Value: id,
				})
			}
			links1, err := svc.database.ListLink(ctx, model.ListParams{
				Conditions: conditions,
				Pagination: &model.Pagination{},
			})
			if err != nil {
				return err
			}
			for _, link := range links1 {
				links[link.ID] = link
			}
			for _, r := range result {
				r.Links = make([]*model.Link, 0)
			}
			for _, link := range links0 {
				for _, r := range result {
					if r.ID == link.LinkID {
						r.Links = append(r.Links, links[link.LinkID])
					}
				}
			}
			return nil
		})
		g.Go(func() error {
			chapters, err := svc.listComicChapter(gctx, model.ListParams{
				Conditions: conds,
				Pagination: &model.Pagination{},
			})
			if err != nil {
				return err
			}
			for _, r := range result {
				r.Chapters = []*model.ComicChapter{}
			}
			for _, chapter := range chapters {
				for _, r := range result {
					if r.ID == chapter.ComicID {
						r.Chapters = append(r.Chapters, chapter)
					}
				}
			}
			return nil
		})
		if err := g.Wait(); err != nil {
			return nil, err
		}
	}

	return result, nil
}

func (svc Service) CountComic(ctx context.Context, conds any) (int, error) {
	return svc.database.CountComic(ctx, conds)
}

func (svc Service) ExistsComicByCode(ctx context.Context, code string) (bool, error) {
	return svc.database.ExistsComic(ctx, model.DBConditionalKV{
		Key:   model.DBComicCode,
		Value: code,
	})
}

// Comic Link

func (svc Service) AddComicLink(ctx context.Context, data model.AddComicLink, v *model.ComicLink) error {
	if !svc.oauth.HasPermissionContext(ctx, svc.oauth.TokenPermissionKey("write")) {
		return model.GenericError("missing admin permission to add comic link")
	}

	if err := data.Validate(); err != nil {
		return err
	}

	return svc.database.AddComicLink(ctx, data, v)
}

func (svc Service) GetComicLinkBySID(ctx context.Context, sid model.ComicLinkSID) (*model.ComicLink, error) {
	var comicID any
	switch {
	case sid.ComicID != nil:
		comicID = sid.ComicID
	case sid.ComicCode != nil:
		comicID = model.DBComicCodeToID(*sid.ComicCode)
	}
	var linkID any
	switch {
	case sid.LinkID != nil:
		linkID = sid.LinkID
	case sid.LinkSID != nil:
		linkID = model.DBLinkSIDToID(*sid.LinkSID)
	}
	result, err := svc.database.GetComicLink(ctx, map[string]any{
		model.DBComicGenericComicID: comicID,
		model.DBLinkGenericLinkID:   linkID,
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (svc Service) UpdateComicLinkBySID(ctx context.Context, sid model.ComicLinkSID, data model.SetComicLink, v *model.ComicLink) error {
	if !svc.oauth.HasPermissionContext(ctx, svc.oauth.TokenPermissionKey("write")) {
		return model.GenericError("missing admin permission to update comic link")
	}

	if err := data.Validate(); err != nil {
		return err
	}

	var comicID any
	switch {
	case sid.ComicID != nil:
		comicID = sid.ComicID
	case sid.ComicCode != nil:
		comicID = model.DBComicCodeToID(*sid.ComicCode)
	}
	var linkID any
	switch {
	case sid.LinkID != nil:
		linkID = sid.LinkID
	case sid.LinkSID != nil:
		linkID = model.DBLinkSIDToID(*sid.LinkSID)
	}
	if err := svc.database.UpdateComicLink(ctx, data, map[string]any{
		model.DBComicGenericComicID: comicID,
		model.DBLinkGenericLinkID:   linkID,
	}, v); err != nil {
		return err
	}

	return nil
}

func (svc Service) DeleteComicLinkBySID(ctx context.Context, sid model.ComicLinkSID) error {
	if !svc.oauth.HasPermissionContext(ctx, svc.oauth.TokenPermissionKey("write")) {
		return model.GenericError("missing admin permission to delete comic link")
	}

	var comicID any
	switch {
	case sid.ComicID != nil:
		comicID = sid.ComicID
	case sid.ComicCode != nil:
		comicID = model.DBComicCodeToID(*sid.ComicCode)
	}
	var linkID any
	switch {
	case sid.LinkID != nil:
		linkID = sid.LinkID
	case sid.LinkSID != nil:
		linkID = model.DBLinkSIDToID(*sid.LinkSID)
	}
	return svc.database.DeleteComicLink(ctx, map[string]any{
		model.DBComicGenericComicID: comicID,
		model.DBLinkGenericLinkID:   linkID,
	}, nil)
}

func (svc Service) listComicLink(ctx context.Context, params model.ListParams) ([]*model.ComicLink, error) {
	result, err := svc.database.ListComicLink(ctx, params)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (svc Service) ListComicLink(ctx context.Context, params model.ListParams) ([]*model.ComicLink, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	params.OrderBys = slices.DeleteFunc(params.OrderBys, func(ob model.OrderBy) bool {
		switch field := ob.Field.(type) {
		case string:
			return !slices.Contains(model.ComicLinkOrderByAllow, field)
		}
		return true
	})
	if len(params.OrderBys) > model.ComicLinkOrderBysMax {
		params.OrderBys = params.OrderBys[:model.ComicLinkOrderBysMax]
	}
	if pagination := params.Pagination; pagination != nil {
		if pagination.Limit > model.ComicLinkPaginationMax {
			pagination.Limit = model.ComicLinkPaginationMax
		}
	}

	return svc.database.ListComicLink(ctx, params)
}

func (svc Service) CountComicLink(ctx context.Context, conds any) (int, error) {
	return svc.database.CountComicLink(ctx, conds)
}

//
// Comic Chapter
//

func (svc Service) AddComicChapter(ctx context.Context, data model.AddComicChapter, v *model.ComicChapter) error {
	if !svc.oauth.HasPermissionContext(ctx, svc.oauth.TokenPermissionKey("write")) {
		return model.GenericError("missing admin permission to add comic chapter")
	}

	if err := data.Validate(); err != nil {
		return err
	}

	if v != nil {
		v.Links = []*model.Link{}
	}

	return svc.database.AddComicChapter(ctx, data, v)
}

func (svc Service) getComicChapter(ctx context.Context, conds any) (*model.ComicChapter, error) {
	result, err := svc.database.GetComicChapter(ctx, conds)
	if err != nil {
		return nil, err
	}

	links0, err := svc.listComicChapterLink(ctx, model.ListParams{
		Conditions: model.DBConditionalKV{Key: model.DBComicChapterGenericChapterID, Value: result.ID},
		Pagination: &model.Pagination{},
	})
	if err != nil {
		return nil, err
	}
	result.Links = []*model.Link{}
	if len(links0) > 0 {
		conditions := make([]any, len(links0)+2)
		conditions = append(conditions, model.DBLogicalOR{})
		for _, link := range links0 {
			conditions = append(conditions, model.DBConditionalKV{
				Key:   model.DBGenericID,
				Value: link.LinkID,
			})
		}
		links1, err := svc.database.ListLink(ctx, model.ListParams{
			Conditions: conditions,
			Pagination: &model.Pagination{},
		})
		if err != nil {
			return nil, err
		}
		result.Links = links1
	}

	return result, nil
}

func (svc Service) GetComicChapterBySID(ctx context.Context, sid model.ComicChapterSID) (*model.ComicChapter, error) {
	var comicID any
	switch {
	case sid.ComicID != nil:
		comicID = sid.ComicID
	case sid.ComicCode != nil:
		comicID = model.DBComicCodeToID(*sid.ComicCode)
	}
	var version any
	switch {
	case sid.Version != nil:
		version = sid.Version
	default:
		version = model.DBIsNull{}
	}
	return svc.getComicChapter(ctx, map[string]any{
		model.DBComicGenericComicID: comicID,
		model.DBComicChapterChapter: sid.Chapter,
		model.DBComicChapterVersion: version,
	})
}

func (svc Service) updateComicChapter(ctx context.Context, data model.SetComicChapter, conds any, v *model.ComicChapter) error {
	if !svc.oauth.HasPermissionContext(ctx, svc.oauth.TokenPermissionKey("write")) {
		return model.GenericError("missing admin permission to update comic chapter")
	}

	if err := data.Validate(); err != nil {
		return err
	}

	if err := svc.database.UpdateComicChapter(ctx, data, conds, v); err != nil {
		return err
	}

	if v != nil {
		links0, err := svc.listComicChapterLink(ctx, model.ListParams{
			Conditions: model.DBConditionalKV{Key: model.DBComicChapterGenericChapterID, Value: v.ID},
			Pagination: &model.Pagination{},
		})
		if err != nil {
			return err
		}
		v.Links = []*model.Link{}
		if len(links0) > 0 {
			conditions := make([]any, len(links0)+2)
			conditions = append(conditions, model.DBLogicalOR{})
			for _, link := range links0 {
				conditions = append(conditions, model.DBConditionalKV{
					Key:   model.DBGenericID,
					Value: link.LinkID,
				})
			}
			links1, err := svc.database.ListLink(ctx, model.ListParams{
				Conditions: conditions,
				Pagination: &model.Pagination{},
			})
			if err != nil {
				return err
			}
			v.Links = links1
		}
	}

	return nil
}

func (svc Service) UpdateComicChapterBySID(ctx context.Context, sid model.ComicChapterSID, data model.SetComicChapter, v *model.ComicChapter) error {
	var comicID any
	switch {
	case sid.ComicID != nil:
		comicID = sid.ComicID
	case sid.ComicCode != nil:
		comicID = model.DBComicCodeToID(*sid.ComicCode)
	}
	var version any
	switch {
	case sid.Version != nil:
		version = sid.Version
	default:
		version = model.DBIsNull{}
	}
	return svc.updateComicChapter(ctx, data, map[string]any{
		model.DBComicGenericComicID: comicID,
		model.DBComicChapterChapter: sid.Chapter,
		model.DBComicChapterVersion: version,
	}, v)
}

func (svc Service) deleteComicChapter(ctx context.Context, conds any) error {
	if !svc.oauth.HasPermissionContext(ctx, svc.oauth.TokenPermissionKey("write")) {
		return model.GenericError("missing admin permission to delete comic chapter")
	}

	return svc.database.DeleteComicChapter(ctx, conds, nil)
}

func (svc Service) DeleteComicChapterBySID(ctx context.Context, sid model.ComicChapterSID) error {
	var comicID any
	switch {
	case sid.ComicID != nil:
		comicID = sid.ComicID
	case sid.ComicCode != nil:
		comicID = model.DBComicCodeToID(*sid.ComicCode)
	}
	var version any
	switch {
	case sid.Version != nil:
		version = sid.Version
	default:
		version = model.DBIsNull{}
	}
	return svc.deleteComicChapter(ctx, map[string]any{
		model.DBComicGenericComicID: comicID,
		model.DBComicChapterChapter: sid.Chapter,
		model.DBComicChapterVersion: version,
	})
}

func (svc Service) listComicChapter(ctx context.Context, params model.ListParams) ([]*model.ComicChapter, error) {
	result, err := svc.database.ListComicChapter(ctx, params)
	if err != nil {
		return nil, err
	}

	if len(result) > 0 {
		conds := make([]any, len(result)+1)
		conds = append(conds, model.DBLogicalOR{})
		for _, r := range result {
			conds = append(conds, model.DBConditionalKV{
				Key:   model.DBComicChapterGenericChapterID,
				Value: r.ID,
			})
		}
		links0, err := svc.listComicChapterLink(ctx, model.ListParams{
			Conditions: conds,
			Pagination: &model.Pagination{},
		})
		if err != nil {
			return nil, err
		}
		links := map[uint]*model.Link{}
		for _, link := range links0 {
			links[link.LinkID] = nil
		}
		conditions := make([]any, len(links0)+2)
		conditions = append(conditions, model.DBLogicalOR{})
		for id := range links {
			conditions = append(conditions, model.DBConditionalKV{
				Key:   model.DBGenericID,
				Value: id,
			})
		}
		links1, err := svc.database.ListLink(ctx, model.ListParams{
			Conditions: conditions,
			Pagination: &model.Pagination{},
		})
		if err != nil {
			return nil, err
		}
		for _, link := range links1 {
			links[link.ID] = link
		}
		for _, r := range result {
			r.Links = make([]*model.Link, 0)
		}
		for _, link := range links0 {
			for _, r := range result {
				if r.ID == link.ChapterID {
					r.Links = append(r.Links, links[link.LinkID])
				}
			}
		}
	}

	return result, nil
}

func (svc Service) ListComicChapter(ctx context.Context, params model.ListParams) ([]*model.ComicChapter, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	params.OrderBys = slices.DeleteFunc(params.OrderBys, func(ob model.OrderBy) bool {
		switch field := ob.Field.(type) {
		case string:
			return !slices.Contains(model.ComicChapterOrderByAllow, field)
		}
		return true
	})
	if len(params.OrderBys) > model.ComicChapterOrderBysMax {
		params.OrderBys = params.OrderBys[:model.ComicChapterOrderBysMax]
	}
	if pagination := params.Pagination; pagination != nil {
		if pagination.Limit > model.ComicChapterPaginationMax {
			pagination.Limit = model.ComicChapterPaginationMax
		}
	}

	return svc.listComicChapter(ctx, params)
}

func (svc Service) CountComicChapter(ctx context.Context, conds any) (int, error) {
	return svc.database.CountComicChapter(ctx, conds)
}

func (svc Service) ExistsComicChapterBySID(ctx context.Context, sid model.ComicChapterSID) (bool, error) {
	var comicID any
	switch {
	case sid.ComicID != nil:
		comicID = sid.ComicID
	case sid.ComicCode != nil:
		comicID = model.DBComicCodeToID(*sid.ComicCode)
	}
	var version any
	switch {
	case sid.Version != nil:
		version = sid.Version
	default:
		version = model.DBIsNull{}
	}
	return svc.database.ExistsComicChapter(ctx, map[string]any{
		model.DBComicGenericComicID: comicID,
		model.DBComicChapterChapter: sid.Chapter,
		model.DBComicChapterVersion: version,
	})
}

// Comic Chapter Link

func (svc Service) AddComicChapterLink(ctx context.Context, data model.AddComicChapterLink, v *model.ComicChapterLink) error {
	if !svc.oauth.HasPermissionContext(ctx, svc.oauth.TokenPermissionKey("write")) {
		return model.GenericError("missing admin permission to add comic chapter link")
	}

	if err := data.Validate(); err != nil {
		return err
	}

	return svc.database.AddComicChapterLink(ctx, data, v)
}

func (svc Service) GetComicChapterLinkBySID(ctx context.Context, sid model.ComicChapterLinkSID) (*model.ComicChapterLink, error) {
	var chapterID any
	switch {
	case sid.ChapterID != nil:
		chapterID = sid.ChapterID
	case sid.ChapterSID != nil:
		chapterID = model.DBComicChapterSIDToID(*sid.ChapterSID)
	}
	var linkID any
	switch {
	case sid.LinkID != nil:
		linkID = sid.LinkID
	case sid.LinkSID != nil:
		linkID = model.DBLinkSIDToID(*sid.LinkSID)
	}
	result, err := svc.database.GetComicChapterLink(ctx, map[string]any{
		model.DBComicChapterGenericChapterID: chapterID,
		model.DBLinkGenericLinkID:            linkID,
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (svc Service) UpdateComicChapterLinkBySID(ctx context.Context, sid model.ComicChapterLinkSID, data model.SetComicChapterLink, v *model.ComicChapterLink) error {
	if !svc.oauth.HasPermissionContext(ctx, svc.oauth.TokenPermissionKey("write")) {
		return model.GenericError("missing admin permission to update comic chapter link")
	}

	if err := data.Validate(); err != nil {
		return err
	}

	var chapterID any
	switch {
	case sid.ChapterID != nil:
		chapterID = sid.ChapterID
	case sid.ChapterSID != nil:
		chapterID = model.DBComicChapterSIDToID(*sid.ChapterSID)
	}
	var linkID any
	switch {
	case sid.LinkID != nil:
		linkID = sid.LinkID
	case sid.LinkSID != nil:
		linkID = model.DBLinkSIDToID(*sid.LinkSID)
	}
	if err := svc.database.UpdateComicChapterLink(ctx, data, map[string]any{
		model.DBComicChapterGenericChapterID: chapterID,
		model.DBLinkGenericLinkID:            linkID,
	}, v); err != nil {
		return err
	}

	return nil
}

func (svc Service) DeleteComicChapterLinkBySID(ctx context.Context, sid model.ComicChapterLinkSID) error {
	if !svc.oauth.HasPermissionContext(ctx, svc.oauth.TokenPermissionKey("write")) {
		return model.GenericError("missing admin permission to delete comic chapter link")
	}

	var chapterID any
	switch {
	case sid.ChapterID != nil:
		chapterID = sid.ChapterID
	case sid.ChapterSID != nil:
		chapterID = model.DBComicChapterSIDToID(*sid.ChapterSID)
	}
	var linkID any
	switch {
	case sid.LinkID != nil:
		linkID = sid.LinkID
	case sid.LinkSID != nil:
		linkID = model.DBLinkSIDToID(*sid.LinkSID)
	}
	return svc.database.DeleteComicChapterLink(ctx, map[string]any{
		model.DBComicChapterGenericChapterID: chapterID,
		model.DBLinkGenericLinkID:            linkID,
	}, nil)
}

func (svc Service) listComicChapterLink(ctx context.Context, params model.ListParams) ([]*model.ComicChapterLink, error) {
	result, err := svc.database.ListComicChapterLink(ctx, params)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (svc Service) ListComicChapterLink(ctx context.Context, params model.ListParams) ([]*model.ComicChapterLink, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	params.OrderBys = slices.DeleteFunc(params.OrderBys, func(ob model.OrderBy) bool {
		switch field := ob.Field.(type) {
		case string:
			return !slices.Contains(model.ComicChapterLinkOrderByAllow, field)
		}
		return true
	})
	if len(params.OrderBys) > model.ComicChapterLinkOrderBysMax {
		params.OrderBys = params.OrderBys[:model.ComicChapterLinkOrderBysMax]
	}
	if pagination := params.Pagination; pagination != nil {
		if pagination.Limit > model.ComicChapterLinkPaginationMax {
			pagination.Limit = model.ComicChapterLinkPaginationMax
		}
	}

	return svc.listComicChapterLink(ctx, params)
}

func (svc Service) CountComicChapterLink(ctx context.Context, conds any) (int, error) {
	return svc.database.CountComicChapterLink(ctx, conds)
}
