package migrator

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

func RunMigrations(db *sql.DB) error {
	migrationsDir := "./db/migrations"

	if err := goose.Up(db, migrationsDir); err != nil {
		return err
	}

	log.Println("Migrations ran successfully.")
	return nil
}
