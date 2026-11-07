package server

import (
	"notification-service/env"
	"notification-service/src/common/log"
	"notification-service/src/factory"
	"notification-service/src/restful/handler"
	"notification-service/src/restful/middleware"
	"notification-service/src/restful/router"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
)

type Restful struct {
	App *fiber.App
}

func NewRestful(f *factory.Factory) *Restful {
	m := middleware.NewMiddleware()

	App := fiber.New(fiber.Config{
		ErrorHandler: m.ErrorHandler,
	})

	App.Use(logger.New())
	App.Use(recover.New())

	notifHandler := handler.NewNotifHandler(f)
	router.Notif(App, notifHandler, m)

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
