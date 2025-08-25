package migration

import (
	"embed"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func MigrateTables(db *sqlx.DB) error {
	// sqlx.DBから*sql.DBを取得
	sqlDB := db.DB
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("mysql"); err != nil {
		return fmt.Errorf("set dialect: %w", err)
	}

	if err := goose.Up(sqlDB, "migrations"); err != nil {
		return fmt.Errorf("up migration: %w", err)
	}

	return nil
}
