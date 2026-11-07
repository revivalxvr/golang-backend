package server

import (
	"bag-service/env"
	"bag-service/src/common/log"
	"bag-service/src/factory"
	"bag-service/src/restful/handler"
	"bag-service/src/restful/middleware"
	"bag-service/src/restful/router"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
)

type Restful struct {
	App *fiber.App
}

func NewRestful(f *factory.Factory) *Restful {
	m := middleware.NewMiddleware(f)

	App := fiber.New(fiber.Config{
		ErrorHandler: m.ErrorHandler,
	})

	App.Use(logger.New())
	App.Use(m.Log)
	App.Use(recover.New())

	cartHandler := handler.NewAddCart(f)
	router.CartRouter(App, cartHandler, m)

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
