package model

import (
	"slices"
	"strconv"
	"time"

	bagicore "github.com/mahmudindes/orenocomic-bagicore"
	"github.com/mahmudindes/orenocomic-bagicore/internal/utila"
)

func init() {
	copy(WebsiteOrderByAllow, GenericOrderByAllow)
}

const (
	WebsiteDomainMax     = 32
	WebsiteNameMax       = 48
	WebsiteOrderBysMax   = 3
	WebsitePaginationDef = 10
	WebsitePaginationMax = 50
	DBWebsite            = bagicore.ID + "." + "website"
	DBWebsiteDomain      = "domain"
	DBWebsiteName        = "name"
	DBWebsiteMachineTL   = "machine_tl"
)

var (
	WebsiteOrderByAllow = []string{
		DBWebsiteDomain,
		DBWebsiteName,
		DBWebsiteMachineTL,
	}

	WebsiteSetNullAllow = []string{
		DBWebsiteMachineTL,
	}

	DBWebsiteDomainToID = func(domain string) DBQueryValue {
		return DBQueryValue{
			Table:      DBWebsite,
			Expression: DBGenericID,
			ZeroValue:  0,
			Conditions: DBConditionalKV{Key: DBWebsiteDomain, Value: domain},
		}
	}
)

type (
	Website struct {
		ID          uint        `json:"id"`
		Domain      string      `json:"domain"`
		Name        string      `json:"name"`
		TLLanguages []*Language `db:"-" json:"tlLanguages"`
		MachineTL   *bool       `json:"machineTL"`
		CreatedAt   time.Time   `json:"createdAt"`
		UpdatedAt   *time.Time  `json:"updatedAt"`
	}

	AddWebsite struct {
		Domain    string
		Name      string
		MachineTL *bool
	}

	SetWebsite struct {
		Domain    *string
		Name      *string
		MachineTL *bool
		SetNull   []string
	}
)

func (m AddWebsite) Validate() error {
	return (SetWebsite{
		Domain:    &m.Domain,
		Name:      &m.Name,
		MachineTL: m.MachineTL,
	}).Validate()
}

func (m SetWebsite) Validate() error {
	if m.Domain != nil {
		if *m.Domain == "" {
			return GenericError("domain cannot be empty")
		}

		if len(*m.Domain) > WebsiteDomainMax {
			max := strconv.FormatInt(WebsiteDomainMax, 10)
			return GenericError("domain must be at most " + max + " characters long")
		}

		if !utila.ValidDomain(*m.Domain) {
			return GenericError("domain is not valid")
		}
	}

	if m.Name != nil {
		if *m.Name == "" {
			return GenericError("name cannot be empty")
		}

		if len(*m.Name) > WebsiteNameMax {
			max := strconv.FormatInt(WebsiteNameMax, 10)
			return GenericError("name must be at most " + max + " characters long")
		}
	}

	for _, key := range m.SetNull {
		if !slices.Contains(WebsiteSetNullAllow, key) {
			return GenericError("set null " + key + " is not recognized")
		}
	}

	return nil
}

func init() {
	copy(WebsiteTLLanguageOrderByAllow, GenericOrderByAllow)
}

const (
	DBWebsiteGenericWebsiteID      = "website_id"
	WebsiteTLLanguageOrderBysMax   = 3
	WebsiteTLLanguagePaginationDef = 10
	WebsiteTLLanguagePaginationMax = 50
	DBWebsiteTLLanguage            = bagicore.ID + "." + "website_tllanguage"
)

var (
	WebsiteTLLanguageOrderByAllow = []string{
		DBLanguageGenericLanguageID,
	}
)

type (
	WebsiteTLLanguage struct {
		WebsiteID    uint       `json:"-"`
		LanguageID   uint       `json:"languageID"`
		LanguageIETF string     `json:"languageIETF"`
		CreatedAt    time.Time  `json:"createdAt"`
		UpdatedAt    *time.Time `json:"updatedAt"`
	}
	AddWebsiteTLLanguage struct {
		WebsiteID     *uint
		WebsiteDomain *string
		LanguageID    *uint
		LanguageIETF  *string
	}
	SetWebsiteTLLanguage struct {
		WebsiteID     *uint
		WebsiteDomain *string
		LanguageID    *uint
		LanguageIETF  *string
	}
	WebsiteTLLanguageSID struct {
		WebsiteID     *uint
		WebsiteDomain *string
		LanguageID    *uint
		LanguageIETF  *string
	}
)

func (m AddWebsiteTLLanguage) Validate() error {
	if m.WebsiteID == nil && m.WebsiteDomain == nil {
		return GenericError("either website id or website domain must exist")
	}

	if m.LanguageID == nil && m.LanguageIETF == nil {
		return GenericError("either language id or language ietf must exist")
	}

	return (&SetWebsiteTLLanguage{
		WebsiteID:     m.WebsiteID,
		WebsiteDomain: m.WebsiteDomain,
		LanguageID:    m.LanguageID,
		LanguageIETF:  m.LanguageIETF,
	}).Validate()
}

func (m SetWebsiteTLLanguage) Validate() error {
	if err := (SetWebsite{Domain: m.WebsiteDomain}).Validate(); err != nil {
		return GenericError("website " + err.Error())
	}

	if err := (SetLanguage{IETF: m.LanguageIETF}).Validate(); err != nil {
		return GenericError("language " + err.Error())
	}

	return nil
}
