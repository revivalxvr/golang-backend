package middleware

import (
	"product-service/src/common/dto/request"
	"product-service/src/common/errors"
	"time"

	"github.com/google/uuid"

	"github.com/gofiber/fiber/v3"
)

func (m *Middleware) Log(c fiber.Ctx) error {
	now := time.Now()
	requestID := uuid.New().String()
	c.Locals("request_id", requestID)

	err := c.Next()

	//ambil dahulu status code
	statusCode := c.Response().StatusCode()
	var errMessage string
	if err != nil {
		if e, ok := err.(*fiber.Error); ok {
			statusCode = e.Code
		} else if e, ok := err.(*errors.Response); ok {
			statusCode = e.HttpCode
		} else {
			statusCode = 500
		}
		errMessage = err.Error()
	}

	var userID *string
	if id, ok := c.Locals("user_id").(string); ok {
		userID = &id
	}

	m.restFullLogPublisher.Publish(request.RestFullLog{
		Timestamp:  time.Now().Format(time.RFC3339),
		RequestID:  requestID,
		UserID:     userID,
		Method:     c.Method(),
		Url:        c.OriginalURL(),
		Body:       c.Locals("request_body"),
		StatusCode: statusCode,
		ClientIP:   c.IP(),
		UserAgent:  c.Get("User-Agent"),
		Latency:    float64(time.Since(now).Milliseconds()),
		Error:      errMessage,
	})

	return err
}
