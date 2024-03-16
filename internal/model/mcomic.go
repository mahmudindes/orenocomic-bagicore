package model

import (
	"slices"
	"strconv"
	"time"

	bagicore "github.com/mahmudindes/orenocomic-bagicore"
)

func init() {
	ComicOrderByAllow = append(ComicOrderByAllow, GenericOrderByAllow...)
}

const (
	ComicCodeLength    = 8
	ComicOrderBysMax   = 3
	ComicPaginationDef = 10
	ComicPaginationMax = 50
	DBComic            = bagicore.ID + "." + "comic"
	DBComicCode        = "code"
)

var (
	ComicOrderByAllow = []string{
		DBComicCode,
	}

	DBComicCodeToID = func(code string) DBQueryValue {
		return DBQueryValue{
			Table:      DBComic,
			Expression: DBGenericID,
			ZeroValue:  0,
			Conditions: DBConditionalKV{Key: DBComicCode, Value: code},
		}
	}
)

type (
	Comic struct {
		ID        uint            `json:"id"`
		Code      string          `json:"code"`
		Links     []*Link         `db:"-" json:"links"`
		Chapters  []*ComicChapter `db:"-" json:"chapters"`
		CreatedAt time.Time       `json:"createdAt"`
		UpdatedAt *time.Time      `json:"updatedAt"`
	}

	AddComic struct {
		Code string
	}

	SetComic struct {
		Code *string
	}
)

func (m AddComic) Validate() error {
	return (SetComic{
		Code: &m.Code,
	}).Validate()
}

func (m SetComic) Validate() error {
	if m.Code != nil {
		if *m.Code == "" {
			return GenericError("code cannot be empty")
		}

		if len(*m.Code) != ComicCodeLength {
			length := strconv.Itoa(ComicCodeLength)
			return GenericError("code must be " + length + " characters long")
		}
	}

	return nil
}

func init() {
	ComicLinkOrderByAllow = append(ComicLinkOrderByAllow, GenericOrderByAllow...)
}

const (
	DBComicGenericComicID  = "comic_id"
	ComicLinkOrderBysMax   = 3
	ComicLinkPaginationDef = 10
	ComicLinkPaginationMax = 50
	DBComicLink            = bagicore.ID + "." + "comic_link"
)

var (
	ComicLinkOrderByAllow = []string{
		DBLinkGenericLinkID,
	}
)

type (
	ComicLink struct {
		ComicID           uint       `json:"-"`
		LinkID            uint       `json:"linkID"`
		LinkWebsiteDomain string     `json:"linkWebsiteDomain"`
		LinkRelativeURL   string     `json:"linkRelativeURL"`
		CreatedAt         time.Time  `json:"createdAt"`
		UpdatedAt         *time.Time `json:"updatedAt"`
	}
	AddComicLink struct {
		ComicID   *uint
		ComicCode *string
		LinkID    *uint
		LinkSID   *LinkSID
	}
	SetComicLink struct {
		ComicID   *uint
		ComicCode *string
		LinkID    *uint
		LinkSID   *LinkSID
	}
	ComicLinkSID struct {
		ComicID   *uint
		ComicCode *string
		LinkID    *uint
		LinkSID   *LinkSID
	}
)

func (m AddComicLink) Validate() error {
	if m.ComicID == nil && m.ComicCode == nil {
		return GenericError("either comic id or comic code must exist")
	}

	if m.LinkID == nil && m.LinkSID == nil {
		return GenericError("either link id or link sid must exist")
	}

	return (SetComicLink{
		ComicID:   m.ComicID,
		ComicCode: m.ComicCode,
		LinkID:    m.LinkID,
		LinkSID:   m.LinkSID,
	}).Validate()
}
func (m SetComicLink) Validate() error {
	if err := (SetComic{Code: m.ComicCode}).Validate(); err != nil {
		return GenericError("comic " + err.Error())
	}

	if m.LinkSID != nil {
		if err := (SetLink{
			WebsiteDomain: m.LinkSID.WebsiteDomain,
			RelativeURL:   &m.LinkSID.RelativeURL,
		}).Validate(); err != nil {
			return GenericError("link " + err.Error())
		}
	}

	return nil
}

func init() {
	ComicChapterOrderByAllow = append(ComicChapterOrderByAllow, GenericOrderByAllow...)
}

const (
	ComicChapterChapterMax    = 64
	ComicChapterVersionMax    = 32
	ComicChapterOrderBysMax   = 5
	ComicChapterPaginationDef = 10
	ComicChapterPaginationMax = 50
	DBComicChapter            = bagicore.ID + "." + "comic_chapter"
	DBComicChapterChapter     = "chapter"
	DBComicChapterVersion     = "version"
	DBComicChapterReleasedAt  = "released_at"
)

var (
	ComicChapterOrderByAllow = []string{
		DBComicGenericComicID,
		DBComicChapterChapter,
		DBComicChapterVersion,
		DBComicChapterReleasedAt,
	}

	ComicChapterSetNullAllow = []string{
		DBComicChapterChapter,
		DBComicChapterVersion,
	}

	DBComicChapterSIDToID = func(sid ComicChapterSID) DBQueryValue {
		var comicID any
		switch {
		case sid.ComicID != nil:
			comicID = sid.ComicID
		case sid.ComicCode != nil:
			comicID = DBComicCodeToID(*sid.ComicCode)
		}
		var version any
		switch {
		case sid.Version != nil:
			version = sid.Version
		default:
			version = DBIsNull{}
		}
		return DBQueryValue{
			Table:      DBComicChapter,
			Expression: DBGenericID,
			ZeroValue:  0,
			Conditions: map[string]any{
				DBComicGenericComicID: comicID,
				DBComicChapterChapter: sid.Chapter,
				DBComicChapterVersion: version,
			},
		}
	}
)

