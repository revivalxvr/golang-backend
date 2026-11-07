package middleware

import (
	"os"

	"product-service/src/common/errors"

	"github.com/gofiber/fiber/v3"
	"github.com/h2non/filetype"
)

func (m *Middleware) ValidateProductImage(c fiber.Ctx) error {
	fileName, ok := c.Locals("fileName").(string)
	if !ok || fileName == "" {
		return c.Next()
	}
	file, err := os.Open("./tmp/" + fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	head := make([]byte, 261)
	file.Read(head)

	if !filetype.IsImage(head) {
		return &errors.Response{HttpCode: 400, Message: "Invalid Image"}
	}
	return c.Next()

}
