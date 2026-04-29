package db

import (
	"context"

	"github.com/gofiber/fiber/v3/log"
	"github.com/jackc/pgx/v5/pgxpool"
)

var DbConnection *pgxpool.Pool

func Connect() {
	var err error
	DbConnection, err = pgxpool.New(context.Background(), "postgres://postgres:1234@localhost:5432/fiber?sslmode=disable")
	if err != nil {
		log.Fatalf("Error connecting db: %v", err)
	}

	// Initiating tables
	if err := createProductTable(); err != nil {
		log.Fatalf("Error creating products table: %v", err)
	}
	if err := createUserTable(); err != nil {
		log.Fatalf("Error creating users table: %v", err)
	}

}
