package middleware

import (
	"order-service/src/common/dto/response"
	"order-service/src/common/errors"
	"order-service/src/common/log"

	"github.com/gofiber/fiber/v3"
)

func (m *Middleware) ErrorHandler(c fiber.Ctx, err error) error {

	if e, ok := err.(*errors.Response); ok {
		return c.Status(e.HttpCode).JSON(response.CommonResponse{Message: e.Message})
	}

	if e, ok := err.(*fiber.Error); ok {
		return c.Status(e.Code).JSON(response.CommonResponse{Message: e.Message})
	}

	log.Logger.Error(err.Error())

	return c.Status(500).JSON(response.CommonResponse{Message: "internal server error"})
}
