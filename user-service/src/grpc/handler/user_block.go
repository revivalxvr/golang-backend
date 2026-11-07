package handler

import (
	"context"
	"user-service/src/factory"
	"user-service/src/grpc/service"

	pb "github.com/revivalxvr/protobuf/protogen/user"
)

type userImpl struct {
	service service.UserService
	pb.UnimplementedUserServiceServer
}

func NewUser(f *factory.Factory) pb.UserServiceServer {
	return &userImpl{
		service: service.NewUserService(f),
	}
}

func (p *userImpl) CheckUserBlocked(ctx context.Context, req *pb.CheckUserBlockedReq) (*pb.CheckUserBlockedRes, error) {
	return p.service.CheckUserBlocked(ctx, req)
}
