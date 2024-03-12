package oauth

import (
	"context"
	"slices"
	"strings"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwt"

	"github.com/mahmudindes/orenocomic-bagicore/internal/model"
)

type accessToken struct {
	Subject    string
	Expiration time.Time
	Others     map[string]any
}

func (at accessToken) Claim(name string) (any, bool) {
	switch name {
	case jwt.SubjectKey:
		return at.Subject, true
	case jwt.ExpirationKey:
		return at.Expiration, true
	default:
		v, ok := at.Others[name]
		return v, ok
	}
}

func (at accessToken) HasScope(scope string) bool {
	if s0, ok := at.Claim("scope"); ok {
		s1, _ := s0.(string)
		return slices.Contains(strings.Split(s1, " "), scope)
	}
	return false
}

func (at accessToken) HasPermission(permission string) bool {
	if p0, ok := at.Claim("permissions"); ok {
		switch p1 := p0.(type) {
		case []any: // Auth0 RBAC
			for _, p2 := range p1 {
				if p3, _ := p2.(string); p3 == permission {
					return true
				}
			}
		case string:
			return slices.Contains(strings.Split(p1, " "), permission)
		}
	}
	return false
}

func (oa OAuth) parseAccessToken(ctx context.Context, token string) (*accessToken, error) {
	jwtParseOpts := []jwt.ParseOption{
		jwt.WithContext(ctx),
		jwt.WithKeySet(oa.jwks),
		jwt.WithIssuer(oa.issuer),
	}
	if oa.audience != "" {
		jwtParseOpts = append(jwtParseOpts, jwt.WithAudience(oa.audience))
	}
	result, err := jwt.ParseString(token, jwtParseOpts...)
	if err != nil {
		switch {
		case oa.IsTokenExpiredError(err):
			return nil, model.WrappedError(err, "expired access token")
		case oa.IsTokenValidationError(err):
			return nil, model.WrappedError(err, "invalid access token")
		}
		return nil, err
	}
	return &accessToken{
		Subject:    result.Subject(),
		Expiration: result.Expiration(),
		Others:     result.PrivateClaims(),
	}, nil
}
