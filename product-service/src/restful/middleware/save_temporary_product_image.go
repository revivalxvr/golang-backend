package middleware

import (
	"product-service/src/common/errors"
	"product-service/src/common/utils"

	"github.com/gofiber/fiber/v3"
)

func (m *Middleware) SaveTemporaryProductImage(c fiber.Ctx) error {
	file, err := c.FormFile("product_image")
	if err != nil {
		return c.Next()
	}
	var maxSize int64 = 1 * 1000 * 1000 //Mb
	if file.Size > maxSize {
		return &errors.Response{Message: "file size too large", HttpCode: 400}
	}
	fileName := utils.CreateUniqueFileName(file.Filename)
	path := "./tmp/" + fileName

	if err := utils.CheckExisDir("./tmp"); err != nil {
		return err
	}
	if err := c.SaveFile(file, path); err != nil {
		return err
	}
	c.Locals("fileName", fileName)
	return c.Next()
}
