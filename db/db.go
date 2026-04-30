package db

import (
	"context"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v3/log"
	"github.com/jackc/pgx/v5/pgxpool"
)

var DbConnection *pgxpool.Pool

func Connect() {
	var err error
	connstring := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB"))
	DbConnection, err = pgxpool.New(context.Background(), connstring)
	if err != nil {
		log.Fatalf("Error connecting db: %v", err)
	}

	// Initiating tables
	if err := createCategoriesTable(); err != nil {
		log.Fatalf("Error creating categories table: %v", err)
	}
	if err := createProductsTable(); err != nil {
		log.Fatalf("Error creating products table: %v", err)
	}
	if err := createUserTable(); err != nil {
		log.Fatalf("Error creating users table: %v", err)
	}

}
