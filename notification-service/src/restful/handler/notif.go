package handler

import (
	"encoding/json"
	"notification-service/src/common/dto/request"
	"notification-service/src/common/dto/response"
	"notification-service/src/common/errors"
	"notification-service/src/common/log"
	"notification-service/src/factory"
	"notification-service/src/publisher"

	"github.com/gofiber/fiber/v3"
)

type NotifHandler struct {
	rabbitMQ *publisher.RabitMQ
}

func NewNotifHandler(f *factory.Factory) *NotifHandler {
	return &NotifHandler{
		rabbitMQ: f.RabbitMQPublisher,
	}
}

func (n *NotifHandler) MidtransTx(c fiber.Ctx) error {
	midtransTx, ok := c.Locals("midtrans_tx").(*request.MidtransTx)
	if !ok {
		return &errors.Response{
			HttpCode: 400,
			Message:  "failed to validate midtrans transaction",
		}
	}
	b, _ := json.MarshalIndent(midtransTx, "", "  ")
	log.Logger.Info("midtrans transaction: " + string(b))

	// Publish event ke RabbitMQ
	// Sesuaikan Exchange Name dan Keyname dengan yang di consumer order-service
	exchange := "notification"
	routingKey := "midtrans-tx"

	if err := n.rabbitMQ.Publish(exchange, routingKey, midtransTx); err != nil {
		log.Logger.Error("failed to publish notification to rabbitmq: " + err.Error())
		return &errors.Response{
			HttpCode: 500,
			Message:  "failed to publish notification to message broker",
		}
	}

	log.Logger.Info("successfully published notification to RabbitMQ exchange: " + exchange)

	return c.Status(200).JSON(response.CommonResponse{
		Message: "notification sent successfully",
	})

}
