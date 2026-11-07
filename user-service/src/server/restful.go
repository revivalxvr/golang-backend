package server

import (
	"user-service/env"
	"user-service/src/common/log"
	"user-service/src/factory"
	"user-service/src/restful/handler"
	"user-service/src/restful/middleware"
	"user-service/src/restful/router"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
)

type Restful struct {
	App *fiber.App
}

// function ini akan di panggil di server/main.go
func NewRestful(f *factory.Factory) *Restful {
	m := middleware.NewMiddleware(f) //panggil function dari middleware

	App := fiber.New(fiber.Config{
		ErrorHandler: m.ErrorHandler,
	})

	App.Use(logger.New())
	App.Use(m.Log)
	App.Use(recover.New())

	//panggil function NewAuthHandler yg ada di handler
	authHandler := handler.NewAuthHandler(f)

	router.Auth(App, authHandler)

	//panggil function NewProfile yg ada di handler
	profileHandler := handler.NewProfile(f)
	router.Profile(App, profileHandler, m)

	return &Restful{
		App: App,
	}
}

func (s *Restful) Start() {
	if err := s.App.Listen(env.Conf.CurrentApp.RestfulAddr); err != nil {
		log.Logger.Fatal(err.Error())

	}
}

func (r *Restful) Stop() {
	if err := r.App.Shutdown(); err != nil {
		log.Logger.Fatal(err.Error())
	}
}
