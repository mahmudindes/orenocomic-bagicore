package oauth

import (
	"context"
	"errors"
)

type (
	ctxAccessToken    struct{}
	ctxAccessTokenRaw struct{}
)

func (oa OAuth) ContextAccessToken(ctx context.Context, token string) context.Context {
	ctx = context.WithValue(ctx, ctxAccessToken{}, new(accessToken))
	ctx = context.WithValue(ctx, ctxAccessTokenRaw{}, token)
	return ctx
}

func (oa OAuth) parseTokenContext(ctx context.Context) (*accessToken, error) {
	aTokenRaw, aTokenOK := ctx.Value(ctxAccessTokenRaw{}).(string)
	if aTokenOK {
		cache, err := oa.cTokenCache.GetToken(ctx, aTokenRaw)
		switch {
		case err != nil:
			oa.logger.ErrMessage(err, "Parse get context token cache failed.")
		case cache != nil:
			return cache, nil
		}
		aToken, err := oa.parseAccessToken(ctx, aTokenRaw)
		if err != nil {
			return nil, err
		}
		go func() {
			if err := oa.cTokenCache.SetToken(context.Background(), aTokenRaw, aToken); err != nil {
				oa.logger.ErrMessage(err, "Parse set context token cache failed.")
			}
		}()
		return aToken, nil
	}
	return nil, nil
}

func (oa OAuth) getAccessTokenContext(ctx context.Context) (*accessToken, error) {
	if token, ok := ctx.Value(ctxAccessToken{}).(*accessToken); ok {
		return token, nil
	}
	return nil, errors.New("access token context not exists")
}

func (oa OAuth) setAccessTokenContext(ctx context.Context, token accessToken) error {
	tokenCtx, err := oa.getAccessTokenContext(ctx)
	if err == nil {
		*tokenCtx = token
	}
	return err
}

func (oa OAuth) ProcessTokenContext(ctx context.Context) (bool, error) {
	token, err := oa.parseTokenContext(ctx)
	switch {
	case err != nil:
		return false, err
	case token != nil:
		return true, oa.setAccessTokenContext(ctx, *token)
	}
	return false, nil
}

func (oa OAuth) HasPermissionContext(ctx context.Context, permission string) bool {
	token, err := oa.getAccessTokenContext(ctx)
	if err != nil {
		return false
	}
	return token.HasPermission(permission)
}
