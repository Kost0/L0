package repository

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

// RunMigrations launches migrations to database
// Accepts:
//   - db: database
//   - dbName: name of database
//
// Returns:
//   - error if something wrong
func RunMigrations(db *sql.DB, dbName string) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("could not create migration driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		dbName,
		driver,
	)
	if err != nil {
		return fmt.Errorf("could not create migration: %w", err)
	}

	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("could not apply migration: %w", err)
	}

	return nil
}
