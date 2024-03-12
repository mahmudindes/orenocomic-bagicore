package service

import (
	"context"
	"slices"

	"github.com/mahmudindes/orenocomic-bagicore/internal/model"
)

func (svc Service) AddLink(ctx context.Context, data model.AddLink, v *model.Link) error {
	if !svc.oauth.HasPermissionContext(ctx, svc.oauth.TokenPermissionKey("write")) {
		return model.GenericError("missing admin permission to add link")
	}

	if err := data.Validate(); err != nil {
		return err
	}

	if v != nil {
		v.TLLanguages = []*model.Language{}
	}

	return svc.database.AddLink(ctx, data, v)
}

func (svc Service) GetLinkBySID(ctx context.Context, sid model.LinkSID) (*model.Link, error) {
	var websiteID any
	switch {
	case sid.WebsiteID != nil:
		websiteID = sid.WebsiteID
	case sid.WebsiteDomain != nil:
		websiteID = model.DBWebsiteDomainToID(*sid.WebsiteDomain)
	}
	result, err := svc.database.GetLink(ctx, map[string]any{
		model.DBWebsiteGenericWebsiteID: websiteID,
		model.DBLinkRelativeURL:         sid.RelativeURL,
	})
	if err != nil {
		return nil, err
	}

	tlLanguages0, err := svc.database.ListLinkTLLanguage(ctx, model.ListParams{
		Conditions: model.DBConditionalKV{Key: model.DBLinkGenericLinkID, Value: result.ID},
		Pagination: &model.Pagination{},
	})
	if err != nil {
		return nil, err
	}
	if len(tlLanguages0) > 0 {
		conditions := make([]any, len(tlLanguages0)+2)
		conditions = append(conditions, model.DBLogicalOR{})
		for _, language := range tlLanguages0 {
			conditions = append(conditions, model.DBConditionalKV{
				Key:   model.DBGenericID,
				Value: language.LanguageID,
			})
		}
		tlLanguages1, err := svc.database.ListLanguage(ctx, model.ListParams{
			Conditions: conditions,
			Pagination: &model.Pagination{},
		})
		if err != nil {
			return nil, err
		}
		result.TLLanguages = tlLanguages1
	}

	return result, nil
}

func (svc Service) UpdateLinkBySID(ctx context.Context, sid model.LinkSID, data model.SetLink, v *model.Link) error {
	if !svc.oauth.HasPermissionContext(ctx, svc.oauth.TokenPermissionKey("write")) {
		return model.GenericError("missing admin permission to update link")
	}

	if err := data.Validate(); err != nil {
		return err
	}

	var websiteID any
	switch {
	case sid.WebsiteID != nil:
		websiteID = sid.WebsiteID
	case sid.WebsiteDomain != nil:
		websiteID = model.DBWebsiteDomainToID(*sid.WebsiteDomain)
	}
	if err := svc.database.UpdateLink(ctx, data, map[string]any{
		model.DBWebsiteGenericWebsiteID: websiteID,
		model.DBLinkRelativeURL:         sid.RelativeURL,
	}, v); err != nil {
		return err
	}

	if v != nil {
		tlLanguages0, err := svc.database.ListLinkTLLanguage(ctx, model.ListParams{
			Conditions: model.DBConditionalKV{Key: model.DBLinkGenericLinkID, Value: v.ID},
			Pagination: &model.Pagination{},
		})
		if err != nil {
			return err
		}
		if len(tlLanguages0) > 0 {
			conditions := make([]any, len(tlLanguages0)+2)
			conditions = append(conditions, model.DBLogicalOR{})
			for _, language := range tlLanguages0 {
				conditions = append(conditions, model.DBConditionalKV{
					Key:   model.DBGenericID,
					Value: language.LanguageID,
				})
			}
			tlLanguages1, err := svc.database.ListLanguage(ctx, model.ListParams{
				Conditions: conditions,
				Pagination: &model.Pagination{},
			})
			if err != nil {
				return err
			}
			v.TLLanguages = tlLanguages1
		}
	}

	return nil
}

