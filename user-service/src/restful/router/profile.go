package router

import (
	"user-service/env"
	"user-service/src/restful/handler"
	"user-service/src/restful/middleware"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
)

func Profile(app *fiber.App, h *handler.Profile, m *middleware.Middleware) {
	origin := env.Conf.CurrentApp.AllowedOrigin

	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{origin},
		AllowMethods:     []string{"GET", "POST", "DELETE", "PUT", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	}))

	app.Get("/api/v1/profile", m.AccessToken, h.GetProfile)
	//update profile for admin and user
	app.Put("/api/v1/profile-admin", m.AccessToken, m.Admin, h.UpdateProfileAdmin)
	app.Put("/api/v1/profile-user", m.AccessToken, m.UserRole, h.UpdateProfileUser)

	//update password for admin and user
	app.Patch("/api/v1/profile-admin/password", m.AccessToken, m.Admin, h.UpdatePasswordAdmin)
	app.Patch("/api/v1/profile-user/password", m.AccessToken, m.UserRole, h.UpdatePasswordUser)
}
