package client

import (
	"context"
	"crypto/tls"
	"order-service/env"
	"order-service/src/common/log"
	"order-service/src/common/pkg/cbreaker"
	"order-service/src/grpc/interceptor"

	pbi "github.com/revivalxvr/protobuf/protogen/cart"
	"github.com/sony/gobreaker/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

type Cart interface {
	GetProductsForCart(ctx context.Context, req *pbi.GetProductsForCartReq) (*pbi.GetProductsForCartRes, error)
	Close()
}

type cartImpl struct {
	client pbi.CartServiceClient
	conn   *grpc.ClientConn
	cb     *gobreaker.CircuitBreaker[any]
}

func NewCart(i *interceptor.Unary, tlsConf *tls.Config) Cart {
	cb := cbreaker.NewGrpc("product service-grpc")

	var opts []grpc.DialOption
	opts = append(
		opts,
		grpc.WithTransportCredentials(credentials.NewTLS(tlsConf)),
		grpc.WithChainUnaryInterceptor(
			i.AddMetaData,
		),
	)
	//pastikan koneksi nya ke product service
	conn, err := grpc.NewClient(env.Conf.ProductService.GrpcAddr, opts...)
	if err != nil {
		log.Logger.Fatal(err.Error())
	}
	client := pbi.NewCartServiceClient(conn)
	return &cartImpl{
		client: client,
		conn:   conn,
		cb:     cb,
	}
}
func (p *cartImpl) GetProductsForCart(ctx context.Context, req *pbi.GetProductsForCartReq) (*pbi.GetProductsForCartRes, error) {
	res, err := p.cb.Execute(func() (any, error) {
		return p.client.GetProductsForCart(ctx, req)
	})
	if err != nil {
		return nil, err
	}
	result, ok := res.(*pbi.GetProductsForCartRes)
	if !ok {
		return nil, status.Errorf(codes.Internal, "invalid response type")
	}
	return result, nil
}

func (p *cartImpl) Close() {
	if err := p.conn.Close(); err != nil {
		log.Logger.Error(err.Error())
	}
}
