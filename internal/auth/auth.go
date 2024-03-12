package auth

import (
	"context"

	"github.com/mahmudindes/orenocomic-bagicore/internal/auth/oauth"
	"github.com/mahmudindes/orenocomic-bagicore/internal/logger"
)

type (
	auth struct {
		OAuth *oauth.OAuth
	}

	Config struct {
		OAuth oauth.Config `conf:"oauth"`
	}

	redis oauth.Redis
)

func New(ctx context.Context, rdb redis, cfg Config, log logger.Logger) (*auth, error) {
	oa, err := oauth.New(ctx, rdb, cfg.OAuth, log.WithName("OAuth"))
	if err != nil {
		return nil, err
	}
	return &auth{OAuth: oa}, nil
}
