package factory

import (
	"notification-service/src/publisher"
)

type Factory struct {
	RabbitMQPublisher *publisher.RabitMQ
}

func NewFactory(pr *publisher.RabitMQ) *Factory {
	return &Factory{
		RabbitMQPublisher: pr,
	}
}
