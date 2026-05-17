package migration

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

//go:embed "migrations/*.sql"
var embedMigrations embed.FS

func Run(dsn string, debug bool) error {
	if dsn == "" {
		return errors.New("dsn required")
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database when running migration: %w", err)
	}
	defer func() {
		_ = db.Close()
	}()

	goose.SetBaseFS(embedMigrations)

	err = goose.SetDialect("postgres")
	if err != nil {
		return fmt.Errorf("failed to set postgres dialect: %w", err)
	}

	if debug {
		// Should not be set to True in release
		err = goose.Reset(db, "migrations")
		if err != nil {
			return fmt.Errorf("failed to reset migration table: %w", err)
		}
	}

	if err := goose.Up(db, "migrations"); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}
