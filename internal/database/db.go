package database

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
)

func InitDB(connString string) *pgx.Conn {
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, connString)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	return conn
}