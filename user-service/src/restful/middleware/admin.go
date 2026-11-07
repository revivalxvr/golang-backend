package middleware

import (
	"user-service/src/common/errors"

	"github.com/gofiber/fiber/v3"
)

func (m *Middleware) Admin(c fiber.Ctx) error {
	role, ok := c.Locals("role").(string)

	// Jika role tidak ada (nil) atau bukan "ADMIN", tolak dengan 403 Forbidden
	if !ok || role != "ADMIN" {
		return &errors.Response{HttpCode: 403, Message: "forbidden"}
	}
	return c.Next()
}