type (
	ComicChapter struct {
		ID         uint       `json:"id"`
		ComicID    uint       `json:"comicID"`
		ComicCode  string     `json:"comicCode"`
		Chapter    string     `json:"chapter"`
		Version    *string    `json:"version"`
		ReleasedAt time.Time  `json:"releasedAt"`
		Links      []*Link    `db:"-" json:"links"`
		CreatedAt  time.Time  `json:"createdAt"`
		UpdatedAt  *time.Time `json:"updatedAt"`
	}

	AddComicChapter struct {
		ComicID    *uint
		ComicCode  *string
		Chapter    string
		Version    *string
		ReleasedAt time.Time
	}

	SetComicChapter struct {
		ComicID    *uint
		ComicCode  *string
		Chapter    *string
		Version    *string
		ReleasedAt *time.Time
		SetNull    []string
	}

	ComicChapterSID struct {
		ComicID   *uint
		ComicCode *string
		Chapter   string
		Version   *string
	}
)

func (m AddComicChapter) Validate() error {
	if m.ComicID == nil && m.ComicCode == nil {
		return GenericError("either comic id or comic code must exist")
	}

	return (SetComicChapter{
		ComicID:    m.ComicID,
		ComicCode:  m.ComicCode,
		Chapter:    &m.Chapter,
		Version:    m.Version,
		ReleasedAt: &m.ReleasedAt,
	}).Validate()
}

func (m SetComicChapter) Validate() error {
	if err := (SetComic{Code: m.ComicCode}).Validate(); err != nil {
		return GenericError("comic " + err.Error())
	}

	if m.Chapter != nil {
		if *m.Chapter == "" {
			return GenericError("chapter cannot be empty")
		}

		if len(*m.Chapter) > ComicChapterChapterMax {
			max := strconv.FormatInt(ComicChapterChapterMax, 10)
			return GenericError("chapter must be at most " + max + " characters long")
		}
	}

	if m.Version != nil {
		if *m.Version == "" {
			return GenericError("version cannot be empty")
		}

		if len(*m.Version) > ComicChapterVersionMax {
			max := strconv.FormatInt(ComicChapterVersionMax, 10)
			return GenericError("version must be at most " + max + " characters long")
		}
	}

	for _, key := range m.SetNull {
		if !slices.Contains(ComicChapterSetNullAllow, key) {
			return GenericError("set null " + key + " is not recognized")
		}
	}

	return nil
}

func init() {
	ComicChapterLinkOrderByAllow = append(ComicChapterLinkOrderByAllow, GenericOrderByAllow...)
}

const (
	DBComicChapterGenericChapterID = "chapter_id"
	ComicChapterLinkOrderBysMax    = 3
	ComicChapterLinkPaginationDef  = 10
	ComicChapterLinkPaginationMax  = 50
	DBComicChapterLink             = bagicore.ID + "." + "comic_chapter_link"
)

var (
	ComicChapterLinkOrderByAllow = []string{
		DBWebsiteGenericWebsiteID,
	}
)

type (
	ComicChapterLink struct {
		ChapterID         uint       `json:"-"`
		LinkID            uint       `json:"linkID"`
		LinkWebsiteDomain string     `json:"linkWebsiteDomain"`
		LinkRelativeURL   string     `json:"linkRelativeURL"`
		CreatedAt         time.Time  `json:"createdAt"`
		UpdatedAt         *time.Time `json:"updatedAt"`
	}
	AddComicChapterLink struct {
		ChapterID  *uint
		ChapterSID *ComicChapterSID
		LinkID     *uint
		LinkSID    *LinkSID
	}
	SetComicChapterLink struct {
		ChapterID  *uint
		ChapterSID *ComicChapterSID
		LinkID     *uint
		LinkSID    *LinkSID
	}
	ComicChapterLinkSID struct {
		ChapterID  *uint
		ChapterSID *ComicChapterSID
		LinkID     *uint
		LinkSID    *LinkSID
	}
)

func (m AddComicChapterLink) Validate() error {
	if m.ChapterID == nil && m.ChapterSID == nil {
		return GenericError("either chapter id or chapter sid must not empty")
	}

	if m.LinkID == nil && m.LinkSID == nil {
		return GenericError("link id or link sid must not empty")
	}

	return (SetComicChapterLink{
		ChapterID:  m.ChapterID,
		ChapterSID: m.ChapterSID,
		LinkID:     m.LinkID,
		LinkSID:    m.LinkSID,
	}).Validate()
}
func (m SetComicChapterLink) Validate() error {
	if m.ChapterSID != nil {
		if err := (SetComicChapter{
			ComicCode: m.ChapterSID.ComicCode,
			Chapter:   &m.ChapterSID.Chapter,
			Version:   m.ChapterSID.Version,
		}).Validate(); err != nil {
			return GenericError("comic chapter " + err.Error())
		}
	}

	if m.LinkSID != nil {
		if err := (SetLink{
			WebsiteDomain: m.LinkSID.WebsiteDomain,
			RelativeURL:   &m.LinkSID.RelativeURL,
		}).Validate(); err != nil {
			return GenericError("link " + err.Error())
		}
	}

	return nil
}
