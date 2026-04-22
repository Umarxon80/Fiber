package db

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type Product struct {
	Id      int    `json:"id" gorm:"primaryKey"`
	Name    string `json:"name"`
	Desc    string `json:"desc"`
	Price   int    `json:"price"`
	ExpDate string `json:"exp_date"`
}

var db *gorm.DB

func ConnectDb() {
	var err error
	db, err = gorm.Open(postgres.Open("host=localhost user=postgres dbname=fiber password=1234 sslmode=disable"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		PrepareStmt: true,
	})
	if err != nil {
		log.Fatal("Failed to connect to DB")
	}
	log.Debug("Connected to DB")
}
func Migrate() {
	db.AutoMigrate(&Product{})
}

func Add(ctx fiber.Ctx) error {
	var product Product
	if err := ctx.Bind().Body(&product); err != nil {
		log.Error("Wrong input")
		return err
	}
	if product.Name == "" && product.Price < 0 {
		log.Error("Wrong input")
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"error": "invalid input",
		})
	}

	db.Create(&product)
	log.Info("New product crated id: ", product.Id)
	return ctx.JSON(product)
}
func Get(ctx fiber.Ctx) error {
	var products []Product
	db.Find(&products)
	log.Info("Returning all Products")
	return ctx.JSON(products)
}
func GetOne(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	var product Product
	if err := db.First(&product, id).Error; err != nil {
		log.Error("Wrong input")
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Product not found",
		})
	}
	log.Info("Returning one product, id: ", product.Id)
	return ctx.JSON(product)
}
func Patch(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	var body Product
	if err := ctx.Bind().Body(&body); err != nil {
		log.Error("Wrong input")
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Empty input",
		})
	}
	db.Model(&Product{}).Where("id=?", id).Updates(Product{
		Name:    body.Name,
		Price:   body.Price,
		Desc:    body.Desc,
		ExpDate: body.ExpDate,
	})
	log.Info("Product updated, id: ", id)
	return ctx.JSON(fiber.Map{
		"Product updated, id: ": id,
	})
}
func Delete(ctx fiber.Ctx) error {
	result := db.Delete(&Product{}, ctx.Params("id"))

	if result.RowsAffected == 0 {
		log.Error("Wrong input")
		return ctx.Status(404).JSON(fiber.Map{
			"error": "not found",
		})
	}
	log.Info("Deleted product id: ", ctx.Params("id"))
	return ctx.SendString("Deleted successfully")
}
