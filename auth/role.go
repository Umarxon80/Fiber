package auth

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/session"
)

func RoleChecker( allowedRoles []string) fiber.Handler {
	return func(ctx fiber.Ctx) error {
		s:=session.FromContext(ctx)
        userRole := s.Get("role") 
        for _, r := range allowedRoles {
            if userRole == r {
                return ctx.Next()
            }
        }
        return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
            "error": "Forbidden",
        })
    }
}