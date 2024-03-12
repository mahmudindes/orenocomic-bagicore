package service

import (
	"context"
	"slices"

	"github.com/mahmudindes/orenocomic-bagicore/internal/model"
)

func (svc Service) AddWebsite(ctx context.Context, data model.AddWebsite, v *model.Website) error {
	if !svc.oauth.HasPermissionContext(ctx, svc.oauth.TokenPermissionKey("write")) {
		return model.GenericError("missing admin permission to add website")
	}

	if err := data.Validate(); err != nil {
		return err
	}

	if v != nil {
		v.TLLanguages = []*model.Language{}
	}

	return svc.database.AddWebsite(ctx, data, v)
}

func (svc Service) GetWebsiteByDomain(ctx context.Context, domain string) (*model.Website, error) {
	result, err := svc.database.GetWebsite(ctx, model.DBConditionalKV{
		Key:   model.DBWebsiteDomain,
		Value: domain,
	})
	if err != nil {
		return nil, err
	}

	tlLanguages0, err := svc.database.ListWebsiteTLLanguage(ctx, model.ListParams{
		Conditions: model.DBConditionalKV{Key: model.DBWebsiteGenericWebsiteID, Value: result.ID},
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

func (svc Service) UpdateWebsiteByDomain(ctx context.Context, domain string, data model.SetWebsite, v *model.Website) error {
	if !svc.oauth.HasPermissionContext(ctx, svc.oauth.TokenPermissionKey("write")) {
		return model.GenericError("missing admin permission to update website")
	}

	if err := data.Validate(); err != nil {
		return err
	}

	if err := svc.database.UpdateWebsite(ctx, data, model.DBConditionalKV{
		Key:   model.DBWebsiteDomain,
		Value: domain,
	}, v); err != nil {
		return err
	}

	if v != nil {
		tlLanguages0, err := svc.database.ListWebsiteTLLanguage(ctx, model.ListParams{
			Conditions: model.DBConditionalKV{Key: model.DBWebsiteGenericWebsiteID, Value: v.ID},
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

func (svc Service) DeleteWebsiteByDomain(ctx context.Context, domain string) error {
	if !svc.oauth.HasPermissionContext(ctx, svc.oauth.TokenPermissionKey("write")) {
		return model.GenericError("missing admin permission to delete website")
	}

	return svc.database.DeleteWebsite(ctx, model.DBConditionalKV{
		Key:   model.DBWebsiteDomain,
		Value: domain,
	}, nil)
}

func (svc Service) ListWebsite(ctx context.Context, params model.ListParams) ([]*model.Website, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	params.OrderBys = slices.DeleteFunc(params.OrderBys, func(ob model.OrderBy) bool {
		switch field := ob.Field.(type) {
		case string:
			return !slices.Contains(model.WebsiteOrderByAllow, field)
		}
		return true
	})
	if len(params.OrderBys) > model.WebsiteOrderBysMax {
		params.OrderBys = params.OrderBys[:model.WebsiteOrderBysMax]
	}
	if pagination := params.Pagination; pagination != nil {
		if pagination.Limit > model.WebsitePaginationMax {
			pagination.Limit = model.WebsitePaginationMax
		}
	}

	result, err := svc.database.ListWebsite(ctx, params)
	if err != nil {
		return nil, err
	}

	if len(result) > 0 {
		conds := make([]any, len(result)+1)
		conds = append(conds, model.DBLogicalOR{})
		for _, r := range result {
			conds = append(conds, model.DBConditionalKV{
				Key:   model.DBWebsiteGenericWebsiteID,
				Value: r.ID,
			})
		}
		tlLanguages0, err := svc.database.ListWebsiteTLLanguage(ctx, model.ListParams{
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
				if r.ID == tlLanguage.WebsiteID {
					r.TLLanguages = append(r.TLLanguages, tlLanguages[tlLanguage.LanguageID])
				}
			}
		}
	}

	return result, nil
}

func (svc Service) CountWebsite(ctx context.Context, conds any) (int, error) {
	return svc.database.CountWebsite(ctx, conds)
}

func (svc Service) AddWebsiteTLLanguage(ctx context.Context, data model.AddWebsiteTLLanguage, v *model.WebsiteTLLanguage) error {
	if !svc.oauth.HasPermissionContext(ctx, svc.oauth.TokenPermissionKey("write")) {
		return model.GenericError("missing admin permission to add website tl language")
	}

	if err := data.Validate(); err != nil {
		return err
	}

	return svc.database.AddWebsiteTLLanguage(ctx, data, v)
}

func (svc Service) GetWebsiteTLLanguageBySID(ctx context.Context, sid model.WebsiteTLLanguageSID) (*model.WebsiteTLLanguage, error) {
	var websiteID any
	switch {
	case sid.WebsiteID != nil:
		websiteID = sid.WebsiteID
	case sid.WebsiteDomain != nil:
		websiteID = model.DBWebsiteDomainToID(*sid.WebsiteDomain)
	}
	var languageID any
	switch {
	case sid.LanguageID != nil:
		languageID = sid.LanguageID
	case sid.LanguageIETF != nil:
		languageID = model.DBLanguageIETFToID(*sid.LanguageIETF)
	}
	return svc.database.GetWebsiteTLLanguage(ctx, map[string]any{
		model.DBWebsiteGenericWebsiteID:   websiteID,
		model.DBLanguageGenericLanguageID: languageID,
	})
}

func (svc Service) UpdateWebsiteTLLanguageBySID(ctx context.Context, sid model.WebsiteTLLanguageSID, data model.SetWebsiteTLLanguage, v *model.WebsiteTLLanguage) error {
	if !svc.oauth.HasPermissionContext(ctx, svc.oauth.TokenPermissionKey("write")) {
		return model.GenericError("missing admin permission to update website tl language")
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
	var languageID any
	switch {
	case sid.LanguageID != nil:
		languageID = sid.LanguageID
	case sid.LanguageIETF != nil:
		languageID = model.DBLanguageIETFToID(*sid.LanguageIETF)
	}
	return svc.database.UpdateWebsiteTLLanguage(ctx, data, map[string]any{
		model.DBWebsiteGenericWebsiteID:   websiteID,
		model.DBLanguageGenericLanguageID: languageID,
	}, v)
}

func (svc Service) DeleteWebsiteTLLanguageBySID(ctx context.Context, sid model.WebsiteTLLanguageSID) error {
	if !svc.oauth.HasPermissionContext(ctx, svc.oauth.TokenPermissionKey("write")) {
		return model.GenericError("missing admin permission to delete website tl language")
	}

	var websiteID any
	switch {
	case sid.WebsiteID != nil:
		websiteID = sid.WebsiteID
	case sid.WebsiteDomain != nil:
		websiteID = model.DBWebsiteDomainToID(*sid.WebsiteDomain)
	}
	var languageID any
	switch {
	case sid.LanguageID != nil:
		languageID = sid.LanguageID
	case sid.LanguageIETF != nil:
		languageID = model.DBLanguageIETFToID(*sid.LanguageIETF)
	}
	return svc.database.DeleteWebsiteTLLanguage(ctx, map[string]any{
		model.DBWebsiteGenericWebsiteID:   websiteID,
		model.DBLanguageGenericLanguageID: languageID,
	}, nil)
}

func (svc Service) ListWebsiteTLLanguage(ctx context.Context, params model.ListParams) ([]*model.WebsiteTLLanguage, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	params.OrderBys = slices.DeleteFunc(params.OrderBys, func(ob model.OrderBy) bool {
		switch field := ob.Field.(type) {
		case string:
			return !slices.Contains(model.WebsiteTLLanguageOrderByAllow, field)
		}
		return true
	})
	if len(params.OrderBys) > model.WebsiteTLLanguageOrderBysMax {
		params.OrderBys = params.OrderBys[:model.WebsiteTLLanguageOrderBysMax]
	}
	if pagination := params.Pagination; pagination != nil {
		if pagination.Limit > model.WebsiteTLLanguagePaginationMax {
			pagination.Limit = model.WebsiteTLLanguagePaginationMax
		}
	}

	return svc.database.ListWebsiteTLLanguage(ctx, params)
}

func (svc Service) CountWebsiteTLLanguage(ctx context.Context, conds any) (int, error) {
	return svc.database.CountWebsiteTLLanguage(ctx, conds)
}
