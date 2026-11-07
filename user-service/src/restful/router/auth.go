package router

import (
	"user-service/env" // Sesuaikan dengan path import package env Anda
	"user-service/src/restful/handler"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
)

func Auth(app *fiber.App, h *handler.AuthHandler) {
	origin := env.Conf.CurrentApp.AllowedOrigin

	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{origin},
		AllowMethods:     []string{"GET", "POST", "DELETE", "PUT", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	}))

	app.Post("/api/v1/auth/register", h.Register)
	app.Post("/api/v1/auth/register/verify", h.VerifyRegister)
	app.Post("/api/v1/auth/login", h.Login)
	app.Delete("/api/v1/auth/logout", h.Logout)
	app.Patch("/api/v1/auth/refresh-token", h.RefreshToken)
}
