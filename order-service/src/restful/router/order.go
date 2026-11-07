package router

import (
	"order-service/env"
	"order-service/src/restful/handler"
	"order-service/src/restful/middleware"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
)

func OrderRouter(app *fiber.App, h *handler.Order, m *middleware.Middleware) {
	origin := env.Conf.CurrentApp.AllowedOrigin

	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{origin},
		AllowMethods:     []string{"GET", "POST", "DELETE", "PUT", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	}))
	app.Post("/api/v1/orders", m.AccessToken, m.UserRole, h.CreateOrder)   //cuma role USER yang bisa create order
	app.Get("/api/v1/orders", m.AccessToken, m.UserRole, h.GetOrderByUser) //cuma role USER yang bisa get order by user
}
