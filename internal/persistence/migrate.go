package persistence

import (
	"embed"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jmoiron/sqlx"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// RunMigrations applies all pending migrations to the database.
// The driver parameter must be "sqlite" or "pgx".
func RunMigrations(db *sqlx.DB, driver string) error {
	sourceDriver, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		return fmt.Errorf("creating migration source: %w", err)
	}

	var dbDriver database.Driver
	switch driver {
	case "sqlite":
		dbDriver, err = sqlite.WithInstance(db.DB, &sqlite.Config{})
	case "pgx":
		dbDriver, err = pgx.WithInstance(db.DB, &pgx.Config{})
	default:
		return fmt.Errorf("unsupported migration driver: %s", driver)
	}
	if err != nil {
		return fmt.Errorf("creating migration db driver: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", sourceDriver, driver, dbDriver)
	if err != nil {
		return fmt.Errorf("creating migrator: %w", err)
	}

	// Fix dirty state from a previous failed migration.
	version, dirty, _ := m.Version()
	if dirty {
		if err := m.Force(int(version)); err != nil { //nolint:gosec // version is a small migration number, no overflow risk
			return fmt.Errorf("fixing dirty migration: %w", err)
		}
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("running migrations: %w", err)
	}

	return nil
}
