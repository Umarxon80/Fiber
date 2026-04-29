package auth

import (
	"github.com/Umarxon80/Fiber.git/db"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"github.com/gofiber/fiber/v3/middleware/session"
)

type login struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"min=8"`
}

func LogIn(ctx fiber.Ctx) error {
	var lData login
	if err := ctx.Bind().Body(&lData); err != nil {
		log.Error("Wrong input")
		return ctx.JSON(fiber.Map{
			"error": "password or email is incorrect",
		})
	}
	user, err := db.GetUserByEmail(ctx, lData.Email)
	if err != nil {
		log.Error("Wrong input")
		return ctx.JSON(fiber.Map{
			"error": "password or email is incorrect",
		})
	}
	check, err := checkPassword(user.Password, lData.Password)
	if err != nil || !check {
		log.Error("Wrong input")
		return ctx.JSON(fiber.Map{
			"error": "password or email is incorrect",
		})
	}
	s := session.FromContext(ctx)
	s.Set("authenticated", true)
	s.Set("is_admin",user.Is_admin)
	log.Info("User logged in, id: ", user.Id)
	return ctx.SendString("loged in")
}
