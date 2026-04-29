package auth

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"github.com/gofiber/fiber/v3/middleware/session"
)

func LogOut(ctx fiber.Ctx) error {
	s := session.FromContext(ctx)
	if err := s.Reset(); err != nil {
		return ctx.Status(500).JSON(fiber.Map{"error": "Session error"})
	}
	log.Info("User logged out")
	return ctx.SendString("loged out")
}
