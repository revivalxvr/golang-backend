package middleware

import (
	"order-service/src/broker/publisher"
	"order-service/src/factory"
	"order-service/src/grpc/client"
)

type Middleware struct {
	restFullLogPublisher *publisher.Kafka
	userClient           client.User
}

func NewMiddleware(f *factory.Factory) *Middleware {
	return &Middleware{
		restFullLogPublisher: f.RestfullLogPublisher,
		userClient:           f.UserClient,
	}
}
