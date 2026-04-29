package auth

import (
	"github.com/Umarxon80/Fiber.git/db"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"github.com/gofiber/fiber/v3/middleware/session"
	"golang.org/x/crypto/bcrypt"
)

func RequreAuth(ctx fiber.Ctx) error {
	sess := session.FromContext(ctx)
	if sess == nil {
		log.Error("Session expired")
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"error": "Session expired",
		})
	}
	if sess.Get("authenticated") != true {
		log.Error("Session expired")
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"error": "Session expired",
		})
	}
	return ctx.Next()
}

func HashPassword(ctx fiber.Ctx) error {
	var user db.User
	if err := ctx.Bind().Body(&user); err != nil {
		log.Error("Wrong input")
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"error": err,
		})
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 7)
	if err != nil {
		log.Error("Error hashing password: ", user.Password)
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"error": err,
		})
	}
	user.Password = string(hashedPassword)
	ctx.Locals("user", user)
	return ctx.Next()
}

func checkPassword(hash, password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return false, err
	} else {
		return true, nil
	}
}
