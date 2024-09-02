package storage

import (
	"database/sql"
	"embed"
	"log/slog"
	"os"

	_ "github.com/lib/pq"
	goose "github.com/pressly/goose/v3"
)

type DB struct {
	driver *sql.DB
}

//go:embed migrations/*.sql
var embedMigrations embed.FS

func upMigrations(driver *sql.DB) {
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		slog.Error("Error when set dialect", "err", err)
		os.Exit(1)
	}

	if err := goose.Up(driver, "migrations"); err != nil {
		slog.Error("Error when up migrations", "err", err)
		os.Exit(1)
	}
}

func MakeDB(dsn string) DB {
	driver, err := sql.Open("postgres", dsn)
	if err != nil {
		slog.Error("Error when open database")
		os.Exit(1)
	}

	upMigrations(driver)

	return DB{
		driver: driver,
	}
}
