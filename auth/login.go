package auth

import (
	"github.com/Umarxon80/Fiber.git/db"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
)

type login struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"min=8"`
}

func LogIn(ctx fiber.Ctx) error {
	var lData login
	if err := ctx.Bind().Body(&lData); err != nil {
		log.Error("Wrong input")
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "password or email is incorrect",
		})
	}
	user, err := db.GetUserByEmail(ctx, lData.Email)
	if err != nil {
		log.Error("Wrong input")
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "password or email is incorrect",
		})
	}
	check, err := checkPassword(user.Password, lData.Password)
	if err != nil || !check {
		log.Error("Wrong input")
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "password or email is incorrect",
		})
	}
	token, err := generateToken(user.Id, user.Role)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not generate token"})
	}

	return ctx.JSON(fiber.Map{"token": token})
}
