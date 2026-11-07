package server

import (
	"product-service/env"
	"product-service/src/common/log"
	"product-service/src/factory"
	"product-service/src/restful/handler"
	"product-service/src/restful/middleware"
	"product-service/src/restful/router"

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

	productHandler := handler.NewProductHandler(f)
	router.Product(App, m, productHandler)

	categoryHandler := handler.NewCategoryHandler(f)
	router.Category(App, m, categoryHandler)

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
