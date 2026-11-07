package middleware

import (
	"context"
	"fmt"
	"order-service/env"

	"order-service/src/common/constant"
	"order-service/src/common/errors"

	pb "github.com/revivalxvr/protobuf/protogen/user"

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

		ctx := context.WithValue(c.Context(), constant.RequestID, c.Locals("request_id"))
		ctx = context.WithValue(ctx, constant.UserID, userID)

		res, err := m.userClient.CheckUserBlocked(ctx, &pb.CheckUserBlockedReq{
			UserId: userID,
		})
		if err != nil {
			return &errors.Response{Message: err.Error(), HttpCode: 500}
		}
		if res.IsBlocked {
			return &errors.Response{HttpCode: 403, Message: fmt.Sprintf("User blocked: %s", res.Reason)}
		}
	}
	return c.Next()
}
