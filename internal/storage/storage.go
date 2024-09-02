package storage

import (
	"database/sql"
	"log/slog"
	"os"

	_ "github.com/lib/pq"
)

type DB struct {
	driver *sql.DB
}

func MakeDB(dsn string) DB {
	slog.Info("dsn for db", "dsn", dsn)
	driver, err := sql.Open("postgres", dsn)
	if err != nil {
		slog.Error("Error when open database")
		os.Exit(1)
	}

	return DB{
		driver: driver,
	}

}
