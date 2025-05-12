package migrations

import (
	"database/sql"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const (
	CurrentDatabaseMigrationVersion = 1
	DefaultSourceUrl                = "file://internal/database/migrations"
)

func RunDatabaseMigrations(db *sql.DB, source string, version uint) error {
	driver, err := pgx.WithInstance(db, &pgx.Config{})
	if err != nil {
		log.Fatal("Failed to create a new migration driver")
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(source, "gopgtest", driver)
	if err != nil {
		log.Fatalf("failed to create a new migration instance: %v", err)
		return err
	}

	err = m.Migrate(version)
	if err != nil && err != migrate.ErrNoChange {
		log.Fatal("failed to run migrations")
		m.Down()
		return err
	}

	return nil
}
