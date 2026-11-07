package server

import (
	"crypto/tls"
	"net"
	"product-service/src/factory"
	"product-service/src/grpc/handler"
	"product-service/src/grpc/interceptor"

	"product-service/env"
	"product-service/src/common/log"

	cr "github.com/revivalxvr/protobuf/protogen/cart"
	pb "github.com/revivalxvr/protobuf/protogen/product"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type GrpcServer struct {
	server         *grpc.Server
	productHandler pb.ProductServiceServer
	cartHandler    cr.CartServiceServer
}

func NewGrpcServer(f *factory.Factory, i *interceptor.Unary, tlsConf *tls.Config) *GrpcServer {
	productHandler := handler.NewProduct(f)
	cartHandler := handler.NewCart(f)
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			// i.BasicAuthValidator,
			i.Logging,
			i.Recovery,
			i.Error,
		),
		grpc.Creds(credentials.NewTLS(tlsConf)),
	)
	pb.RegisterProductServiceServer(grpcServer, productHandler)
	cr.RegisterCartServiceServer(grpcServer, cartHandler)

	//gunakan untuk mengecek method gRPC yg terdaftar(testing di postman aja hbs tu hapus)
	// reflection.Register(grpcServer)

	return &GrpcServer{
		server:         grpcServer,
		productHandler: productHandler,
		cartHandler:    cartHandler,
	}
}

func (s *GrpcServer) Run() {
	lis, err := net.Listen("tcp", env.Conf.CurrentApp.GrpcAddr) //kalau jadi server harus Listen bukan Dial(liat di folder client order-service)
	if err != nil {
		log.Logger.Fatal(err.Error())
	}
	log.Logger.Info("PRODUCT gRPC server is running on " + env.Conf.CurrentApp.GrpcAddr)
	if err := s.server.Serve(lis); err != nil {
		log.Logger.Fatal(err.Error())
	}
}

func (s *GrpcServer) Stop() {
	s.server.Stop()

}
