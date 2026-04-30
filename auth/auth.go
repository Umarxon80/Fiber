package auth

import (
	"fmt"
	"strings"

	"github.com/Umarxon80/Fiber.git/db"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"golang.org/x/crypto/bcrypt"
)

func RequreAuth(ctx fiber.Ctx) error {
	authHeader := ctx.Get("Authorization")
	if authHeader == "" {
		log.Error("No token provided")
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "no token provided 1"})
	}
	devidedAuth := strings.SplitN(authHeader, " ", 2)
	if len(devidedAuth) != 2 || devidedAuth[0] != "Bearer" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token format"})
	}
	claims, err := checkToken(devidedAuth[1])
	fmt.Print(devidedAuth[1])
	if err != nil {
		log.Error("No token provided")
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "no token provided 2", "err": err.Error()})
	}
	ctx.Locals("userid", claims["id"])
	ctx.Locals("role", claims["role"])
	return ctx.Next()
}

func HashPassword(ctx fiber.Ctx) error {
	var user db.User
	if err := ctx.Bind().Body(&user); err != nil {
		log.Error("Wrong input")
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"error": err.Error(),
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
