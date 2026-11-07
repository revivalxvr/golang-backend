package factory

import (
	"bag-service/src/broker/publisher"
	grpcclinet "bag-service/src/grpc/client"
	"bag-service/src/repository"

	"gorm.io/gorm"
)

type Factory struct {
	Cartrepository       repository.CartRepository
	UserClient           grpcclinet.User
	CartClient           grpcclinet.Cart
	RestFullLogPublisher *publisher.Kafka
}

func NewFactory(db *gorm.DB, uc grpcclinet.User, ct grpcclinet.Cart, rp *publisher.Kafka) *Factory {
	return &Factory{
		Cartrepository:       repository.NewCart(db),
		UserClient:           uc,
		CartClient:           ct,
		RestFullLogPublisher: rp,
	}
}
