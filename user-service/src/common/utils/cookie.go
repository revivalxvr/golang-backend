package utils

import (
	"time"

	"github.com/gofiber/fiber/v3"
)

func ClearCookie(c fiber.Ctx, name, path string) {
	c.Cookie(&fiber.Cookie{
		Name:     name,
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
		Path:     path,
	})
}
