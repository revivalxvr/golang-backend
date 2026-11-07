package service

import (
	"context"
	"product-service/src/factory"
	"product-service/src/repository"

	productpb "github.com/revivalxvr/protobuf/protogen/product"
)

type ProductService interface {
	ReduceStocks(ctx context.Context, req *productpb.ReduceStocksReq) error
	RollbackStocks(ctx context.Context, req *productpb.RollbackStocksReq) error
}

type productImpl struct {
	productRepository repository.ProductRepository
}

func NewProductService(f *factory.Factory) ProductService {
	return &productImpl{
		productRepository: f.ProductRepository,
	}
}

func (p *productImpl) ReduceStocks(ctx context.Context, req *productpb.ReduceStocksReq) error {
	return p.productRepository.ReduceStocks(ctx, req.ProductOrders)
}
func (p *productImpl) RollbackStocks(ctx context.Context, req *productpb.RollbackStocksReq) error {
	return p.productRepository.RollbackStocks(ctx, req.ProductOrders)
}
