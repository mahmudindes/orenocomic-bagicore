package model

import (
	"slices"
	"strconv"
	"time"

	bagicore "github.com/mahmudindes/orenocomic-bagicore"
)

func init() {
	LinkOrderByAllow = append(LinkOrderByAllow, GenericOrderByAllow...)
}

const (
	LinkCodeLength     = 8
	LinkRelativeURLMax = 128
	LinkOrderBysMax    = 3
	LinkPaginationDef  = 10
	LinkPaginationMax  = 50
	DBLink             = bagicore.ID + "." + "link"
	DBLinkRelativeURL  = "relative_url"
	DBLinkMachineTL    = "machine_tl"
)

var (
	LinkOrderByAllow = []string{
		DBWebsiteGenericWebsiteID,
		DBLinkRelativeURL,
		DBLinkMachineTL,
	}

	LinkSetNullAllow = []string{
		DBLinkMachineTL,
	}

	DBLinkSIDToID = func(sid LinkSID) DBQueryValue {
		var websiteID any
		switch {
		case sid.WebsiteID != nil:
			websiteID = sid.WebsiteID
		case sid.WebsiteDomain != nil:
			websiteID = DBWebsiteDomainToID(*sid.WebsiteDomain)
		}
		return DBQueryValue{
			Table:      DBLink,
			Expression: DBGenericID,
			ZeroValue:  0,
			Conditions: map[string]any{
				DBWebsiteGenericWebsiteID: websiteID,
				DBLinkRelativeURL:         sid.RelativeURL,
			},
		}
	}
)

type (
	Link struct {
		ID            uint        `json:"id"`
		WebsiteID     uint        `json:"websiteID"`
		WebsiteDomain string      `json:"websiteDomain"`
		RelativeURL   string      `json:"relativeURL"`
		TLLanguages   []*Language `db:"-" json:"tlLanguages"`
		MachineTL     *bool       `json:"machineTL"`
		CreatedAt     time.Time   `json:"createdAt"`
		UpdatedAt     *time.Time  `json:"updatedAt"`
	}

	AddLink struct {
		WebsiteID     *uint
		WebsiteDomain *string
		RelativeURL   string
		MachineTL     *bool
	}

	SetLink struct {
		WebsiteID     *uint
		WebsiteDomain *string
		RelativeURL   *string
		MachineTL     *bool
		SetNull       []string
	}

	LinkSID struct {
		WebsiteID     *uint
		WebsiteDomain *string
		RelativeURL   string
	}
)

func (m AddLink) Validate() error {
	return (SetLink{
		WebsiteID:     m.WebsiteID,
		WebsiteDomain: m.WebsiteDomain,
		RelativeURL:   &m.RelativeURL,
		MachineTL:     m.MachineTL,
	}).Validate()
}

func (m SetLink) Validate() error {
	if err := (SetWebsite{Domain: m.WebsiteDomain}).Validate(); err != nil {
		return GenericError("website " + err.Error())
	}

	if m.RelativeURL != nil {
		if *m.RelativeURL == "" {
			return GenericError("relative url cannot be empty")
		}

		if len(*m.RelativeURL) > LinkRelativeURLMax {
			max := strconv.FormatInt(LinkRelativeURLMax, 10)
			return GenericError("relative url must be at most " + max + " characters long")
		}
	}

	for _, key := range m.SetNull {
		if !slices.Contains(LinkSetNullAllow, key) {
			return GenericError("set null " + key + " is not recognized")
		}
	}

	return nil
}

func init() {
	LinkTLLanguageOrderByAllow = append(LinkTLLanguageOrderByAllow, GenericOrderByAllow...)
}

const (
	DBLinkGenericLinkID         = "link_id"
	LinkTLLanguageOrderBysMax   = 3
	LinkTLLanguagePaginationDef = 10
	LinkTLLanguagePaginationMax = 50
	DBLinkTLLanguage            = bagicore.ID + "." + "link_tllanguage"
)

var (
	LinkTLLanguageOrderByAllow = []string{
		DBLanguageGenericLanguageID,
	}
)

type (
	LinkTLLanguage struct {
		LinkID       uint       `json:"-"`
		LanguageID   uint       `json:"languageID"`
		LanguageIETF string     `json:"languageIETF"`
		CreatedAt    time.Time  `json:"createdAt"`
		UpdatedAt    *time.Time `json:"updatedAt"`
	}
	AddLinkTLLanguage struct {
		LinkID       *uint
		LinkSID      *LinkSID
		LanguageID   *uint
		LanguageIETF *string
	}
	SetLinkTLLanguage struct {
		LinkID       *uint
		LinkSID      *LinkSID
		LanguageID   *uint
		LanguageIETF *string
	}
	LinkTLLanguageSID struct {
		LinkID       *uint
		LinkSID      *LinkSID
		LanguageID   *uint
		LanguageIETF *string
	}
)

func (m AddLinkTLLanguage) Validate() error {
	if m.LinkID == nil && m.LinkSID == nil {
		return GenericError("either link id or link sid must exist")
	}

	if m.LanguageID == nil && m.LanguageIETF == nil {
		return GenericError("either language id or language ietf must exist")
	}

	return (&SetLinkTLLanguage{
		LinkID:       m.LinkID,
		LinkSID:      m.LinkSID,
		LanguageID:   m.LanguageID,
		LanguageIETF: m.LanguageIETF,
	}).Validate()
}

func (m SetLinkTLLanguage) Validate() error {
	if m.LinkSID != nil {
		if err := (SetLink{
			WebsiteDomain: m.LinkSID.WebsiteDomain,
			RelativeURL:   &m.LinkSID.RelativeURL,
		}).Validate(); err != nil {
			return GenericError("link " + err.Error())
		}
	}

	if err := (SetLanguage{IETF: m.LanguageIETF}).Validate(); err != nil {
		return GenericError("language " + err.Error())
	}

	return nil
}
