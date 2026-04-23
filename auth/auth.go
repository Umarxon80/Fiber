package auth

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"github.com/gofiber/fiber/v3/middleware/session"
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
