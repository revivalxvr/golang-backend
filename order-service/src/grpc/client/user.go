package client

import (
	"context"
	"crypto/tls"
	"order-service/env"
	"order-service/src/common/log"
	"order-service/src/common/pkg/cbreaker"
	"order-service/src/grpc/interceptor"

	pb "github.com/revivalxvr/protobuf/protogen/user"
	"github.com/sony/gobreaker/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

type User interface {
	CheckUserBlocked(ctx context.Context, req *pb.CheckUserBlockedReq) (*pb.CheckUserBlockedRes, error)
	Close()
}
type userImpl struct {
	client pb.UserServiceClient
	conn   *grpc.ClientConn
	cb     *gobreaker.CircuitBreaker[any]
}

func NewUser(i *interceptor.Unary, tlsConf *tls.Config) User {
	cb := cbreaker.NewGrpc("user service-grpc") //untuk circuit breaker, nama nya user service-grpc

	var opts []grpc.DialOption //kalau jadi clinet kita harus DialOption (liat di folder server product-service)
	opts = append(
		opts,
		// grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithTransportCredentials(credentials.NewTLS(tlsConf)),
		grpc.WithChainUnaryInterceptor(
			// i.AddBasicAuth, // add meta data di interceptor
			i.AddMetaData, //add meta data
		),
	)
	conn, err := grpc.NewClient(env.Conf.UserService.GrpcAddr, opts...)
	if err != nil {
		log.Logger.Fatal(err.Error())
	}
	client := pb.NewUserServiceClient(conn)
	return &userImpl{
		client: client,
		conn:   conn,
		cb:     cb,
	}
}

func (p *userImpl) CheckUserBlocked(ctx context.Context, req *pb.CheckUserBlockedReq) (*pb.CheckUserBlockedRes, error) {
	res, err := p.cb.Execute(func() (any, error) {
		return p.client.CheckUserBlocked(ctx, req)
	})
	if err != nil {
		return nil, err
	}
	result, ok := res.(*pb.CheckUserBlockedRes)
	if !ok {
		return nil, status.Errorf(codes.Internal, "invalid response type")
	}
	return result, nil
}

func (p *userImpl) Close() {
	if err := p.conn.Close(); err != nil {
		log.Logger.Error(err.Error())
	}
}
