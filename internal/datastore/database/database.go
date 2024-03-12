package database

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"

	bagicore "github.com/mahmudindes/orenocomic-bagicore"
	"github.com/mahmudindes/orenocomic-bagicore/embedded"
)

type (
	Database struct {
		client *pgxpool.Pool
	}

	Config struct {
		URL      string `conf:"url"`
		Provider string `conf:"provider"`
	}
)

func New(ctx context.Context, cfg Config) (*Database, error) {
	config, err := pgxpool.ParseConfig(cfg.URL)
	if err != nil {
		return nil, err
	}

	client, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	db := &Database{client: client}

	if err := db.Migrate(ctx, strings.ToLower(cfg.Provider)); err != nil {
		return nil, err
	}

	return db, nil
}

func (db Database) Close() error {
	if db.client != nil {
		db.client.Close()
	}

	return nil
}

func (db Database) Migrate(ctx context.Context, provider string) error {
	if err := goose.SetDialect("pgx"); err != nil {
		return err
	}
	switch provider {
	case "crdb", "cockroachdb":
		goose.SetBaseFS(embedded.CRDBMigrations)
	case "pg", "postgres", "postgresql":
		goose.SetBaseFS(embedded.PGMigrations)
	default:
		return errors.New("database provider " + provider + " not supported")
	}
	goose.SetTableName(bagicore.ID + "_version")
	goose.SetLogger(&migrationLogger{})

	if err := goose.UpContext(ctx, stdlib.OpenDBFromPool(db.client), "migrations"); err != nil {
		return err
	}

	goose.SetBaseFS(nil)

	return nil
}

type migrationLogger struct{}

func (l *migrationLogger) Fatalf(format string, v ...interface{}) { panic(fmt.Sprintf(format, v...)) }
func (l *migrationLogger) Printf(format string, v ...interface{}) {}
