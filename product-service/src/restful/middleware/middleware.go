package middleware

import (
	"product-service/src/factory"
	"product-service/src/grpc/client"
	"product-service/src/publisher"
)

type Middleware struct {
	userClient           client.User
	restFullLogPublisher *publisher.Kafka
}

func NewMiddleware(f *factory.Factory) *Middleware {
	return &Middleware{
		userClient:           f.UserClient,
		restFullLogPublisher: f.RestFullLogPublisher,
	}
}
