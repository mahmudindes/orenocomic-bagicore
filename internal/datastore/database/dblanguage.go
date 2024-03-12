package database

import (
	"context"
	"errors"

	"github.com/mahmudindes/orenocomic-bagicore/internal/model"
)

const (
	NameErrLanguageKey = "language_ietf_key"
)

func (db Database) AddLanguage(ctx context.Context, data model.AddLanguage, v *model.Language) error {
	if err := db.GenericAdd(ctx, model.DBLanguage, map[string]any{
		model.DBLanguageIETF: data.IETF,
		model.DBLanguageName: data.Name,
	}, v); err != nil {
		return languageSetError(err)
	}
	return nil
}

func (db Database) GetLanguage(ctx context.Context, conds any) (*model.Language, error) {
	var result model.Language
	if err := db.GenericGet(ctx, model.DBLanguage, conds, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (db Database) UpdateLanguage(ctx context.Context, data model.SetLanguage, conds any, v *model.Language) error {
	data0 := map[string]any{}
	if data.IETF != nil {
		data0[model.DBLanguageIETF] = data.IETF
	}
	if data.Name != nil {
		data0[model.DBLanguageName] = data.Name
	}
	if err := db.GenericUpdate(ctx, model.DBLanguage, data0, conds, v); err != nil {
		return languageSetError(err)
	}
	return nil
}

func (db Database) DeleteLanguage(ctx context.Context, conds any, v *model.Language) error {
	return db.GenericDelete(ctx, model.DBLanguage, conds, v)
}

func (db Database) ListLanguage(ctx context.Context, params model.ListParams) ([]*model.Language, error) {
	result := []*model.Language{}
	if len(params.OrderBys) < 1 {
		params.OrderBys = append(params.OrderBys, model.OrderBy{Field: model.DBLanguageIETF})
	}
	if params.Pagination == nil {
		params.Pagination = &model.Pagination{Page: 1, Limit: model.LanguagePaginationDef}
	}
	if err := db.GenericList(ctx, model.DBLanguage, params, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (db Database) CountLanguage(ctx context.Context, conds any) (int, error) {
	return db.GenericCount(ctx, model.DBLanguage, conds)
}

func languageSetError(err error) error {
	var errDatabase model.DatabaseError
	if errors.As(err, &errDatabase) {
		if errDatabase.Code == CodeErrExists && errDatabase.Name == NameErrLanguageKey {
			return model.GenericError("same ietf already exists")
		}
	}
	return err
}
