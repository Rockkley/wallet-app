package repository

import (
	"context"
	"github.com/jackc/pgx/v5"
)

type Database struct {
	Conn *pgx.Conn
}

func NewDatabase(connStr string) (*Database, error) {
	conn, err := pgx.Connect(context.Background(), connStr)

	if err != nil {

		return nil, err
	}

	return &Database{
		Conn: conn,
	}, nil
}
