package middleware

import (
	"order-service/src/common/errors"

	"github.com/gofiber/fiber/v3"
)

func (m *Middleware) UserRole(c fiber.Ctx) error {
	role, ok := c.Locals("role").(string)

	// Jika role tidak ada (nil) atau bukan "USER", tolak dengan 403 Forbidden
	if !ok || role != "USER" {
		return &errors.Response{HttpCode: 403, Message: "forbidden"}
	}
	return c.Next()
}
