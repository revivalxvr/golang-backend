package middleware

import (
	"crypto/sha512"
	"encoding/hex"
	"notification-service/env"
	"notification-service/src/common/dto/request"
	"notification-service/src/common/errors"

	"github.com/gofiber/fiber/v3"
)

func (m *Middleware) ValidateMidtransTx(c fiber.Ctx) error {
	req := new(request.MidtransTx)
	if err := c.Bind().Body(req); err != nil {
		return err
	}
	//SHA512(order_id+status_code+gross_amount+ServerKey)
	key := req.OrderId + req.StatusCode + req.GrossAmount + env.Conf.Midtrans.ServerKey

	hash := sha512.New()
	hash.Write([]byte(key))
	signature := hash.Sum(nil)
	sigratureHashString := hex.EncodeToString(signature)

	if sigratureHashString != req.SingratureKey {
		return &errors.Response{
			HttpCode: 400,
			Message:  "invalid signature key",
		}
	}

	req.OrderId = req.Metadata.OriginalOrderId

	c.Locals("midtrans_tx", req)
	return c.Next()
}
