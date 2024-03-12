package service

import (
	"context"
	"slices"

	"github.com/mahmudindes/orenocomic-bagicore/internal/model"
)

func (svc Service) AddLanguage(ctx context.Context, data model.AddLanguage, v *model.Language) error {
	if !svc.oauth.HasPermissionContext(ctx, svc.oauth.TokenPermissionKey("write")) {
		return model.GenericError("missing admin permission to add language")
	}

	if err := data.Validate(); err != nil {
		return err
	}

	return svc.database.AddLanguage(ctx, data, v)
}

func (svc Service) GetLanguageByID(ctx context.Context, id uint) (*model.Language, error) {
	return svc.database.GetLanguage(ctx, model.DBConditionalKV{
		Key:   model.DBGenericID,
		Value: id,
	})
}

func (svc Service) GetLanguageByIETF(ctx context.Context, ietf string) (*model.Language, error) {
	return svc.database.GetLanguage(ctx, model.DBConditionalKV{
		Key:   model.DBLanguageIETF,
		Value: ietf,
	})
}

func (svc Service) UpdateLanguageByIETF(ctx context.Context, ietf string, data model.SetLanguage, v *model.Language) error {
	if !svc.oauth.HasPermissionContext(ctx, svc.oauth.TokenPermissionKey("write")) {
		return model.GenericError("missing admin permission to update language")
	}

	if err := data.Validate(); err != nil {
		return err
	}

	return svc.database.UpdateLanguage(ctx, data, model.DBConditionalKV{
		Key:   model.DBLanguageIETF,
		Value: ietf,
	}, v)
}

func (svc Service) DeleteLanguageByIETF(ctx context.Context, ietf string) error {
	if !svc.oauth.HasPermissionContext(ctx, svc.oauth.TokenPermissionKey("write")) {
		return model.GenericError("missing admin permission to delete language")
	}

	return svc.database.DeleteLanguage(ctx, model.DBConditionalKV{
		Key:   model.DBLanguageIETF,
		Value: ietf,
	}, nil)
}

func (svc Service) ListLanguage(ctx context.Context, params model.ListParams) ([]*model.Language, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	params.OrderBys = slices.DeleteFunc(params.OrderBys, func(ob model.OrderBy) bool {
		switch field := ob.Field.(type) {
		case string:
			return !slices.Contains(model.LanguageOrderByAllow, field)
		}
		return true
	})
	if len(params.OrderBys) > model.LanguageOrderBysMax {
		params.OrderBys = params.OrderBys[:model.LanguageOrderBysMax]
	}
	if pagination := params.Pagination; pagination != nil {
		if pagination.Limit > model.LanguagePaginationMax {
			pagination.Limit = model.LanguagePaginationMax
		}
	}

	return svc.database.ListLanguage(ctx, params)
}

func (svc Service) CountLanguage(ctx context.Context, conds any) (int, error) {
	return svc.database.CountLanguage(ctx, conds)
}
