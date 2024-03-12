package datastore

import (
	"context"

	"github.com/mahmudindes/orenocomic-bagicore/internal/datastore/database"
	"github.com/mahmudindes/orenocomic-bagicore/internal/datastore/redis"
)

type (
	datastore struct {
		Database *database.Database
		Redis    *redis.Redis
	}

	Config struct {
		Database database.Config `conf:"database"`
		Redis    redis.Config    `conf:"redis"`
	}
)

func New(ctx context.Context, cfg Config) (datastore, error) {
	ds := datastore{}

	db, err := database.New(ctx, cfg.Database)
	if err != nil {
		return ds, err
	}
	ds.Database = db

	rd, err := redis.New(ctx, cfg.Redis)
	if err != nil {
		return ds, err
	}
	ds.Redis = rd

	return ds, nil
}

func (ds datastore) Stop() error {
	if ds.Database != nil {
		if err := ds.Database.Close(); err != nil {
			return err
		}
	}
	return nil
}
