package server

import (
	"crypto/tls"
	"net"
	"user-service/src/factory"
	"user-service/src/grpc/handler"
	"user-service/src/grpc/interceptor"

	"user-service/env"
	"user-service/src/common/log"

	pb "github.com/revivalxvr/protobuf/protogen/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type GrpcServer struct {
	server      *grpc.Server
	userHandler pb.UserServiceServer
}

func NewGrpcServer(f *factory.Factory, tlsConf *tls.Config) *GrpcServer {
	userHandler := handler.NewUser(f)

	i := interceptor.NewUnary(f)

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			// i.BasicAuthValidator,
			i.Logging,
			i.Recovery,
			i.Error,
		),
		grpc.Creds(credentials.NewTLS(tlsConf)),
	)

	pb.RegisterUserServiceServer(grpcServer, userHandler)

	//jgn gunakan di production karena mengekspos semua method gRPC ke public (testing di postman aja)
	//  reflection.Register(grpcServer)

	return &GrpcServer{
		server:      grpcServer,
		userHandler: userHandler,
	}
}

func (s *GrpcServer) Run() {
	lis, err := net.Listen("tcp", env.Conf.CurrentApp.GrpcAddr) //kalau jadi server harus Listen bukan Dial(liat di folder client order-service)
	if err != nil {
		log.Logger.Fatal(err.Error())
	}
	log.Logger.Info("USER gRPC server is running on " + env.Conf.CurrentApp.GrpcAddr)
	if err := s.server.Serve(lis); err != nil {
		log.Logger.Fatal(err.Error())
	}
}

func (s *GrpcServer) Stop() {
	s.server.Stop()

}
