package rapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/mahmudindes/orenocomic-bagicore/internal/model"
)

func modelLanguage(m *model.Language) Language {
	return Language{
		ID:        m.ID,
		IETF:      m.IETF,
		Name:      m.Name,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func (api *api) AddLanguage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := api.logger.WithContext(ctx)

	var data model.AddLanguage
	switch r.Header.Get("Content-Type") {
	case "application/json":
		var data0 AddLanguageJSONRequestBody
		if err := json.NewDecoder(r.Body).Decode(&data0); err != nil {
			responseErr(w, "Bad request body.", http.StatusBadRequest)
			log.ErrMessage(err, "Add language decode json body failed.")
			return
		}
		data = model.AddLanguage{
			IETF: data0.IETF,
			Name: data0.Name,
		}
	case "application/x-www-form-urlencoded":
		if err := r.ParseForm(); err != nil {
			responseErr(w, "Bad request body.", http.StatusBadRequest)
			log.ErrMessage(err, "Add language parse form failed.")
			return
		}
		var data0 AddLanguageFormdataRequestBody
		if err := formDecode(r.PostForm, &data0); err != nil {
			responseErr(w, "Bad form data.", http.StatusBadRequest)
			log.ErrMessage(err, "Add language decode form data failed.")
			return
		}
		data = model.AddLanguage{
			IETF: data0.IETF,
			Name: data0.Name,
		}
	}

	result := new(model.Language)
	if err := api.service.AddLanguage(ctx, data, result); err != nil {
		responseServiceErr(w, err)
		log.ErrMessage(err, "Add language failed.")
		return
	}

	w.Header().Set("Location", r.URL.Path+"/"+result.IETF)
	response(w, modelLanguage(result), http.StatusCreated)
}

func (api *api) GetLanguage(w http.ResponseWriter, r *http.Request, ietf string) {
	ctx := r.Context()
	log := api.logger.WithContext(ctx)

	result, err := api.service.GetLanguageByIETF(ctx, ietf)
	if err != nil {
		responseServiceErr(w, err)
		log.ErrMessage(err, "Get language failed.")
		return
	}

	response(w, modelLanguage(result), http.StatusOK)
}

func (api *api) UpdateLanguage(w http.ResponseWriter, r *http.Request, ietf string) {
	ctx := r.Context()
	log := api.logger.WithContext(ctx)

	var data model.SetLanguage
	switch r.Header.Get("Content-Type") {
	case "application/json":
		var data0 UpdateLanguageJSONRequestBody
		if err := json.NewDecoder(r.Body).Decode(&data0); err != nil {
			responseErr(w, "Bad request body.", http.StatusBadRequest)
			log.ErrMessage(err, "Update language decode json body failed.")
			return
		}
		data = model.SetLanguage{
			IETF: data0.IETF,
			Name: data0.Name,
		}
	case "application/x-www-form-urlencoded":
		if err := r.ParseForm(); err != nil {
			responseErr(w, "Bad request body.", http.StatusBadRequest)
			log.ErrMessage(err, "Update language parse form failed.")
			return
		}
		var data0 UpdateLanguageFormdataRequestBody
		if err := formDecode(r.PostForm, &data0); err != nil {
			responseErr(w, "Bad form data.", http.StatusBadRequest)
			log.ErrMessage(err, "Update language decode form data failed.")
			return
		}
		data = model.SetLanguage{
			IETF: data0.IETF,
			Name: data0.Name,
		}
	}

	result := new(model.Language)
	if err := api.service.UpdateLanguageByIETF(ctx, ietf, data, result); err != nil {
		if errors.As(err, &model.ErrNotFound) {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		responseServiceErr(w, err)
		log.ErrMessage(err, "Update language failed.")
		return
	}

	w.Header().Set("Location", r.URL.Path+"/"+result.IETF)
	response(w, modelLanguage(result), http.StatusOK)
}

func (api *api) DeleteLanguage(w http.ResponseWriter, r *http.Request, ietf string) {
	ctx := r.Context()
	log := api.logger.WithContext(ctx)

	if err := api.service.DeleteLanguageByIETF(ctx, ietf); err != nil {
		responseServiceErr(w, err)
		log.ErrMessage(err, "Delete language failed.")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (api *api) ListLanguage(w http.ResponseWriter, r *http.Request, params ListLanguageParams) {
	ctx := r.Context()
	log := api.logger.WithContext(ctx)

	pagination := model.Pagination{Page: 1, Limit: 10}
	if params.Page != nil {
		pagination.Page = *params.Page
	}
	if params.Limit != nil {
		pagination.Limit = *params.Limit
	}

	var orderBys model.OrderBys
	if params.OrderBy != nil {
		orderBys = queryOrderBys(*params.OrderBy)
	}

	totalCountCh := make(chan int, 1)
	go func() {
		count, err := api.service.CountLanguage(ctx, nil)
		if err != nil {
			totalCountCh <- -1
			log.ErrMessage(err, "Count language failed.")
			return
		}
		totalCountCh <- count
	}()

	result0, err := api.service.ListLanguage(ctx, model.ListParams{
		OrderBys:   orderBys,
		Pagination: &pagination,
	})
	if err != nil {
		responseServiceErr(w, err)
		log.ErrMessage(err, "List language failed.")
		return
	}

	wHeader := w.Header()
	wHeader.Set("X-Total-Count", strconv.Itoa(<-totalCountCh))
	wHeader.Set("X-Pagination-Limit", strconv.Itoa(pagination.Limit))
	var result []Language
	for _, r := range result0 {
		result = append(result, modelLanguage(r))
	}
	response(w, result, http.StatusOK)
}
