package validation

import (
	"github.com/Umarxon80/Fiber.git/db"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
)

func ValidateUserInput(ctx fiber.Ctx) error {
	var user db.User

	if err := ctx.Bind().JSON(&user); err != nil {
		log.Error("Wrong input")
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "wrong input",
			"err":err.Error(),
		})
	}

	val := validator.New(validator.WithRequiredStructEnabled())
	if err := val.Struct(user); err != nil {
		log.Error("Validation failed")
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Next()
}

func ValidateProductInput(ctx fiber.Ctx) error {
	var product db.Product

	if err := ctx.Bind().JSON(&product); err != nil {
		log.Error("Wrong input")
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "wrong input",
						"err":err.Error(),

		})
	}

	val := validator.New(validator.WithRequiredStructEnabled())
	if err := val.Struct(product); err != nil {
		log.Error("Validation failed")
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Next()
}

func ValidateCategoryInput(ctx fiber.Ctx) error {
	var product db.Category

	if err := ctx.Bind().JSON(&product); err != nil {
		log.Error("Wrong input")
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "wrong input",
						"err":err.Error(),

		})
	}

	val := validator.New(validator.WithRequiredStructEnabled())
	if err := val.Struct(product); err != nil {
		log.Error("Validation failed")
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Next()
}
