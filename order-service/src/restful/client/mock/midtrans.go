package mockrestfulclient

import (
	"context"
	"order-service/src/common/dto/request"
	"order-service/src/common/dto/response"

	"github.com/stretchr/testify/mock"
)

type Midtrans struct {
	mock.Mock
}

func NewMidtrans() *Midtrans {
	return &Midtrans{
		Mock: mock.Mock{},
	}
}

func (c *Midtrans) CreateTransaction(ctx context.Context, order *request.MidtransTransaction) (*response.MidtransTx, error) {
	arguments := c.Mock.Called(ctx, order)

	if arguments.Get(0) == nil {
		return nil, arguments.Error(1)
	}

	return arguments.Get(0).(*response.MidtransTx), arguments.Error(1)
}
