package mockgrpcclient

import (
	"context"

	pb "github.com/revivalxvr/protobuf/protogen/product"
	"github.com/stretchr/testify/mock"
)

type Product struct {
	mock.Mock
}

func NewProduct() *Product {
	return &Product{
		Mock: mock.Mock{},
	}
}

func (p *Product) ReduceStocks(ctx context.Context, data *pb.ReduceStocksReq) error {
	args := p.Called(ctx, data)
	return args.Error(0)
}
func (p *Product) RollbackStock(ctx context.Context, data *pb.RollbackStocksReq) error {
	args := p.Called(ctx, data)
	return args.Error(0)
}

func (p *Product) Close() {}
