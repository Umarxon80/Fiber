package db

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
)

type Product struct {
	Id         uint      `json:"id"`
	Name       string    `json:"name"`
	Desc       string    `json:"description"`
	Price      int       `json:"price"  validate:"required"`
	ExpDate    time.Time `json:"exp_date"`
	Category   Category  `json:"category" validate:"omitempty"`
	CategoryId uint      `json:"category_id"  validate:"required"`
}

func createProductsTable() error {
	_, err := DbConnection.Exec(context.Background(), `
        CREATE TABLE IF NOT EXISTS products (
            id       SERIAL PRIMARY KEY,
            name     TEXT           NOT NULL,
            description    TEXT,
            price    NUMERIC(10, 2) NOT NULL,
            exp_date DATE,
			category_id INT REFERENCES categories(id) ON DELETE SET NULL
        )
    `)
	return err
}

// // Products CRUD

func CreateProduct(ctx fiber.Ctx) error {
	var newProd Product
	var id uint
	if err := ctx.Bind().Body(&newProd); err != nil {
		log.Error("Wrong input", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	err := DbConnection.QueryRow(context.Background(), `
		INSERT INTO products (name, description, price, exp_date,category_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`, newProd.Name, newProd.Desc, newProd.Price, newProd.ExpDate, newProd.CategoryId).Scan(&id)
	if err != nil {
		log.Error("Database error", err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	log.Info("New product crated id: ", id)
	return ctx.Status(fiber.StatusCreated).JSON(id)
}
func GetProducts(ctx fiber.Ctx) error {
	var products []Product

	rows, err := DbConnection.Query(context.Background(), `SELECT p.id,p.name,p.description,p.price,p.exp_date,c.id,c.name,p.category_id FROM products p LEFT JOIN categories c on p.category_id=c.id`)
	if err != nil {
		log.Errorf("Error getting all products %v ", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	for rows.Next() {
		var buff Product
		err := rows.Scan(&buff.Id, &buff.Name, &buff.Desc, &buff.Price, &buff.ExpDate, &buff.Category.Id, &buff.Category.Name, &buff.CategoryId)
		if err != nil {
			log.Errorf("Error getting all products %v ", err)
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		products = append(products, buff)
	}
	log.Info("Returning all products")
	return ctx.JSON(products)
}
func GetOneProduct(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	var pr Product
	if err := DbConnection.QueryRow(context.Background(), `SELECT p.id,p.name,p.price,p.description,p.exp_date,c.id,c.name,p.category_id FROM products p LEFT JOIN categories c on p.category_id=c.id where p.id=$1`, id).Scan(&pr.Id, &pr.Name, &pr.Price, &pr.Desc, &pr.ExpDate, &pr.Category.Id, &pr.Category.Name, &pr.CategoryId); err != nil {
		log.Errorf("Wrong input; %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
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
        UPDATE products 
		set name=$1, description=$2, price=$3, exp_date=$4, category_id=$5
        WHERE id=$6
        RETURNING id, name, description, price, exp_date
    `, product.Name, product.Desc, product.Price, product.ExpDate, product.CategoryId, id)
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
        DELETE FROM products
        WHERE id=$1
    `, id)
	if err != nil {
		log.Error("Error deleting product, ", err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	if ct.RowsAffected() < 1 {
		log.Error("Product with this id does not exists")
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Not found",
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"Product deleted, id": id,
	})
}
