package rapi

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	middleware "github.com/oapi-codegen/nethttp-middleware"

	"github.com/mahmudindes/orenocomic-bagicore/internal/model"
	"github.com/mahmudindes/orenocomic-bagicore/internal/utila"
)

func (api *api) Authentication(ctx context.Context, input *openapi3filter.AuthenticationInput) error {
	switch input.SecuritySchemeName {
	case SecuritySchemeBearerAuth:
		valid, err := api.oauth.ProcessTokenContext(ctx)
		switch {
		case api.oauth.IsTokenExpiredError(err):
			return errors.New("bearer authentication token expired")
		case errors.As(err, &model.ErrGeneric):
			return fmt.Errorf("bearer authentication failed: %w", err)
		case err != nil:
			api.logger.ErrMessage(err, "Bearer authentication proccess token context failed.")
			return errors.New("bearer authentication failed")
		case !valid:
			return errors.New("bearer authentication invalid")
		}
		return nil
	}
	return fmt.Errorf("security scheme %s is not supported", input.SecuritySchemeName)
}

func Middleware(s *openapi3.T, af openapi3filter.AuthenticationFunc) func(http.Handler) http.Handler {
	return middleware.OapiRequestValidatorWithOptions(s, &middleware.Options{
		ErrorHandler:          middlewareErrorHandler,
		SilenceServersWarning: true,
		Options:               openapi3filter.Options{AuthenticationFunc: af},
	})
}

func middlewareErrorHandler(w http.ResponseWriter, message string, statusCode int) {
	responseErr(w, utila.CapitalPeriod(message), statusCode)
}
