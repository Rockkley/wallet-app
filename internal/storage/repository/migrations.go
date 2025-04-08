package repository

import (
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"log"

	"github.com/pressly/goose/v3"
)

func RunMigrations(conn *pgx.Conn) error {
	connConfig := conn.Config()

	db := stdlib.OpenDB(*connConfig)
	defer db.Close()

	if err := goose.SetDialect("postgres"); err != nil {

		return fmt.Errorf("error setting dialect: %w", err)
	}

	if err := goose.Up(db, "migrations"); err != nil {

		return fmt.Errorf("error while migration up: %w", err)
	}

	log.Println("migration done successful")

	return nil
}
