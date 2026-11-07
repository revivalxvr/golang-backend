package router

import (
	"product-service/env"
	"product-service/src/restful/handler"
	"product-service/src/restful/middleware"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
)

func Category(app *fiber.App, m *middleware.Middleware, h *handler.Category) {

	origin := env.Conf.CurrentApp.AllowedOrigin

	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{origin},
		AllowMethods:     []string{"GET", "POST", "DELETE", "PUT", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	}))
	app.Post("/api/v1/categories", m.AccessToken, m.Admin, h.CreateCategory)
	app.Get("/api/v1/categories", h.FindCategory)
	app.Put("/api/v1/categories/:id", m.AccessToken, m.Admin, h.UpdateCategory)
	app.Delete("/api/v1/categories/:id", m.AccessToken, m.Admin, h.DeleteCategory)

}