func (svc Service) DeleteLinkBySID(ctx context.Context, sid model.LinkSID) error {
	if !svc.oauth.HasPermissionContext(ctx, svc.oauth.TokenPermissionKey("write")) {
		return model.GenericError("missing admin permission to delete link")
	}

	var websiteID any
	switch {
	case sid.WebsiteID != nil:
		websiteID = sid.WebsiteID
	case sid.WebsiteDomain != nil:
		websiteID = model.DBWebsiteDomainToID(*sid.WebsiteDomain)
	}
	return svc.database.DeleteLink(ctx, map[string]any{
		model.DBWebsiteGenericWebsiteID: websiteID,
		model.DBLinkRelativeURL:         sid.RelativeURL,
	}, nil)
}

func (svc Service) ListLink(ctx context.Context, params model.ListParams) ([]*model.Link, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	params.OrderBys = slices.DeleteFunc(params.OrderBys, func(ob model.OrderBy) bool {
		switch field := ob.Field.(type) {
		case string:
			return !slices.Contains(model.LinkOrderByAllow, field)
		}
		return true
	})
	if len(params.OrderBys) > model.LinkOrderBysMax {
		params.OrderBys = params.OrderBys[:model.LinkOrderBysMax]
	}
	if pagination := params.Pagination; pagination != nil {
		if pagination.Limit > model.LinkPaginationMax {
			pagination.Limit = model.LinkPaginationMax
		}
	}

	result, err := svc.database.ListLink(ctx, params)
	if err != nil {
		return nil, err
	}

	if len(result) > 0 {
		conds := make([]any, len(result)+1)
		conds = append(conds, model.DBLogicalOR{})
		for _, r := range result {
			conds = append(conds, model.DBConditionalKV{
				Key:   model.DBLinkGenericLinkID,
				Value: r.ID,
			})
		}
		tlLanguages0, err := svc.database.ListLinkTLLanguage(ctx, model.ListParams{
			Conditions: conds,
			Pagination: &model.Pagination{},
		})
		if err != nil {
			return nil, err
		}
		tlLanguages := map[uint]*model.Language{}
		for _, tlLanguage := range tlLanguages0 {
			tlLanguages[tlLanguage.LanguageID] = nil
		}
		conditions := make([]any, len(tlLanguages0)+2)
		conditions = append(conditions, model.DBLogicalOR{})
		for id := range tlLanguages {
			conditions = append(conditions, model.DBConditionalKV{
				Key:   model.DBGenericID,
				Value: id,
			})
		}
		tlLanguages1, err := svc.database.ListLanguage(ctx, model.ListParams{
			Conditions: conditions,
			Pagination: &model.Pagination{},
		})
		if err != nil {
			return nil, err
		}
		for _, tlLanguage := range tlLanguages1 {
			tlLanguages[tlLanguage.ID] = tlLanguage
		}
		for _, r := range result {
			r.TLLanguages = []*model.Language{}
		}
		for _, tlLanguage := range tlLanguages0 {
			for _, r := range result {
				if r.ID == tlLanguage.LinkID {
					r.TLLanguages = append(r.TLLanguages, tlLanguages[tlLanguage.LanguageID])
				}
			}
		}
	}

	return result, nil
}

func (svc Service) CountLink(ctx context.Context, conds any) (int, error) {
	return svc.database.CountLink(ctx, conds)
}

func (svc Service) AddLinkTLLanguage(ctx context.Context, data model.AddLinkTLLanguage, v *model.LinkTLLanguage) error {
	if !svc.oauth.HasPermissionContext(ctx, svc.oauth.TokenPermissionKey("write")) {
		return model.GenericError("missing admin permission to add link tl language")
	}

	if err := data.Validate(); err != nil {
		return err
	}

	return svc.database.AddLinkTLLanguage(ctx, data, v)
}

