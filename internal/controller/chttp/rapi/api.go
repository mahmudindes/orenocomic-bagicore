package rapi

import (
	"context"

	"github.com/mahmudindes/orenocomic-bagicore/internal/logger"
	"github.com/mahmudindes/orenocomic-bagicore/internal/model"
)

type (
	api struct {
		service Service
		oauth   OAuth
		logger  logger.Logger
	}

	Service interface {
		AddLanguage(ctx context.Context, data model.AddLanguage, v *model.Language) error
		GetLanguageByID(ctx context.Context, id uint) (*model.Language, error)
		GetLanguageByIETF(ctx context.Context, ietf string) (*model.Language, error)
		UpdateLanguageByIETF(ctx context.Context, ietf string, data model.SetLanguage, v *model.Language) error
		DeleteLanguageByIETF(ctx context.Context, ietf string) error
		ListLanguage(ctx context.Context, params model.ListParams) ([]*model.Language, error)
		CountLanguage(ctx context.Context, conds any) (int, error)

		AddWebsite(ctx context.Context, data model.AddWebsite, v *model.Website) error
		GetWebsiteByDomain(ctx context.Context, domain string) (*model.Website, error)
		UpdateWebsiteByDomain(ctx context.Context, domain string, data model.SetWebsite, v *model.Website) error
		DeleteWebsiteByDomain(ctx context.Context, domain string) error
		ListWebsite(ctx context.Context, params model.ListParams) ([]*model.Website, error)
		CountWebsite(ctx context.Context, conds any) (int, error)
		AddWebsiteTLLanguage(ctx context.Context, data model.AddWebsiteTLLanguage, v *model.WebsiteTLLanguage) error
		GetWebsiteTLLanguageBySID(ctx context.Context, sid model.WebsiteTLLanguageSID) (*model.WebsiteTLLanguage, error)
		UpdateWebsiteTLLanguageBySID(ctx context.Context, sid model.WebsiteTLLanguageSID, data model.SetWebsiteTLLanguage, v *model.WebsiteTLLanguage) error
		DeleteWebsiteTLLanguageBySID(ctx context.Context, sid model.WebsiteTLLanguageSID) error

		AddLink(ctx context.Context, data model.AddLink, v *model.Link) error
		GetLinkBySID(ctx context.Context, sid model.LinkSID) (*model.Link, error)
		UpdateLinkBySID(ctx context.Context, sid model.LinkSID, data model.SetLink, v *model.Link) error
		DeleteLinkBySID(ctx context.Context, sid model.LinkSID) error
		ListLink(ctx context.Context, params model.ListParams) ([]*model.Link, error)
		CountLink(ctx context.Context, conds any) (int, error)
		AddLinkTLLanguage(ctx context.Context, data model.AddLinkTLLanguage, v *model.LinkTLLanguage) error
		GetLinkTLLanguageBySID(ctx context.Context, sid model.LinkTLLanguageSID) (*model.LinkTLLanguage, error)
		UpdateLinkTLLanguageBySID(ctx context.Context, sid model.LinkTLLanguageSID, data model.SetLinkTLLanguage, v *model.LinkTLLanguage) error
		DeleteLinkTLLanguageBySID(ctx context.Context, sid model.LinkTLLanguageSID) error

		// Comic
		AddComic(ctx context.Context, data model.AddComic, v *model.Comic) error
		GetComicByCode(ctx context.Context, code string) (*model.Comic, error)
		UpdateComicByCode(ctx context.Context, code string, data model.SetComic, v *model.Comic) error
		DeleteComicByCode(ctx context.Context, code string) error
		ListComic(ctx context.Context, params model.ListParams) ([]*model.Comic, error)
		CountComic(ctx context.Context, conds any) (int, error)
		ExistsComicByCode(ctx context.Context, code string) (bool, error)
		AddComicLink(ctx context.Context, data model.AddComicLink, v *model.ComicLink) error
		GetComicLinkBySID(ctx context.Context, sid model.ComicLinkSID) (*model.ComicLink, error)
		UpdateComicLinkBySID(ctx context.Context, sid model.ComicLinkSID, data model.SetComicLink, v *model.ComicLink) error
		DeleteComicLinkBySID(ctx context.Context, sid model.ComicLinkSID) error
		// Comic Chapter
		AddComicChapter(ctx context.Context, data model.AddComicChapter, v *model.ComicChapter) error
		GetComicChapterBySID(ctx context.Context, sid model.ComicChapterSID) (*model.ComicChapter, error)
		UpdateComicChapterBySID(ctx context.Context, sid model.ComicChapterSID, data model.SetComicChapter, v *model.ComicChapter) error
		DeleteComicChapterBySID(ctx context.Context, sid model.ComicChapterSID) error
		ListComicChapter(ctx context.Context, params model.ListParams) ([]*model.ComicChapter, error)
		CountComicChapter(ctx context.Context, conds any) (int, error)
		ExistsComicChapterBySID(ctx context.Context, sid model.ComicChapterSID) (bool, error)
		AddComicChapterLink(ctx context.Context, data model.AddComicChapterLink, v *model.ComicChapterLink) error
		GetComicChapterLinkBySID(ctx context.Context, sid model.ComicChapterLinkSID) (*model.ComicChapterLink, error)
		UpdateComicChapterLinkBySID(ctx context.Context, sid model.ComicChapterLinkSID, data model.SetComicChapterLink, v *model.ComicChapterLink) error
		DeleteComicChapterLinkBySID(ctx context.Context, sid model.ComicChapterLinkSID) error
	}

	OAuth interface {
		ProcessTokenContext(ctx context.Context) (bool, error)
		IsTokenExpiredError(err error) bool
	}
)

const SecuritySchemeBearerAuth = "BearerAuth"

var _ ServerInterface = (*api)(nil)

func NewAPI(svc Service, oa OAuth, log logger.Logger) *api {
	return &api{service: svc, oauth: oa, logger: log}
}
