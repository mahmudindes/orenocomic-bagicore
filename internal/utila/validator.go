package utila

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var (
	Validate    = validator.New(validator.WithRequiredStructEnabled())
	ValidDomain = regexp.MustCompile(`^(?:[0-9\p{L}](?:[0-9\p{L}-]{0,61}[0-9\p{L}])?\.)+[0-9\p{L}][0-9\p{L}-]{0,61}[0-9\p{L}]$`).MatchString
)
