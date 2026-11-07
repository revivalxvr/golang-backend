package service

import (
	"context"
	"fmt"
	"time"
	"user-service/src/cache"
	"user-service/src/factory"
	"user-service/src/repository"

	pb "github.com/revivalxvr/protobuf/protogen/user"
)

type UserService interface {
	CheckUserBlocked(ctx context.Context, req *pb.CheckUserBlockedReq) (*pb.CheckUserBlockedRes, error)
}

type userImpl struct {
	UserBlockRepository repository.UserBlock
	cacheUserBlocked    cache.Cache
}

func NewUserService(f *factory.Factory) UserService {
	return &userImpl{
		UserBlockRepository: f.UserBlockRepository,
		cacheUserBlocked:    f.CacheRepository,
	}
}

func (p *userImpl) CheckUserBlocked(ctx context.Context, req *pb.CheckUserBlockedReq) (*pb.CheckUserBlockedRes, error) {
	userBlockedInfo, err := p.cacheUserBlocked.GetUserBlockInfo(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	if userBlockedInfo != nil {
		return &pb.CheckUserBlockedRes{
			IsBlocked: userBlockedInfo.IsBlocked,
			Reason:    userBlockedInfo.Reason,
		}, nil
	}

	userBlock, err := p.UserBlockRepository.Find(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	var isBlocked bool
	var reason string

	if userBlock != nil {
		isBlocked = true
		reason = userBlock.Reason
	} else {
		isBlocked = false
		reason = ""
	}

	err = p.cacheUserBlocked.Set(ctx, fmt.Sprint("user block info:", req.UserId), &pb.CheckUserBlockedRes{
		IsBlocked: isBlocked,
		Reason:    reason,
	}, 5*time.Minute)
	if err != nil {
		return nil, err
	}

	if userBlock == nil {
		return &pb.CheckUserBlockedRes{
			IsBlocked: false,
			Reason:    "",
		}, nil
	}
	return &pb.CheckUserBlockedRes{
		IsBlocked: true,
		Reason:    userBlock.Reason,
	}, nil
}
