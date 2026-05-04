// Package persistence provides database access for the subject-data service.
package persistence

import (
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib" // PostgreSQL driver
	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite" // Pure-Go SQLite driver (CGO_ENABLED=0)
)

// NewSQLiteDB opens a SQLite database at the given path and returns a sqlx.DB.
func NewSQLiteDB(path string) (*sqlx.DB, error) {
	db, err := sqlx.Open("sqlite", path+"?_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)&_pragma=foreign_keys(1)")
	if err != nil {
		return nil, fmt.Errorf("opening sqlite db: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("pinging sqlite db: %w", err)
	}
	return db, nil
}

// NewPostgresDB opens a PostgreSQL database with the given DSN and returns a sqlx.DB.
func NewPostgresDB(dsn string) (*sqlx.DB, error) {
	db, err := sqlx.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("opening postgres db: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("pinging postgres db: %w", err)
	}
	return db, nil
}