func (svc Service) GetLinkTLLanguageBySID(ctx context.Context, sid model.LinkTLLanguageSID) (*model.LinkTLLanguage, error) {
	var websiteID any
	switch {
	case sid.LinkID != nil:
		websiteID = sid.LinkID
	case sid.LinkSID != nil:
		websiteID = model.DBLinkSIDToID(*sid.LinkSID)
	}
	var languageID any
	switch {
	case sid.LanguageID != nil:
		languageID = sid.LanguageID
	case sid.LanguageIETF != nil:
		languageID = model.DBLanguageIETFToID(*sid.LanguageIETF)
	}
	return svc.database.GetLinkTLLanguage(ctx, map[string]any{
		model.DBLinkGenericLinkID:         websiteID,
		model.DBLanguageGenericLanguageID: languageID,
	})
}

func (svc Service) UpdateLinkTLLanguageBySID(ctx context.Context, sid model.LinkTLLanguageSID, data model.SetLinkTLLanguage, v *model.LinkTLLanguage) error {
	if !svc.oauth.HasPermissionContext(ctx, svc.oauth.TokenPermissionKey("write")) {
		return model.GenericError("missing admin permission to update link tl language")
	}

	if err := data.Validate(); err != nil {
		return err
	}

	var websiteID any
	switch {
	case sid.LinkID != nil:
		websiteID = sid.LinkID
	case sid.LinkSID != nil:
		websiteID = model.DBLinkSIDToID(*sid.LinkSID)
	}
	var languageID any
	switch {
	case sid.LanguageID != nil:
		languageID = sid.LanguageID
	case sid.LanguageIETF != nil:
		languageID = model.DBLanguageIETFToID(*sid.LanguageIETF)
	}
	return svc.database.UpdateLinkTLLanguage(ctx, data, map[string]any{
		model.DBLinkGenericLinkID:         websiteID,
		model.DBLanguageGenericLanguageID: languageID,
	}, v)
}

func (svc Service) DeleteLinkTLLanguageBySID(ctx context.Context, sid model.LinkTLLanguageSID) error {
	if !svc.oauth.HasPermissionContext(ctx, svc.oauth.TokenPermissionKey("write")) {
		return model.GenericError("missing admin permission to delete link tl language")
	}

	var websiteID any
	switch {
	case sid.LinkID != nil:
		websiteID = sid.LinkID
	case sid.LinkSID != nil:
		websiteID = model.DBLinkSIDToID(*sid.LinkSID)
	}
	var languageID any
	switch {
	case sid.LanguageID != nil:
		languageID = sid.LanguageID
	case sid.LanguageIETF != nil:
		languageID = model.DBLanguageIETFToID(*sid.LanguageIETF)
	}
	return svc.database.DeleteLinkTLLanguage(ctx, map[string]any{
		model.DBLinkGenericLinkID:         websiteID,
		model.DBLanguageGenericLanguageID: languageID,
	}, nil)
}

func (svc Service) ListLinkTLLanguage(ctx context.Context, params model.ListParams) ([]*model.LinkTLLanguage, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	params.OrderBys = slices.DeleteFunc(params.OrderBys, func(ob model.OrderBy) bool {
		switch field := ob.Field.(type) {
		case string:
			return !slices.Contains(model.LinkTLLanguageOrderByAllow, field)
		}
		return true
	})
	if len(params.OrderBys) > model.LinkTLLanguageOrderBysMax {
		params.OrderBys = params.OrderBys[:model.LinkTLLanguageOrderBysMax]
	}
	if pagination := params.Pagination; pagination != nil {
		if pagination.Limit > model.LinkTLLanguagePaginationMax {
			pagination.Limit = model.LinkTLLanguagePaginationMax
		}
	}

	return svc.database.ListLinkTLLanguage(ctx, params)
}

func (svc Service) CountLinkTLLanguage(ctx context.Context, conds any) (int, error) {
	return svc.database.CountLinkTLLanguage(ctx, conds)
}
