package rapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"

	"github.com/mahmudindes/orenocomic-bagicore/internal/model"
)

func modelLink(m *model.Link) Link {
	return Link{
		ID:            m.ID,
		WebsiteID:     m.WebsiteID,
		WebsiteDomain: m.WebsiteDomain,
		RelativeURL:   m.RelativeURL,
		TLLanguages:   slicesModel(m.TLLanguages, modelLanguage),
		MachineTL:     m.MachineTL,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
	}
}

func (api *api) AddLink(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := api.logger.WithContext(ctx)

	var data model.AddLink
	switch r.Header.Get("Content-Type") {
	case "application/json":
		var data0 AddLinkJSONRequestBody
		if err := json.NewDecoder(r.Body).Decode(&data0); err != nil {
			responseErr(w, "Bad request body.", http.StatusBadRequest)
			log.ErrMessage(err, "Add link decode json body failed.")
			return
		}
		data = model.AddLink{
			WebsiteID:     data0.WebsiteID,
			WebsiteDomain: data0.WebsiteDomain,
			RelativeURL:   data0.RelativeURL,
			MachineTL:     data0.MachineTL,
		}
	case "application/x-www-form-urlencoded":
		if err := r.ParseForm(); err != nil {
			responseErr(w, "Bad request body.", http.StatusBadRequest)
			log.ErrMessage(err, "Add link parse form failed.")
			return
		}
		var data0 AddLinkFormdataRequestBody
		if err := formDecode(r.PostForm, &data0); err != nil {
			responseErr(w, "Bad form data.", http.StatusBadRequest)
			log.ErrMessage(err, "Add link decode form data failed.")
			return
		}
		data = model.AddLink{
			WebsiteID:     data0.WebsiteID,
			WebsiteDomain: data0.WebsiteDomain,
			RelativeURL:   data0.RelativeURL,
			MachineTL:     data0.MachineTL,
		}
	}

	result := new(model.Link)
	if err := api.service.AddLink(ctx, data, result); err != nil {
		responseServiceErr(w, err)
		log.ErrMessage(err, "Add link failed.")
		return
	}

	w.Header().Set("Location", r.URL.Path+"/"+result.WebsiteDomain+"-"+url.QueryEscape(result.RelativeURL))
	response(w, modelLink(result), http.StatusCreated)
}

func (api *api) GetLink(w http.ResponseWriter, r *http.Request, websiteDomain string, relativeURL string) {
	ctx := r.Context()
	log := api.logger.WithContext(ctx)

	relativeURL, err := url.QueryUnescape(relativeURL)
	if err != nil {
		responseErr(w, "Invalid link relative url.", http.StatusBadRequest)
		return
	}

	result, err := api.service.GetLinkBySID(ctx, model.LinkSID{
		WebsiteDomain: &websiteDomain,
		RelativeURL:   relativeURL,
	})
	if err != nil {
		responseServiceErr(w, err)
		log.ErrMessage(err, "Get link failed.")
		return
	}

	response(w, modelLink(result), http.StatusOK)
}

