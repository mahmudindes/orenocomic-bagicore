package oauth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"

	"github.com/mahmudindes/orenocomic-bagicore/internal/logger"
)

type (
	OAuth struct {
		issuer           string
		audience         string
		jwks             jwk.Set
		permissionPrefix string
		httpClient       *http.Client
		cTokenCache      cTokenCacheStore
		logger           logger.Logger
	}

	Config struct {
		Issuer           string `conf:"issuer"`
		Audience         string `conf:"audience"`
		PermissionPrefix string `conf:"permission_prefix"`
	}

	Redis interface {
		GobGet(ctx context.Context, key string, v any) error
		GobSet(ctx context.Context, key string, v any, exp time.Duration) error
	}
)

func New(ctx context.Context, rdb Redis, cfg Config, log logger.Logger) (*OAuth, error) {
	data := struct {
		Issuer  string `json:"issuer"`
		JWKSURI string `json:"jwks_uri"`
	}{}

	client := &http.Client{Timeout: 15 * time.Second}

	oauthMetadata := cfg.Issuer + ".well-known/oauth-authorization-server"
	reqOM, err := http.NewRequestWithContext(ctx, http.MethodGet, oauthMetadata, nil)
	if err != nil {
		return nil, err
	}
	resOM, err := client.Do(reqOM)
	if err != nil {
		return nil, err
	}
	defer resOM.Body.Close()

	if resOM.StatusCode < 400 {
		if err := json.NewDecoder(resOM.Body).Decode(&data); err != nil {
			return nil, err
		}
	}

	if resOM.StatusCode >= 400 {
		oidcDiscovery := cfg.Issuer + ".well-known/openid-configuration"
		reqOD, err := http.NewRequestWithContext(ctx, http.MethodGet, oidcDiscovery, nil)
		if err != nil {
			return nil, err
		}
		resOD, err := client.Do(reqOD)
		if err != nil {
			return nil, err
		}
		defer resOD.Body.Close()

		if err := json.NewDecoder(resOD.Body).Decode(&data); err != nil {
			return nil, err
		}
	}

	if data.Issuer != cfg.Issuer {
		return nil, fmt.Errorf("issuer did not match, expected %q got %q", cfg.Issuer, data.Issuer)
	}

	jwkc := jwk.NewCache(ctx)
	if err := jwkc.Register(data.JWKSURI, jwk.WithHTTPClient(client)); err != nil {
		return nil, err
	}

	jwks, err := jwkc.Refresh(ctx, data.JWKSURI)
	if err != nil {
		return nil, err
	}

	return &OAuth{
		issuer:           cfg.Issuer,
		audience:         cfg.Audience,
		jwks:             jwks,
		permissionPrefix: cfg.PermissionPrefix,
		httpClient:       client,
		cTokenCache:      newCTokenCache(rdb),
		logger:           log,
	}, nil
}

func (oa OAuth) TokenPermissionKey(s ...string) string {
	permission := oa.permissionPrefix
	for _, key := range s {
		permission += "." + key
	}
	return permission
}

func (oa OAuth) IsTokenExpiredError(err error) bool {
	return errors.Is(err, jwt.ErrTokenExpired())
}

func (oa OAuth) IsTokenValidationError(err error) bool {
	return jwt.IsValidationError(err)
}
