package factory

import (
	"google.golang.org/api/gmail/v1"

	"email-service/src/broker/publisher"
	"email-service/src/common/pkg/oauth"
)

type Factory struct {
	GmailService      *gmail.Service
	RabbitMqPublisher *publisher.Kafka
}

func NewFactory(kf *publisher.Kafka) *Factory {
	return &Factory{
		GmailService:      oauth.NewGmailService(),
		RabbitMqPublisher: kf,
	}
}
