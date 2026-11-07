package handler

import (
	"log"
	"order-service/src/broker/publisher"
	"order-service/src/common/dto/event"
	"order-service/src/common/dto/request"
	"order-service/src/factory"
	"order-service/src/restful/service"

	"order-service/src/common/dto/response"

	"github.com/google/uuid"

	"github.com/gofiber/fiber/v3"
)

type Order struct {
	service  service.Order
	rabbitMQ *publisher.RabitMQ
}

func NewOrder(f *factory.Factory) *Order {
	return &Order{
		service:  service.NewOrder(f),
		rabbitMQ: f.RabbitMQPublisher,
	}
}

func (h *Order) CreateOrder(c fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok {
		return fiber.NewError(
			fiber.StatusUnauthorized,
			"invalid user ID",
		)
	}
	req := new(request.CreateOrder)

	if err := c.Bind().Body(req); err != nil {
		return err
	}

	// 1. Jalankan proses pembuatan order.
	res, err := h.service.CreateOrder(c, req, userID)
	if err != nil {
		return err
	}
	// 2. Ambil product ID dari item yang di-checkout.
	productIDs := make([]string, 0, len(req.ProductOrders))
	for _, product := range req.ProductOrders {
		productIDs = append(productIDs, product.ProductID)
	}
	// 3. Siapkan event order.created untuk RabbitMQ
	orderEvent := event.OrderCreated{
		EventID:    uuid.NewString(),
		OrderID:    req.Order.ID,
		UserID:     userID,
		ProductIDs: productIDs,
	}

	// 4. Publish event ke RabbitMQ
	if err := h.rabbitMQ.Publish(
		"order",
		"order.created",
		orderEvent,
	); err != nil {
		log.Printf(
			"failed to publish order.created event: %v",
			err,
		)

		// Order sudah berhasil dibuat.
		// Jangan menganggap order gagal hanya karena publish gagal.
		// Untuk production, gunakan Transactional Outbox.
	}
	// 5. Kembalikan respons kepada user.
	return c.Status(200).JSON(response.CommonResponse{
		Message: "success create order",
		Data:    res,
	})
}

// GetOrderByUser
func (h *Order) GetOrderByUser(c fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok {
		return fiber.NewError(
			fiber.StatusUnauthorized,
			"invalid user ID",
		)
	}

	res, err := h.service.GetOrderByUser(c, userID)
	if err != nil {
		return err
	}

	return c.Status(200).JSON(response.CommonResponse{
		Message: "success get order by user",
		Data:    res,
	})
}