func (api *api) UpdateLink(w http.ResponseWriter, r *http.Request, websiteDomain string, relativeURL string) {
	ctx := r.Context()
	log := api.logger.WithContext(ctx)

	relativeURL, err := url.QueryUnescape(relativeURL)
	if err != nil {
		responseErr(w, "Invalid link relative url.", http.StatusBadRequest)
		return
	}

	var data model.SetLink
	switch r.Header.Get("Content-Type") {
	case "application/json":
		var data0 UpdateLinkJSONRequestBody
		if err := json.NewDecoder(r.Body).Decode(&data0); err != nil {
			responseErr(w, "Bad request body.", http.StatusBadRequest)
			log.ErrMessage(err, "Update link decode json body failed.")
			return
		}
		data = model.SetLink{
			WebsiteID:     data0.WebsiteID,
			WebsiteDomain: data0.WebsiteDomain,
			RelativeURL:   data0.RelativeURL,
			MachineTL:     data0.MachineTL,
			SetNull:       data0.SetNull,
		}
	case "application/x-www-form-urlencoded":
		if err := r.ParseForm(); err != nil {
			responseErr(w, "Bad request body.", http.StatusBadRequest)
			log.ErrMessage(err, "Update link parse form failed.")
			return
		}
		var data0 UpdateLinkFormdataRequestBody
		if err := formDecode(r.PostForm, &data0); err != nil {
			responseErr(w, "Bad form data.", http.StatusBadRequest)
			log.ErrMessage(err, "Update link decode form data failed.")
			return
		}
		data = model.SetLink{
			WebsiteID:     data0.WebsiteID,
			WebsiteDomain: data0.WebsiteDomain,
			RelativeURL:   data0.RelativeURL,
			MachineTL:     data0.MachineTL,
			SetNull:       data0.SetNull,
		}
	}

	result := new(model.Link)
	if err := api.service.UpdateLinkBySID(ctx, model.LinkSID{
		WebsiteDomain: &websiteDomain,
		RelativeURL:   relativeURL,
	}, data, result); err != nil {
		if errors.As(err, &model.ErrNotFound) {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		responseServiceErr(w, err)
		log.ErrMessage(err, "Update link failed.")
		return
	}

	w.Header().Set("Location", r.URL.Path+"/"+result.WebsiteDomain+"-"+url.QueryEscape(result.RelativeURL))
	response(w, modelLink(result), http.StatusOK)
}

func (api *api) DeleteLink(w http.ResponseWriter, r *http.Request, websiteDomain string, relativeURL string) {
	ctx := r.Context()
	log := api.logger.WithContext(ctx)

	relativeURL, err := url.QueryUnescape(relativeURL)
	if err != nil {
		responseErr(w, "Invalid link relative url.", http.StatusBadRequest)
		return
	}

	if err := api.service.DeleteLinkBySID(ctx, model.LinkSID{
		WebsiteDomain: &websiteDomain,
		RelativeURL:   relativeURL,
	}); err != nil {
		responseServiceErr(w, err)
		log.ErrMessage(err, "Delete link failed.")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (api *api) ListLink(w http.ResponseWriter, r *http.Request, params ListLinkParams) {
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
		count, err := api.service.CountLink(ctx, nil)
		if err != nil {
			totalCountCh <- -1
			log.ErrMessage(err, "Count link failed.")
			return
		}
		totalCountCh <- count
	}()

	result0, err := api.service.ListLink(ctx, model.ListParams{
		OrderBys:   orderBys,
		Pagination: &pagination,
	})
	if err != nil {
		responseServiceErr(w, err)
		log.ErrMessage(err, "List link failed.")
		return
	}

	wHeader := w.Header()
	wHeader.Set("X-Total-Count", strconv.Itoa(<-totalCountCh))
	wHeader.Set("X-Pagination-Limit", strconv.Itoa(pagination.Limit))
	var result []Link
	for _, r := range result0 {
		result = append(result, modelLink(r))
	}
	response(w, result, http.StatusOK)
}

func modelLinkTLLanguage(m *model.LinkTLLanguage) LinkTLLanguage {
	return LinkTLLanguage{
		LanguageID:   m.LanguageID,
		LanguageIETF: m.LanguageIETF,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}

func (api *api) AddLinkTLLanguage(w http.ResponseWriter, r *http.Request, websiteDomain string, relativeURL string) {
	ctx := r.Context()
	log := api.logger.WithContext(ctx)

	relativeURL, err := url.QueryUnescape(relativeURL)
	if err != nil {
		responseErr(w, "Invalid link relative url.", http.StatusBadRequest)
		return
	}

	var data model.AddLinkTLLanguage
	switch r.Header.Get("Content-Type") {
	case "application/json":
		var data0 AddLinkTLLanguageJSONRequestBody
		if err := json.NewDecoder(r.Body).Decode(&data0); err != nil {
			responseErr(w, "Bad request body.", http.StatusBadRequest)
			log.ErrMessage(err, "Add link tl language decode json body failed.")
			return
		}
		data = model.AddLinkTLLanguage{
			LinkID:       nil,
			LinkSID:      &model.LinkSID{WebsiteDomain: &websiteDomain, RelativeURL: relativeURL},
			LanguageID:   data0.LanguageID,
			LanguageIETF: data0.LanguageIETF,
		}
	case "application/x-www-form-urlencoded":
		if err := r.ParseForm(); err != nil {
			responseErr(w, "Bad request body.", http.StatusBadRequest)
			log.ErrMessage(err, "Add link tl language parse form failed.")
			return
		}
		var data0 AddWebsiteTLLanguageFormdataRequestBody
		if err := formDecode(r.PostForm, &data0); err != nil {
			responseErr(w, "Bad form data.", http.StatusBadRequest)
			log.ErrMessage(err, "Add link tl language decode form data failed.")
			return
		}
		data = model.AddLinkTLLanguage{
			LinkID:       nil,
			LinkSID:      &model.LinkSID{WebsiteDomain: &websiteDomain, RelativeURL: relativeURL},
			LanguageID:   data0.LanguageID,
			LanguageIETF: data0.LanguageIETF,
		}
	}

	result := new(model.LinkTLLanguage)
	if err := api.service.AddLinkTLLanguage(ctx, data, result); err != nil {
		responseServiceErr(w, err)
		log.ErrMessage(err, "Add link tl language failed.")
		return
	}

	w.Header().Set("Location", r.URL.Path+"/"+result.LanguageIETF)
	response(w, modelLinkTLLanguage(result), http.StatusCreated)
}

func (api *api) GetLinkTLLanguage(w http.ResponseWriter, r *http.Request, websiteDomain string, relativeURL string, ietf string) {
	ctx := r.Context()
	log := api.logger.WithContext(ctx)

	relativeURL, err := url.QueryUnescape(relativeURL)
	if err != nil {
		responseErr(w, "Invalid link relative url.", http.StatusBadRequest)
		return
	}

	result, err := api.service.GetLinkTLLanguageBySID(ctx, model.LinkTLLanguageSID{
		LinkSID:      &model.LinkSID{WebsiteDomain: &websiteDomain, RelativeURL: relativeURL},
		LanguageIETF: &ietf,
	})
	if err != nil {
		responseServiceErr(w, err)
		log.ErrMessage(err, "Get link tl language failed.")
		return
	}

	response(w, modelLinkTLLanguage(result), http.StatusOK)
}

func (api *api) UpdateLinkTLLanguage(w http.ResponseWriter, r *http.Request, websiteDomain string, relativeURL string, ietf string) {
	ctx := r.Context()
	log := api.logger.WithContext(ctx)

	relativeURL, err := url.QueryUnescape(relativeURL)
	if err != nil {
		responseErr(w, "Invalid link relative url.", http.StatusBadRequest)
		return
	}

	var data model.SetLinkTLLanguage
	switch r.Header.Get("Content-Type") {
	case "application/json":
		var data0 UpdateLinkTLLanguageJSONRequestBody
		if err := json.NewDecoder(r.Body).Decode(&data0); err != nil {
			responseErr(w, "Bad request body.", http.StatusBadRequest)
			log.ErrMessage(err, "Update link tl language decode json body failed.")
			return
		}
		data = model.SetLinkTLLanguage{
			LinkID:       nil,
			LinkSID:      nil,
			LanguageID:   data0.LanguageID,
			LanguageIETF: data0.LanguageIETF,
		}
	case "application/x-www-form-urlencoded":
		if err := r.ParseForm(); err != nil {
			responseErr(w, "Bad request body.", http.StatusBadRequest)
			log.ErrMessage(err, "Update link tl language parse form failed.")
			return
		}
		var data0 UpdateLinkTLLanguageFormdataRequestBody
		if err := formDecode(r.PostForm, &data0); err != nil {
			responseErr(w, "Bad form data.", http.StatusBadRequest)
			log.ErrMessage(err, "Update link tl language decode form data failed.")
			return
		}
		data = model.SetLinkTLLanguage{
			LinkID:       nil,
			LinkSID:      nil,
			LanguageID:   data0.LanguageID,
			LanguageIETF: data0.LanguageIETF,
		}
	}

	result := new(model.LinkTLLanguage)
	if err := api.service.UpdateLinkTLLanguageBySID(ctx, model.LinkTLLanguageSID{
		LinkSID:      &model.LinkSID{WebsiteDomain: &websiteDomain, RelativeURL: relativeURL},
		LanguageIETF: &ietf,
	}, data, result); err != nil {
		if errors.As(err, &model.ErrNotFound) {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		responseServiceErr(w, err)
		log.ErrMessage(err, "Update link tl language failed.")
		return
	}

	w.Header().Set("Location", r.URL.Path+"/"+result.LanguageIETF)
	response(w, modelLinkTLLanguage(result), http.StatusOK)
}

func (api *api) DeleteLinkTLLanguage(w http.ResponseWriter, r *http.Request, websiteDomain string, relativeURL string, ietf string) {
	ctx := r.Context()
	log := api.logger.WithContext(ctx)

	relativeURL, err := url.QueryUnescape(relativeURL)
	if err != nil {
		responseErr(w, "Invalid link relative url.", http.StatusBadRequest)
		return
	}

	if err := api.service.DeleteLinkTLLanguageBySID(ctx, model.LinkTLLanguageSID{
		LinkSID:      &model.LinkSID{WebsiteDomain: &websiteDomain, RelativeURL: relativeURL},
		LanguageIETF: &ietf,
	}); err != nil {
		responseServiceErr(w, err)
		log.ErrMessage(err, "Delete link tl language failed.")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
