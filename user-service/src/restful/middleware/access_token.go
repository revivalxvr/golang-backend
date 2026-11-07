package middleware

import (
	"context"
	"fmt"
	"time"
	"user-service/env"
	"user-service/src/common/dto/response"
	"user-service/src/common/errors"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

func (m *Middleware) AccessToken(c fiber.Ctx) error {
	accessToken := c.Cookies("access_token")

	if accessToken == "" {
		return &errors.Response{Message: "unauthorized", HttpCode: 401}
	}

	jwtToken, err := jwt.Parse(accessToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, &errors.Response{Message: "invalid token", HttpCode: 401}
		}
		return env.Conf.Jwt.PublicKey, nil
	})
	if err != nil {
		return &errors.Response{Message: err.Error(), HttpCode: 401}
	}
	if claims, ok := jwtToken.Claims.(jwt.MapClaims); ok && jwtToken.Valid {
		userID, ok := claims["user_id"].(string)
		if !ok || userID == "" {
			return &errors.Response{Message: "invalid token claims", HttpCode: 401}
		}

		c.Locals("user_id", userID)
		c.Locals("role", claims["role"])

		userBlockInfo, err := m.checkUserBlock(c.Context(), userID)
		if err != nil {
			return &errors.Response{HttpCode: 500, Message: err.Error()}
		}
		if userBlockInfo.IsBlocked {
			return &errors.Response{HttpCode: 403, Message: fmt.Sprintf("User blocked: %s", userBlockInfo.Reason)}
		}
	}
	return c.Next()
}

func (m *Middleware) checkUserBlock(ctx context.Context, userID string) (*response.UserBlockResponse, error) {
	userBlockInfo, err := m.cacheRepository.GetUserBlockInfo(ctx, userID)
	if err != nil {
		return nil, err
	}

	if userBlockInfo != nil {
		return &response.UserBlockResponse{
			IsBlocked: userBlockInfo.IsBlocked, Reason: userBlockInfo.Reason}, nil
	}

	userBlockLog, err := m.userBlockLogRepository.Find(ctx, userID)
	if err != nil {
		return nil, err
	}

	var isBlocked bool
	var reason string

	if userBlockLog != nil {
		isBlocked = true
		reason = userBlockLog.Reason
	}

	err = m.cacheRepository.Set(ctx, fmt.Sprintf("user_block_info:%s", userID), &response.UserBlockResponse{IsBlocked: isBlocked, Reason: reason}, 30*time.Minute)
	if err != nil {
		return nil, err
	}

	return &response.UserBlockResponse{IsBlocked: isBlocked, Reason: reason}, nil
}
