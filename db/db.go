package db

import (
	"context"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Product struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Desc    string `json:"desc"`
	Price   int    `json:"price"`
	ExpDate string `json:"exp_date"`
}

var DbConnection *pgxpool.Pool

func Connect() {
	var err error
	DbConnection, err = pgxpool.New(context.Background(), "postgres://postgres:1234@localhost:5432/fiber?sslmode=disable")
	if err != nil {
		log.Fatalf("Error connecting db: %v", err)
	}
	if err := createProductTable(context.Background()); err != nil {
		log.Fatalf("Error creating table: %v", err)
	}
	defer DbConnection.Close()
}
func createProductTable(ctx context.Context) error {
	_, err := DbConnection.Exec(ctx, `
        CREATE TABLE IF NOT EXISTS products (
            id       SERIAL PRIMARY KEY,
            name     TEXT           NOT NULL,
            "description"     TEXT,
            price    NUMERIC(10, 2) NOT NULL,
            exp_date DATE
        )
    `)
	return err
}

// // PRODUCT CRUD
func CreateProduct(ctx fiber.Ctx) error {
	var newProd Product
	if err := ctx.Bind().Body(&newProd); err != nil {
		log.Error("Wrong input")
		return err
	}
	err := DbConnection.QueryRow(context.Background(), `
		INSERT INTO product (name, description, price, exp_date)
		VALUES ($1, $2, $3, $4)
		RETURNING id, name, description, price, exp_date
	`, newProd.Name, newProd.Desc, newProd.Price, newProd.ExpDate).Scan(&newProd.Id, &newProd.Name, &newProd.Desc, &newProd.Price, &newProd.ExpDate)
	if err != nil {
		log.Fatal(err)
	}
	log.Info("New product crated id: ", newProd.Id)
	return ctx.Status(fiber.StatusCreated).JSON(newProd)
}
func GetProducts(ctx fiber.Ctx) error {
	var products []Product

	rows, err := DbConnection.Query(context.Background(), `SELECT * FROM product`)
	if err != nil {
		log.Errorf("Error getting all products %v ", err)
		return err
	}
	for rows.Next() {
		var buff Product
		err := rows.Scan(&buff.Id, &buff.Name, &buff.Desc, &buff.Price, &buff.ExpDate)
		if err != nil {
			log.Errorf("Error getting all products %v ", err)
			return err
		}
		products = append(products, buff)
	}
	log.Info("Returning all products")
	return ctx.JSON(products)
}
func GetOneProduct(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	var pr Product
	if err := DbConnection.QueryRow(context.Background(), `SELECT * FROM product where id=$1`, id).Scan(&pr.Id, &pr.Name, &pr.Desc, &pr.Price, &pr.ExpDate); err != nil {
		log.Errorf("Wrong input; %v", err)
	}
	log.Info("Returning one product, id: ", pr.Id)
	return ctx.JSON(pr)
}
func PatchProduct(ctx fiber.Ctx) error {
	var product Product
	id := ctx.Params("id")
	if err := ctx.Bind().Body(&product); err != nil {
		log.Error("Wrong input")
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Wrong input",
		})
	}
	efRows, err := DbConnection.Exec(context.Background(), `
        UPDATE product 
		set name=$1, description=$2, price=$3, exp_date=$4
        WHERE id=$5
        RETURNING id, name, description, price, exp_date
    `, product.Name, product.Desc, product.Price, product.ExpDate, id)
	if efRows.RowsAffected() < 1 {
		log.Warn("Not found, id: ", id)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Not found",
		})
	}
	if err != nil {
		log.Fatal(err)
	}
	log.Info("Product updated, id: ", id)
	return ctx.JSON(fiber.Map{
		"Product updated, id: ": id,
	})
}
func DeleteProduct(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	ct, err := DbConnection.Exec(context.Background(), `
        DELETE FROM product 
        WHERE id=$1
    `, id)
	if err != nil {
		log.Error("Error deleting product, ", err)
		return err
	}
	if ct.RowsAffected() < 1 {
		log.Fatal("Product with this id does not exists")
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Not found",
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"Product deleted, id": id,
	})
}
