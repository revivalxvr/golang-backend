package handler

import (
	"context"
	"product-service/src/factory"
	"product-service/src/grpc/service"

	pb "github.com/revivalxvr/protobuf/protogen/cart"
)

type cartImpl struct {
	service service.CartService
	pb.UnimplementedCartServiceServer
}

func NewCart(f *factory.Factory) pb.CartServiceServer {
	return &cartImpl{
		service: service.NewCartService(f),
	}
}

func (c *cartImpl) GetProductsForCart(ctx context.Context, req *pb.GetProductsForCartReq) (*pb.GetProductsForCartRes, error) {
	return c.service.GetProductsForCart(ctx, req)
}
