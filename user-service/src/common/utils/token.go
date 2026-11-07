package utils

import (
	"time"
	"user-service/env"

	"github.com/golang-jwt/jwt/v5"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

func GenerateAccessToken(userID string, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(30 * time.Minute).Unix(),
	})

	accessToken, err := token.SignedString(env.Conf.Jwt.PrivateKey)
	if err != nil {
		return "", err
	}
	return accessToken, nil
}

func GenerateRefreshToken(userID string) (string, error) {
	id, err := gonanoid.New()
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"id":      id,
		"user_id": userID,
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(),
	})

	refreshToken, err := token.SignedString(env.Conf.Jwt.PrivateKey)
	if err != nil {
		return "", err
	}

	return refreshToken, nil
}
