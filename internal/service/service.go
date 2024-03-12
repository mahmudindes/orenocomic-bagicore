package service

import (
	"context"

	"github.com/mahmudindes/orenocomic-bagicore/internal/model"
)

type (
	Service struct {
		database database
		oauth    oauth
	}

	database interface {
		AddLanguage(ctx context.Context, data model.AddLanguage, v *model.Language) error
		GetLanguage(ctx context.Context, conds any) (*model.Language, error)
		UpdateLanguage(ctx context.Context, data model.SetLanguage, conds any, v *model.Language) error
		DeleteLanguage(ctx context.Context, conds any, v *model.Language) error
		ListLanguage(ctx context.Context, params model.ListParams) ([]*model.Language, error)
		CountLanguage(ctx context.Context, conds any) (int, error)

		AddWebsite(ctx context.Context, data model.AddWebsite, v *model.Website) error
		GetWebsite(ctx context.Context, conds any) (*model.Website, error)
		UpdateWebsite(ctx context.Context, data model.SetWebsite, conds any, v *model.Website) error
		DeleteWebsite(ctx context.Context, conds any, v *model.Website) error
		ListWebsite(ctx context.Context, params model.ListParams) ([]*model.Website, error)
		CountWebsite(ctx context.Context, conds any) (int, error)
		ExistsWebsite(ctx context.Context, conds any) (bool, error)
		AddWebsiteTLLanguage(ctx context.Context, data model.AddWebsiteTLLanguage, v *model.WebsiteTLLanguage) error
		GetWebsiteTLLanguage(ctx context.Context, conds any) (*model.WebsiteTLLanguage, error)
		UpdateWebsiteTLLanguage(ctx context.Context, data model.SetWebsiteTLLanguage, conds any, v *model.WebsiteTLLanguage) error
		DeleteWebsiteTLLanguage(ctx context.Context, conds any, v *model.WebsiteTLLanguage) error
		ListWebsiteTLLanguage(ctx context.Context, params model.ListParams) ([]*model.WebsiteTLLanguage, error)
		CountWebsiteTLLanguage(ctx context.Context, conds any) (int, error)

		AddLink(ctx context.Context, data model.AddLink, v *model.Link) error
		GetLink(ctx context.Context, conds any) (*model.Link, error)
		UpdateLink(ctx context.Context, data model.SetLink, conds any, v *model.Link) error
		DeleteLink(ctx context.Context, conds any, v *model.Link) error
		ListLink(ctx context.Context, params model.ListParams) ([]*model.Link, error)
		CountLink(ctx context.Context, conds any) (int, error)
		ExistsLink(ctx context.Context, conds any) (bool, error)
		AddLinkTLLanguage(ctx context.Context, data model.AddLinkTLLanguage, v *model.LinkTLLanguage) error
		GetLinkTLLanguage(ctx context.Context, conds any) (*model.LinkTLLanguage, error)
		UpdateLinkTLLanguage(ctx context.Context, data model.SetLinkTLLanguage, conds any, v *model.LinkTLLanguage) error
		DeleteLinkTLLanguage(ctx context.Context, conds any, v *model.LinkTLLanguage) error
		ListLinkTLLanguage(ctx context.Context, params model.ListParams) ([]*model.LinkTLLanguage, error)
		CountLinkTLLanguage(ctx context.Context, conds any) (int, error)

		// Comic
		AddComic(ctx context.Context, data model.AddComic, v *model.Comic) error
		GetComic(ctx context.Context, conds any) (*model.Comic, error)
		UpdateComic(ctx context.Context, data model.SetComic, conds any, v *model.Comic) error
		DeleteComic(ctx context.Context, conds any, v *model.Comic) error
		ListComic(ctx context.Context, params model.ListParams) ([]*model.Comic, error)
		CountComic(ctx context.Context, conds any) (int, error)
		ExistsComic(ctx context.Context, conds any) (bool, error)
		AddComicLink(ctx context.Context, data model.AddComicLink, v *model.ComicLink) error
		GetComicLink(ctx context.Context, conds any) (*model.ComicLink, error)
		UpdateComicLink(ctx context.Context, data model.SetComicLink, conds any, v *model.ComicLink) error
		DeleteComicLink(ctx context.Context, params any, v *model.ComicLink) error
		ListComicLink(ctx context.Context, params model.ListParams) ([]*model.ComicLink, error)
		CountComicLink(ctx context.Context, conds any) (int, error)
		// Comic Chapter
		AddComicChapter(ctx context.Context, data model.AddComicChapter, v *model.ComicChapter) error
		GetComicChapter(ctx context.Context, conds any) (*model.ComicChapter, error)
		UpdateComicChapter(ctx context.Context, data model.SetComicChapter, conds any, v *model.ComicChapter) error
		DeleteComicChapter(ctx context.Context, conds any, v *model.ComicChapter) error
		ListComicChapter(ctx context.Context, params model.ListParams) ([]*model.ComicChapter, error)
		CountComicChapter(ctx context.Context, conds any) (int, error)
		ExistsComicChapter(ctx context.Context, conds any) (bool, error)
		AddComicChapterLink(ctx context.Context, data model.AddComicChapterLink, v *model.ComicChapterLink) error
		GetComicChapterLink(ctx context.Context, conds any) (*model.ComicChapterLink, error)
		UpdateComicChapterLink(ctx context.Context, data model.SetComicChapterLink, conds any, v *model.ComicChapterLink) error
		DeleteComicChapterLink(ctx context.Context, conds any, v *model.ComicChapterLink) error
		ListComicChapterLink(ctx context.Context, params model.ListParams) ([]*model.ComicChapterLink, error)
		CountComicChapterLink(ctx context.Context, conds any) (int, error)
	}

	oauth interface {
		HasPermissionContext(ctx context.Context, permission string) bool
		TokenPermissionKey(s ...string) string
	}
)

func New(db database, oa oauth) Service {
	return Service{database: db, oauth: oa}
}
