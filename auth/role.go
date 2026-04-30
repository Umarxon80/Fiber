package auth

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
)

func RoleChecker(allowedRoles []string) fiber.Handler {
	return func(ctx fiber.Ctx) error {
		userRole := ctx.Locals("role")
		for _, r := range allowedRoles {
			if userRole == r {
				return ctx.Next()
			}
		}
		log.Error("Not allowev method id: ", ctx.Locals("id"))
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Forbidden",
		})
	}
}
