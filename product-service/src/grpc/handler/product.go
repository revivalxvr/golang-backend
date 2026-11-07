package handler

import (
	"context"
	"product-service/src/factory"
	"product-service/src/grpc/service"

	productpb "github.com/revivalxvr/protobuf/protogen/product"
	"google.golang.org/protobuf/types/known/emptypb"
)

type productImpl struct {
	service service.ProductService
	productpb.UnimplementedProductServiceServer
}

func NewProduct(f *factory.Factory) productpb.ProductServiceServer {
	return &productImpl{
		service: service.NewProductService(f),
	}
}

func (p *productImpl) ReduceStocks(ctx context.Context, req *productpb.ReduceStocksReq) (*emptypb.Empty, error) {
	err := p.service.ReduceStocks(ctx, req)
	return &emptypb.Empty{}, err
}
func (p *productImpl) RollbackStocks(ctx context.Context, req *productpb.RollbackStocksReq) (*emptypb.Empty, error) {
	err := p.service.RollbackStocks(ctx, req)
	return &emptypb.Empty{}, err
}
