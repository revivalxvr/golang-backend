package middleware

import (
	"bag-service/src/broker/publisher"
	"bag-service/src/factory"
	"bag-service/src/grpc/client"
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
