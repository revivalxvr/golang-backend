package router

import (
	"product-service/env"
	"product-service/src/restful/handler"
	"product-service/src/restful/middleware"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
)

func Product(app *fiber.App, m *middleware.Middleware, h *handler.Product) {

	origin := env.Conf.CurrentApp.AllowedOrigin

	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{origin},
		AllowMethods:     []string{"GET", "POST", "DELETE", "PUT", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	}))
	app.Post("/api/v1/products", m.AccessToken, m.Admin, m.SaveTemporaryProductImage, m.ValidateProductImage, h.Create)
	app.Put("/api/v1/products/:id", m.AccessToken, m.Admin, m.SaveTemporaryProductImage, m.ValidateProductImage, h.Update)
	app.Delete("/api/v1/products/:id", m.AccessToken, m.Admin, h.Delete)
	app.Get("/api/v1/products", h.FindMany)
}
