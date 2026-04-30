package db

import (
	"context"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
)

type Category struct {
	Id   uint   `json:"id"`
	Name string `json:"name" validate:"required"`
}

func createCategoriesTable() error {
	_, err := DbConnection.Exec(context.Background(), `
        CREATE TABLE IF NOT EXISTS categories (
            id       SERIAL PRIMARY KEY,
            name     TEXT           NOT NULL
        )
    `)
	return err
}

// // Categories CRUD

func CreateCategory(ctx fiber.Ctx) error {
	var newCategory Category
	if err := ctx.Bind().Body(&newCategory); err != nil {
		log.Error("Wrong input", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	err := DbConnection.QueryRow(context.Background(), `
		INSERT INTO categories (name)
		VALUES ($1)
		RETURNING id
	`, newCategory.Name).Scan(&newCategory.Id)
	if err != nil {
		log.Error("Database error", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	log.Info("New Category crated id: ", newCategory.Id)
	return ctx.Status(fiber.StatusCreated).JSON(newCategory)
}
func GetCategories(ctx fiber.Ctx) error {
	var Categories []Category

	rows, err := DbConnection.Query(context.Background(), `SELECT * FROM categories`)
	if err != nil {
		log.Errorf("Error getting all categories %v ", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err,
		})
	}
	for rows.Next() {
		var buff Category
		err := rows.Scan(&buff.Id, &buff.Name)
		if err != nil {
			log.Errorf("Error getting all Categories %v ", err)
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err,
			})
		}
		Categories = append(Categories, buff)
	}
	log.Info("Returning all Categorys")
	return ctx.JSON(Categories)
}
func GetOneCategory(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	var category Category
	if err := DbConnection.QueryRow(context.Background(), `SELECT * FROM categories where id=$1`, id).Scan(&category.Id, &category.Name); err != nil {
		log.Errorf("Wrong input; %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	log.Info("Returning one Category, id: ", category.Id)
	return ctx.JSON(category)
}
func PatchCategory(ctx fiber.Ctx) error {
	var Category Category
	id := ctx.Params("id")
	if err := ctx.Bind().Body(&Category); err != nil {
		log.Error("Wrong input")
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Wrong input",
		})
	}
	efRows, err := DbConnection.Exec(context.Background(), `
        UPDATE categories 
		set name=$1
        WHERE id=$2
        RETURNING id, name
    `, Category.Name, id)
	if efRows.RowsAffected() < 1 {
		log.Warn("Not found, id: ", id)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Not found",
		})
	}
	if err != nil {
		log.Fatal(err)
	}
	log.Info("Category updated, id: ", id)
	return ctx.JSON(fiber.Map{
		"Category updated, id: ": id,
	})
}
func DeleteCategory(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	ct, err := DbConnection.Exec(context.Background(), `
        DELETE FROM categories 
        WHERE id=$1
    `, id)
	if err != nil {
		log.Error("Error deleting Category, ", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err,
		})
	}
	if ct.RowsAffected() < 1 {
		log.Error("Category with this id does not exists")
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Not found",
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"Category deleted, id": id,
	})
}
