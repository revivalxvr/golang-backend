package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"user-service/src/common/dto/request"
	"user-service/src/common/dto/response"
	"user-service/src/common/model"

	"github.com/redis/go-redis/v9"
)

// jgn lupa tambahkan Cache ke factory
type Cache interface {
	Set(ctx context.Context, key string, value any, expiration time.Duration) error
	GetSendOtp(ctx context.Context, key string) (*request.SendOtpReq, error)
	GetRegister(ctx context.Context, key string) (*request.RegisterRequest, error)
	GetUser(ctx context.Context, userID string) (*model.User, error)
	Delete(ctx context.Context, key string) error
	GetUserBlockInfo(ctx context.Context, userID string) (*response.UserBlockResponse, error)
}

type cacheImpl struct {
	redis *redis.ClusterClient
}

func NewCache(redis *redis.ClusterClient) Cache {
	return &cacheImpl{
		redis: redis,
	}
}

func (c *cacheImpl) Set(ctx context.Context, key string, value any, expired time.Duration) error {
	jsonData, err := json.Marshal(value)
	if err != nil && err != redis.Nil {
		return err
	}
	return c.redis.Set(ctx, key, jsonData, expired).Err()
}

func (c *cacheImpl) GetSendOtp(ctx context.Context, key string) (*request.SendOtpReq, error) {
	data, err := c.redis.Get(ctx, key).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}
	if data == "" {
		return nil, nil
	}
	sendOtpReq := new(request.SendOtpReq)
	if err := json.Unmarshal([]byte(data), sendOtpReq); err != nil {
		return nil, err
	}

	return sendOtpReq, nil

}

func (c *cacheImpl) GetRegister(ctx context.Context, key string) (*request.RegisterRequest, error) {
	data, err := c.redis.Get(ctx, key).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}
	if data == "" {
		return nil, nil
	}
	registerReq := new(request.RegisterRequest)
	if err := json.Unmarshal([]byte(data), registerReq); err != nil {
		return nil, err
	}

	return registerReq, nil
}

func (c *cacheImpl) GetUser(ctx context.Context, userID string) (*model.User, error) {
	data, err := c.redis.Get(ctx, fmt.Sprintf("user:%s", userID)).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}
	if data == "" {
		return nil, nil
	}
	registerReq := new(model.User)
	if err := json.Unmarshal([]byte(data), registerReq); err != nil {
		return nil, err
	}

	return registerReq, nil
}
func (c *cacheImpl) GetUserBlockInfo(ctx context.Context, userID string) (*response.UserBlockResponse, error) {
	data, err := c.redis.Get(ctx, fmt.Sprintf("user_block_info:%s", userID)).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}
	if data == "" {
		return nil, nil
	}
	userBlockRes := new(response.UserBlockResponse)
	if err := json.Unmarshal([]byte(data), userBlockRes); err != nil {
		return nil, err
	}

	return userBlockRes, nil
}

func (c *cacheImpl) Delete(ctx context.Context, key string) error {
	return c.redis.Del(ctx, key).Err()
}
