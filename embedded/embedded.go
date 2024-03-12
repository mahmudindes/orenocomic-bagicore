package embedded

import "embed"

var (
	//go:embed default-config.yml
	DefaultConfig []byte

	//go:embed migrations/*-crdb.sql
	CRDBMigrations embed.FS
	//go:embed migrations/*-pg.sql
	PGMigrations embed.FS
)
