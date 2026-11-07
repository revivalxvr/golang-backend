package factory

import (
	"order-service/src/broker/publisher"
	grpcclinet "order-service/src/grpc/client"
	"order-service/src/repository"
	"order-service/src/restful/client"

	"gorm.io/gorm"
)

type Factory struct {
	MidtransClient         client.Midtrans
	OrderRepository        repository.Order
	ProductOrderRepository repository.ProductOrder
	TxRepository           repository.TxBeginner
	ProductClient          grpcclinet.Product
	RestfullLogPublisher   *publisher.Kafka
	UserClient             grpcclinet.User
	CartClient             grpcclinet.Cart
	RabbitMQPublisher      *publisher.RabitMQ
}

func NewFactory(db *gorm.DB, pc grpcclinet.Product, kf *publisher.Kafka, uc grpcclinet.User, pr *publisher.RabitMQ, cc grpcclinet.Cart) *Factory {
	return &Factory{
		MidtransClient:         client.NewMidtrans(),
		OrderRepository:        repository.NewOrder(db),
		ProductClient:          pc,
		TxRepository:           repository.NewTxBeginner(db),
		ProductOrderRepository: repository.NewProductOrder(db),
		RestfullLogPublisher:   kf,
		UserClient:             uc,
		CartClient:             cc,
		RabbitMQPublisher:      pr,
	}
}
