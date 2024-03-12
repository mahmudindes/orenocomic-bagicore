package rapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/mahmudindes/orenocomic-bagicore/internal/model"
)

//
// Comic
//

func modelComic(m *model.Comic) Comic {
	return Comic{
		ID:        m.ID,
		Code:      m.Code,
		Links:     slicesModel(m.Links, modelLink),
		Chapters:  slicesModel(m.Chapters, modelComicChapter),
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func (api *api) AddComic(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := api.logger.WithContext(ctx)

	var data model.AddComic
	switch r.Header.Get("Content-Type") {
	case "application/json":
		var data0 AddComicJSONRequestBody
		if err := json.NewDecoder(r.Body).Decode(&data0); err != nil {
			responseErr(w, "Bad request body.", http.StatusBadRequest)
			log.ErrMessage(err, "Add comic decode json body failed.")
			return
		}
		data = model.AddComic{
			Code: data0.Code,
		}
	case "application/x-www-form-urlencoded":
		if err := r.ParseForm(); err != nil {
			responseErr(w, "Bad request body.", http.StatusBadRequest)
			log.ErrMessage(err, "Add comic parse form failed.")
			return
		}
		var data0 AddComicFormdataRequestBody
		if err := formDecode(r.PostForm, &data0); err != nil {
			responseErr(w, "Bad form data.", http.StatusBadRequest)
			log.ErrMessage(err, "Add comic decode form data failed.")
			return
		}
		data = model.AddComic{
			Code: data0.Code,
		}
	}

	result := new(model.Comic)
	if err := api.service.AddComic(ctx, data, result); err != nil {
		responseServiceErr(w, err)
		log.ErrMessage(err, "Add comic failed.")
		return
	}

	w.Header().Set("Location", r.URL.Path+"/"+result.Code)
	response(w, modelComic(result), http.StatusCreated)
}

func (api *api) GetComic(w http.ResponseWriter, r *http.Request, code string) {
	ctx := r.Context()
	log := api.logger.WithContext(ctx)

	result, err := api.service.GetComicByCode(ctx, code)
	if err != nil {
		responseServiceErr(w, err)
		log.ErrMessage(err, "Get comic failed.")
		return
	}

	response(w, modelComic(result), http.StatusOK)
}

func (api *api) UpdateComic(w http.ResponseWriter, r *http.Request, code string) {
	ctx := r.Context()
	log := api.logger.WithContext(ctx)

	var data model.SetComic
	switch r.Header.Get("Content-Type") {
	case "application/json":
		var data0 UpdateComicJSONRequestBody
		if err := json.NewDecoder(r.Body).Decode(&data0); err != nil {
			responseErr(w, "Bad request body.", http.StatusBadRequest)
			log.ErrMessage(err, "Update comic decode json body failed.")
			return
		}
		data = model.SetComic{
			Code: data0.Code,
		}
	case "application/x-www-form-urlencoded":
		if err := r.ParseForm(); err != nil {
			responseErr(w, "Bad request body.", http.StatusBadRequest)
			log.ErrMessage(err, "Update comic parse form failed.")
			return
		}
		var data0 UpdateComicFormdataRequestBody
		if err := formDecode(r.PostForm, &data0); err != nil {
			responseErr(w, "Bad form data.", http.StatusBadRequest)
			log.ErrMessage(err, "Update comic decode form data failed.")
			return
		}
		data = model.SetComic{
			Code: data0.Code,
		}
	}

	result := new(model.Comic)
	if err := api.service.UpdateComicByCode(ctx, code, data, result); err != nil {
		if errors.As(err, &model.ErrNotFound) {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		responseServiceErr(w, err)
		log.ErrMessage(err, "Update comic failed.")
		return
	}

	w.Header().Set("Location", r.URL.Path+"/"+result.Code)
	response(w, modelComic(result), http.StatusOK)
}

func (api *api) DeleteComic(w http.ResponseWriter, r *http.Request, code string) {
	ctx := r.Context()
	log := api.logger.WithContext(ctx)

	if err := api.service.DeleteComicByCode(ctx, code); err != nil {
		responseServiceErr(w, err)
		log.ErrMessage(err, "Delete comic failed.")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (api *api) ListComic(w http.ResponseWriter, r *http.Request, params ListComicParams) {
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
		count, err := api.service.CountComic(ctx, nil)
		if err != nil {
			totalCountCh <- -1
			log.ErrMessage(err, "Count comic failed.")
			return
		}
		totalCountCh <- count
	}()

	result0, err := api.service.ListComic(ctx, model.ListParams{
		OrderBys:   orderBys,
		Pagination: &pagination,
	})
	if err != nil {
		responseServiceErr(w, err)
		log.ErrMessage(err, "List comic failed.")
		return
	}

	wHeader := w.Header()
	wHeader.Set("X-Total-Count", strconv.Itoa(<-totalCountCh))
	wHeader.Set("X-Pagination-Limit", strconv.Itoa(pagination.Limit))
	var result []Comic
	for _, r := range result0 {
		result = append(result, modelComic(r))
	}
	response(w, result, http.StatusOK)
}

// Comic Link

func modelComicLink(m *model.ComicLink) ComicLink {
	return ComicLink{
		LinkID:            m.LinkID,
		LinkWebsiteDomain: m.LinkWebsiteDomain,
		LinkRelativeURL:   m.LinkRelativeURL,
		CreatedAt:         m.CreatedAt,
		UpdatedAt:         m.UpdatedAt,
	}
}

func (api *api) AddComicLink(w http.ResponseWriter, r *http.Request, code string) {
	ctx := r.Context()
	log := api.logger.WithContext(ctx)

	var data model.AddComicLink
	switch r.Header.Get("Content-Type") {
	case "application/json":
		var data0 AddComicLinkJSONRequestBody
		if err := json.NewDecoder(r.Body).Decode(&data0); err != nil {
			responseErr(w, "Bad request body.", http.StatusBadRequest)
			log.ErrMessage(err, "Add comic link decode json body failed.")
			return
		}
		data = model.AddComicLink{
			ComicID:   nil,
			ComicCode: &code,
			LinkID:    data0.LinkID,
		}
		if data0.LinkWebsiteDomain != nil && data0.LinkRelativeURL != nil {
			relativeURL, err := url.QueryUnescape(*data0.LinkRelativeURL)
			if err != nil {
				responseErr(w, "Invalid link relative url.", http.StatusBadRequest)
				return
			}
			data.LinkSID = &model.LinkSID{
				WebsiteDomain: data0.LinkWebsiteDomain,
				RelativeURL:   relativeURL,
			}
		}
	case "application/x-www-form-urlencoded":
		if err := r.ParseForm(); err != nil {
			responseErr(w, "Bad request body.", http.StatusBadRequest)
			log.ErrMessage(err, "Add comic link parse form failed.")
			return
		}
		var data0 AddComicLinkFormdataRequestBody
		if err := formDecode(r.PostForm, &data0); err != nil {
			responseErr(w, "Bad form data.", http.StatusBadRequest)
			log.ErrMessage(err, "Add comic link decode form data failed.")
			return
		}
		data = model.AddComicLink{
			ComicID:   nil,
			ComicCode: &code,
			LinkID:    data0.LinkID,
		}
		if data0.LinkWebsiteDomain != nil && data0.LinkRelativeURL != nil {
			relativeURL, err := url.QueryUnescape(*data0.LinkRelativeURL)
			if err != nil {
				responseErr(w, "Invalid link relative url.", http.StatusBadRequest)
				return
			}
			data.LinkSID = &model.LinkSID{
				WebsiteDomain: data0.LinkWebsiteDomain,
				RelativeURL:   relativeURL,
			}
		}
	}

	result := new(model.ComicLink)
	if err := api.service.AddComicLink(ctx, data, result); err != nil {
		responseServiceErr(w, err)
		log.ErrMessage(err, "Add comic link failed.")
		return
	}

	w.Header().Set("Location", r.URL.Path+"/"+result.LinkWebsiteDomain+"-"+url.QueryEscape(result.LinkRelativeURL))
	response(w, modelComicLink(result), http.StatusCreated)
}

func (api *api) GetComicLink(w http.ResponseWriter, r *http.Request, code string, websiteDomain string, relativeURL string) {
	ctx := r.Context()
	log := api.logger.WithContext(ctx)

	relativeURL, err := url.QueryUnescape(relativeURL)
	if err != nil {
		responseErr(w, "Invalid link relative url.", http.StatusBadRequest)
		return
	}

	result, err := api.service.GetComicLinkBySID(ctx, model.ComicLinkSID{
		ComicCode: &code,
		LinkSID:   &model.LinkSID{WebsiteDomain: &websiteDomain, RelativeURL: relativeURL},
	})
	if err != nil {
		responseServiceErr(w, err)
		log.ErrMessage(err, "Get comic link failed.")
		return
	}

	response(w, modelComicLink(result), http.StatusOK)
}

func (api *api) UpdateComicLink(w http.ResponseWriter, r *http.Request, code string, websiteDomain string, relativeURL string) {
	ctx := r.Context()
	log := api.logger.WithContext(ctx)

	relativeURL, err := url.QueryUnescape(relativeURL)
	if err != nil {
		responseErr(w, "Invalid link relative url.", http.StatusBadRequest)
		return
	}

	var data model.SetComicLink
	switch r.Header.Get("Content-Type") {
	case "application/json":
		var data0 UpdateComicLinkJSONRequestBody
		if err := json.NewDecoder(r.Body).Decode(&data0); err != nil {
			responseErr(w, "Bad request body.", http.StatusBadRequest)
			log.ErrMessage(err, "Update comic link decode json body failed.")
			return
		}
		data = model.SetComicLink{
			ComicID:   nil,
			ComicCode: nil,
			LinkID:    data0.LinkID,
		}
		if data0.LinkWebsiteDomain != nil && data0.LinkRelativeURL != nil {
			data.LinkSID = &model.LinkSID{
				WebsiteDomain: data0.LinkWebsiteDomain,
				RelativeURL:   *data0.LinkRelativeURL,
			}
		}
	case "application/x-www-form-urlencoded":
		if err := r.ParseForm(); err != nil {
			responseErr(w, "Bad request body.", http.StatusBadRequest)
			log.ErrMessage(err, "Update comic link parse form failed.")
			return
		}
		var data0 UpdateComicLinkFormdataRequestBody
		if err := formDecode(r.PostForm, &data0); err != nil {
			responseErr(w, "Bad form data.", http.StatusBadRequest)
			log.ErrMessage(err, "Update comic link decode form data failed.")
			return
		}
		data = model.SetComicLink{
			ComicID:   nil,
			ComicCode: nil,
			LinkID:    data0.LinkID,
		}
		if data0.LinkWebsiteDomain != nil && data0.LinkRelativeURL != nil {
			data.LinkSID = &model.LinkSID{
				WebsiteDomain: data0.LinkWebsiteDomain,
				RelativeURL:   *data0.LinkRelativeURL,
			}
		}
	}

	result := new(model.ComicLink)
	if err := api.service.UpdateComicLinkBySID(ctx, model.ComicLinkSID{
		ComicCode: &code,
		LinkSID:   &model.LinkSID{WebsiteDomain: &websiteDomain, RelativeURL: relativeURL},
	}, data, result); err != nil {
		if errors.As(err, &model.ErrNotFound) {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		responseServiceErr(w, err)
		log.ErrMessage(err, "Update comic link failed.")
		return
	}

	w.Header().Set("Location", r.URL.Path+"/"+result.LinkWebsiteDomain+"-"+url.QueryEscape(result.LinkRelativeURL))
	response(w, modelComicLink(result), http.StatusOK)
}

func (api *api) DeleteComicLink(w http.ResponseWriter, r *http.Request, code string, websiteDomain string, relativeURL string) {
	ctx := r.Context()
	log := api.logger.WithContext(ctx)

	relativeURL, err := url.QueryUnescape(relativeURL)
	if err != nil {
		responseErr(w, "Invalid link relative url.", http.StatusBadRequest)
		return
	}

	if err := api.service.DeleteComicLinkBySID(ctx, model.ComicLinkSID{
		ComicCode: &code,
		LinkSID:   &model.LinkSID{WebsiteDomain: &websiteDomain, RelativeURL: relativeURL},
	}); err != nil {
		responseServiceErr(w, err)
		log.ErrMessage(err, "Delete comic link failed.")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

//
// Comic Chapter
//

func modelComicChapter(m *model.ComicChapter) ComicChapter {
	return ComicChapter{
		ID:         m.ID,
		Chapter:    m.Chapter,
		Version:    m.Version,
		ReleasedAt: m.ReleasedAt,
		Links:      slicesModel(m.Links, modelLink),
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
	}
}

func (api *api) AddComicChapter(w http.ResponseWriter, r *http.Request, code string) {
	ctx := r.Context()
	log := api.logger.WithContext(ctx)

	var data model.AddComicChapter
	switch r.Header.Get("Content-Type") {
	case "application/json":
		var data0 AddComicChapterJSONRequestBody
		if err := json.NewDecoder(r.Body).Decode(&data0); err != nil {
			responseErr(w, "Bad request body.", http.StatusBadRequest)
			log.ErrMessage(err, "Add comic chapter decode json body failed.")
			return
		}
		data = model.AddComicChapter{
			ComicID:    nil,
			ComicCode:  &code,
			Chapter:    data0.Chapter,
			Version:    data0.Version,
			ReleasedAt: data0.ReleasedAt,
		}
	case "application/x-www-form-urlencoded":
		if err := r.ParseForm(); err != nil {
			responseErr(w, "Bad request body.", http.StatusBadRequest)
			log.ErrMessage(err, "Add comic chapter parse form failed.")
			return
		}
		var data0 AddComicChapterFormdataRequestBody
		if err := formDecode(r.PostForm, &data0); err != nil {
			responseErr(w, "Bad form data.", http.StatusBadRequest)
			log.ErrMessage(err, "Add comic chapter decode form data failed.")
			return
		}
		data = model.AddComicChapter{
			ComicID:    nil,
			ComicCode:  &code,
			Chapter:    data0.Chapter,
			Version:    data0.Version,
			ReleasedAt: data0.ReleasedAt,
		}
	}

	result := new(model.ComicChapter)
	if err := api.service.AddComicChapter(ctx, data, result); err != nil {
		responseServiceErr(w, err)
		log.ErrMessage(err, "Add comic chapter failed.")
		return
	}

	slug := url.QueryEscape(result.Chapter)
	if result.Version != nil {
		slug += "+" + url.QueryEscape(*result.Version)
	}

	w.Header().Set("Location", r.URL.Path+"/"+slug)
	response(w, modelComicChapter(result), http.StatusCreated)
}

func (api *api) GetComicChapter(w http.ResponseWriter, r *http.Request, code string, cv string) {
	ctx := r.Context()
	log := api.logger.WithContext(ctx)

	chapterRaw, versionRaw, versionOK := strings.Cut(cv, "+")
	var version *string
	if versionOK {
		version = &versionRaw
	}
	chapter, err := url.QueryUnescape(chapterRaw)
	if err != nil {
		responseErr(w, "Invalid comic chapter chapter.", http.StatusBadRequest)
		return
	}

	result, err := api.service.GetComicChapterBySID(ctx, model.ComicChapterSID{
		ComicCode: &code,
		Chapter:   chapter,
		Version:   version,
	})
	if err != nil {
		responseServiceErr(w, err)
		log.ErrMessage(err, "Get comic chapter failed.")
		return
	}

	response(w, modelComicChapter(result), http.StatusOK)
}

func (api *api) UpdateComicChapter(w http.ResponseWriter, r *http.Request, code string, cv string) {
	ctx := r.Context()
	log := api.logger.WithContext(ctx)

	var data model.SetComicChapter
	switch r.Header.Get("Content-Type") {
	case "application/json":
		var data0 UpdateComicChapterJSONRequestBody
		if err := json.NewDecoder(r.Body).Decode(&data0); err != nil {
			responseErr(w, "Bad request body.", http.StatusBadRequest)
			log.ErrMessage(err, "Update comic chapter decode json body failed.")
			return
		}
		data = model.SetComicChapter{
			ComicID:    nil,
			ComicCode:  nil,
			Chapter:    data0.Chapter,
			Version:    data0.Version,
			ReleasedAt: data0.ReleasedAt,
			SetNull:    data0.SetNull,
		}
	case "application/x-www-form-urlencoded":
		if err := r.ParseForm(); err != nil {
			responseErr(w, "Bad request body.", http.StatusBadRequest)
			log.ErrMessage(err, "Update comic chapter parse form failed.")
			return
		}
		var data0 UpdateComicChapterFormdataRequestBody
		if err := formDecode(r.PostForm, &data0); err != nil {
			responseErr(w, "Bad form data.", http.StatusBadRequest)
			log.ErrMessage(err, "Update comic chapter decode form data failed.")
			return
		}
		data = model.SetComicChapter{
			ComicID:    nil,
			ComicCode:  nil,
			Chapter:    data0.Chapter,
			Version:    data0.Version,
			ReleasedAt: data0.ReleasedAt,
			SetNull:    data0.SetNull,
		}
	}

	chapterRaw, versionRaw, versionOK := strings.Cut(cv, "+")
	var version *string
	if versionOK {
		version = &versionRaw
	}
	chapter, err := url.QueryUnescape(chapterRaw)
	if err != nil {
		responseErr(w, "Invalid comic chapter chapter.", http.StatusBadRequest)
		return
	}

	result := new(model.ComicChapter)
	if err := api.service.UpdateComicChapterBySID(ctx, model.ComicChapterSID{
		ComicCode: &code,
		Chapter:   chapter,
		Version:   version,
	}, data, result); err != nil {
		if errors.As(err, &model.ErrNotFound) {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		responseServiceErr(w, err)
		log.ErrMessage(err, "Update comic chapter failed.")
		return
	}

	slug := url.QueryEscape(result.Chapter)
	if result.Version != nil {
		slug += "+" + url.QueryEscape(*result.Version)
	}

	w.Header().Set("Location", r.URL.Path+"/"+slug)
	response(w, modelComicChapter(result), http.StatusOK)
}

func (api *api) DeleteComicChapter(w http.ResponseWriter, r *http.Request, code string, cv string) {
	ctx := r.Context()
	log := api.logger.WithContext(ctx)

	chapterRaw, versionRaw, versionOK := strings.Cut(cv, "+")
	var version *string
	if versionOK {
		version = &versionRaw
	}
	chapter, err := url.QueryUnescape(chapterRaw)
	if err != nil {
		responseErr(w, "Invalid comic chapter chapter.", http.StatusBadRequest)
		return
	}

	if err := api.service.DeleteComicChapterBySID(ctx, model.ComicChapterSID{
		ComicCode: &code,
		Chapter:   chapter,
		Version:   version,
	}); err != nil {
		responseServiceErr(w, err)
		log.ErrMessage(err, "Delete comic chapter failed.")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (api *api) ListComicChapter(w http.ResponseWriter, r *http.Request, code string, params ListComicChapterParams) {
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

	conditions := model.DBConditionalKV{
		Key:   model.DBComicGenericComicID,
		Value: model.DBComicCodeToID(code),
	}

	totalCountCh := make(chan int, 1)
	go func() {
		count, err := api.service.CountComicChapter(ctx, conditions)
		if err != nil {
			totalCountCh <- -1
			log.ErrMessage(err, "Count comic chapter failed.")
			return
		}
		totalCountCh <- count
	}()

	result0, err := api.service.ListComicChapter(ctx, model.ListParams{
		Conditions: conditions,
		OrderBys:   orderBys,
		Pagination: &pagination,
	})
	if err != nil {
		responseServiceErr(w, err)
		log.ErrMessage(err, "List comic chapter failed.")
		return
	}

	wHeader := w.Header()
	wHeader.Set("X-Total-Count", strconv.Itoa(<-totalCountCh))
	wHeader.Set("X-Pagination-Limit", strconv.Itoa(pagination.Limit))
	var result []ComicChapter
	for _, r := range result0 {
		result = append(result, modelComicChapter(r))
	}
	response(w, result, http.StatusOK)
}

// Comic Chapter Link

func modelComicChapterLink(m *model.ComicChapterLink) ComicChapterLink {
	return ComicChapterLink{
		LinkID:            m.LinkID,
		LinkWebsiteDomain: m.LinkWebsiteDomain,
		LinkRelativeURL:   m.LinkRelativeURL,
		CreatedAt:         m.CreatedAt,
		UpdatedAt:         m.UpdatedAt,
	}
}

func (api *api) AddComicChapterLink(w http.ResponseWriter, r *http.Request, code string, cv string) {
	ctx := r.Context()
	log := api.logger.WithContext(ctx)

	chapterRaw, versionRaw, versionOK := strings.Cut(cv, "+")
	var version *string
	if versionOK {
		version = &versionRaw
	}
	chapter, err := url.QueryUnescape(chapterRaw)
	if err != nil {
		responseErr(w, "Invalid comic chapter chapter.", http.StatusBadRequest)
		return
	}

	var data model.AddComicChapterLink
	switch r.Header.Get("Content-Type") {
	case "application/json":
		var data0 AddComicChapterLinkJSONRequestBody
		if err := json.NewDecoder(r.Body).Decode(&data0); err != nil {
			responseErr(w, "Bad request body.", http.StatusBadRequest)
			log.ErrMessage(err, "Add comic chapter link decode json body failed.")
			return
		}
		data = model.AddComicChapterLink{
			ChapterID:  nil,
			ChapterSID: &model.ComicChapterSID{ComicCode: &code, Chapter: chapter, Version: version},
			LinkID:     data0.LinkID,
		}
		if data0.LinkWebsiteDomain != nil && data0.LinkRelativeURL != nil {
			data.LinkSID = &model.LinkSID{
				WebsiteDomain: data0.LinkWebsiteDomain,
				RelativeURL:   *data0.LinkRelativeURL,
			}
		}
	case "application/x-www-form-urlencoded":
		if err := r.ParseForm(); err != nil {
			responseErr(w, "Bad request body.", http.StatusBadRequest)
			log.ErrMessage(err, "Add comic chapterlink parse form failed.")
			return
		}
		var data0 AddComicChapterLinkFormdataRequestBody
		if err := formDecode(r.PostForm, &data0); err != nil {
			responseErr(w, "Bad form data.", http.StatusBadRequest)
			log.ErrMessage(err, "Add comic chapter link decode form data failed.")
			return
		}
		data = model.AddComicChapterLink{
			ChapterID:  nil,
			ChapterSID: &model.ComicChapterSID{ComicCode: &code, Chapter: chapter, Version: version},
			LinkID:     data0.LinkID,
		}
		if data0.LinkWebsiteDomain != nil && data0.LinkRelativeURL != nil {
			data.LinkSID = &model.LinkSID{
				WebsiteDomain: data0.LinkWebsiteDomain,
				RelativeURL:   *data0.LinkRelativeURL,
			}
		}
	}

	result := new(model.ComicChapterLink)
	if err := api.service.AddComicChapterLink(ctx, data, result); err != nil {
		responseServiceErr(w, err)
		log.ErrMessage(err, "Add comic chapter link failed.")
		return
	}

	w.Header().Set("Location", r.URL.Path+"/"+result.LinkWebsiteDomain+"-"+url.QueryEscape(result.LinkRelativeURL))
	response(w, modelComicChapterLink(result), http.StatusCreated)
}

func (api *api) GetComicChapterLink(w http.ResponseWriter, r *http.Request, code string, cv string, websiteDomain string, relativeURL string) {
	ctx := r.Context()
	log := api.logger.WithContext(ctx)

	relativeURL, err := url.QueryUnescape(relativeURL)
	if err != nil {
		responseErr(w, "Invalid link relative url.", http.StatusBadRequest)
		return
	}

	chapterRaw, versionRaw, versionOK := strings.Cut(cv, "+")
	var version *string
	if versionOK {
		version = &versionRaw
	}
	chapter, err := url.QueryUnescape(chapterRaw)
	if err != nil {
		responseErr(w, "Invalid comic chapter chapter.", http.StatusBadRequest)
		return
	}

	result, err := api.service.GetComicChapterLinkBySID(ctx, model.ComicChapterLinkSID{
		ChapterSID: &model.ComicChapterSID{ComicCode: &code, Chapter: chapter, Version: version},
		LinkSID:    &model.LinkSID{WebsiteDomain: &websiteDomain, RelativeURL: relativeURL},
	})
	if err != nil {
		responseServiceErr(w, err)
		log.ErrMessage(err, "Get comic chapter link failed.")
		return
	}

	response(w, modelComicChapterLink(result), http.StatusOK)
}

func (api *api) UpdateComicChapterLink(w http.ResponseWriter, r *http.Request, code string, cv string, websiteDomain string, relativeURL string) {
	ctx := r.Context()
	log := api.logger.WithContext(ctx)

	relativeURL, err := url.QueryUnescape(relativeURL)
	if err != nil {
		responseErr(w, "Invalid link relative url.", http.StatusBadRequest)
		return
	}

	chapterRaw, versionRaw, versionOK := strings.Cut(cv, "+")
	var version *string
	if versionOK {
		version = &versionRaw
	}
	chapter, err := url.QueryUnescape(chapterRaw)
	if err != nil {
		responseErr(w, "Invalid comic chapter chapter.", http.StatusBadRequest)
		return
	}

	var data model.SetComicChapterLink
	switch r.Header.Get("Content-Type") {
	case "application/json":
		var data0 UpdateComicChapterLinkJSONRequestBody
		if err := json.NewDecoder(r.Body).Decode(&data0); err != nil {
			responseErr(w, "Bad request body.", http.StatusBadRequest)
			log.ErrMessage(err, "Update comic chapter link decode json body failed.")
			return
		}
		data = model.SetComicChapterLink{
			ChapterID:  nil,
			ChapterSID: nil,
			LinkID:     data0.LinkID,
		}
		if data0.LinkWebsiteDomain != nil && data0.LinkRelativeURL != nil {
			data.LinkSID = &model.LinkSID{
				WebsiteDomain: data0.LinkWebsiteDomain,
				RelativeURL:   *data0.LinkRelativeURL,
			}
		}
	case "application/x-www-form-urlencoded":
		if err := r.ParseForm(); err != nil {
			responseErr(w, "Bad request body.", http.StatusBadRequest)
			log.ErrMessage(err, "Update comic chapter link parse form failed.")
			return
		}
		var data0 UpdateComicChapterLinkFormdataRequestBody
		if err := formDecode(r.PostForm, &data0); err != nil {
			responseErr(w, "Bad form data.", http.StatusBadRequest)
			log.ErrMessage(err, "Update comic chapter link decode form data failed.")
			return
		}
		data = model.SetComicChapterLink{
			ChapterID:  nil,
			ChapterSID: nil,
			LinkID:     data0.LinkID,
		}
		if data0.LinkWebsiteDomain != nil && data0.LinkRelativeURL != nil {
			data.LinkSID = &model.LinkSID{
				WebsiteDomain: data0.LinkWebsiteDomain,
				RelativeURL:   *data0.LinkRelativeURL,
			}
		}
	}

	result := new(model.ComicChapterLink)
	if err := api.service.UpdateComicChapterLinkBySID(ctx, model.ComicChapterLinkSID{
		ChapterSID: &model.ComicChapterSID{ComicCode: &code, Chapter: chapter, Version: version},
		LinkSID:    &model.LinkSID{WebsiteDomain: &websiteDomain, RelativeURL: relativeURL},
	}, data, result); err != nil {
		if errors.As(err, &model.ErrNotFound) {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		responseServiceErr(w, err)
		log.ErrMessage(err, "Update comic chapter link failed.")
		return
	}

	w.Header().Set("Location", r.URL.Path+"/"+result.LinkWebsiteDomain+"-"+url.QueryEscape(result.LinkRelativeURL))
	response(w, modelComicChapterLink(result), http.StatusOK)
}

func (api *api) DeleteComicChapterLink(w http.ResponseWriter, r *http.Request, code string, cv string, websiteDomain string, relativeURL string) {
	ctx := r.Context()
	log := api.logger.WithContext(ctx)

	relativeURL, err := url.QueryUnescape(relativeURL)
	if err != nil {
		responseErr(w, "Invalid link relative url.", http.StatusBadRequest)
		return
	}

	chapterRaw, versionRaw, versionOK := strings.Cut(cv, "+")
	var version *string
	if versionOK {
		version = &versionRaw
	}
	chapter, err := url.QueryUnescape(chapterRaw)
	if err != nil {
		responseErr(w, "Invalid comic chapter chapter.", http.StatusBadRequest)
		return
	}

	if err := api.service.DeleteComicChapterLinkBySID(ctx, model.ComicChapterLinkSID{
		ChapterSID: &model.ComicChapterSID{ComicCode: &code, Chapter: chapter, Version: version},
		LinkSID:    &model.LinkSID{WebsiteDomain: &websiteDomain, RelativeURL: relativeURL},
	}); err != nil {
		responseServiceErr(w, err)
		log.ErrMessage(err, "Delete comic chapter link failed.")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
