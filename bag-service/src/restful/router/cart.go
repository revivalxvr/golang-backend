package router

import (
	"bag-service/env"
	"bag-service/src/restful/handler"
	"bag-service/src/restful/middleware"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
)

func CartRouter(app *fiber.App, h *handler.Cart, m *middleware.Middleware) {
	origin := env.Conf.CurrentApp.AllowedOrigin

	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{origin},
		AllowMethods:     []string{"GET", "POST", "DELETE", "PUT", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	}))
	app.Post("/api/v1/cart", m.AccessToken, m.UserRole, h.AddCart)                      //cuma bisa login dan role harus USER
	app.Get("/api/v1/cart", m.AccessToken, m.UserRole, h.GetCart)                       //cuma bisa login dan role harus USER
	app.Delete("/api/v1/cart/:product_id", m.AccessToken, m.UserRole, h.DeleteCartItem) //cuma bisa login dan role harus USER
}
