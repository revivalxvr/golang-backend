package middleware

import (
	"fmt"
	"runtime/debug"
	"user-service/src/common/dto/response"
	"user-service/src/common/errors"
	"user-service/src/common/log"

	"github.com/gofiber/fiber/v3"
)

func (m *Middleware) ErrorHandler(c fiber.Ctx, err error) error {
	// Menangkap stacktrace lengkap jika terjadi panic di controller/repository
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("=== PANIC TERDETEKSI ===")
			fmt.Printf("Panic: %v\n", r)
			debug.PrintStack()
			fmt.Println("========================")
		}
	}()

	if e, ok := err.(*errors.Response); ok {
		return c.Status(e.HttpCode).JSON(response.CommonResponse{Message: e.Message})
	}

	if e, ok := err.(*fiber.Error); ok {
		return c.Status(e.Code).JSON(response.CommonResponse{Message: e.Message})
	}

	log.Logger.Error(err.Error())
	return c.Status(500).JSON(response.CommonResponse{Message: "internal server error"})
}
