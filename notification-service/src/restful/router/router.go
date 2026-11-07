package router

import (
	"notification-service/src/restful/handler"
	"notification-service/src/restful/middleware"

	"github.com/gofiber/fiber/v3"
)

func Notif(app *fiber.App, h *handler.NotifHandler, m *middleware.Middleware) {
	app.Post("/api/v1/notif/midtrans-tx", m.ValidateMidtransTx, h.MidtransTx)
}
