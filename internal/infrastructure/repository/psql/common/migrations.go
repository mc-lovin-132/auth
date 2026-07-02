package repo

import (
	"database/sql"
	"log"

	"github.com/pressly/goose/v3"
)

func RunMigrations(db *sql.DB) error {
	// Указываем диалект БД
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	// Накатываем миграции из папки "migrations"[citation:9]
	if err := goose.Up(db, "./migrations"); err != nil {
		return err
	}

	log.Println("Migrations applied successfully")
	return nil

}
