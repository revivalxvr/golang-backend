package server

import (
	"order-service/env"
	"order-service/src/common/log"
	"order-service/src/factory"
	"order-service/src/restful/handler"
	"order-service/src/restful/middleware"
	"order-service/src/restful/router"

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

	orderHandler := handler.NewOrder(f)
	router.OrderRouter(App, orderHandler, m)

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
