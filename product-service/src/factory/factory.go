package factory

import (
	grpcclinet "product-service/src/grpc/client"
	"product-service/src/publisher"
	"product-service/src/repository"
	"product-service/src/restful/client"

	"gorm.io/gorm"
)

type Factory struct {
	ImageKitClient       *client.ImageKit
	ProductRepository    repository.ProductRepository
	Cartrepository       repository.CartRepository
	CategoryRepository   repository.CategoryRepository
	GrpcLogPublisher     *publisher.Kafka
	UserClient           grpcclinet.User
	RestFullLogPublisher *publisher.Kafka
}

func NewFactory(
	db *gorm.DB,
	kf *publisher.Kafka,
	uc grpcclinet.User,
	rp *publisher.Kafka,
) *Factory {
	return &Factory{
		ProductRepository:    repository.NewProduct(db),
		Cartrepository:       repository.NewCart(db),
		CategoryRepository:   repository.NewCategory(db),
		ImageKitClient:       client.NewImageKit(),
		GrpcLogPublisher:     kf,
		UserClient:           uc,
		RestFullLogPublisher: rp,
	}
}
