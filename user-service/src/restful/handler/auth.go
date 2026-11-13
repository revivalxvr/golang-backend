package handler

import (
	"encoding/base64"
	"time"

	"github.com/gofiber/fiber/v3"

	"user-service/src/common/dto/request"
	"user-service/src/common/dto/response"
	"user-service/src/common/utils"
	"user-service/src/factory"
	"user-service/src/restful/service"
)

type AuthHandler struct {
	service service.Auth
}

// function ini akan di panggil di server/restful.go
func NewAuthHandler(f *factory.Factory) *AuthHandler {
	return &AuthHandler{
		service: service.NewAuth(f),
	}
}

func (h *AuthHandler) Register(c fiber.Ctx) error {
	req := new(request.RegisterRequest)

	if err := c.Bind().Body(req); err != nil {
		return err
	}

	logBody := *req
	logBody.Password = "******"
	c.Locals("request_body", logBody)

	email, err := h.service.Register(c, req)
	if err != nil {
		return err
	}

	c.Cookie(&fiber.Cookie{
		Name:     "pending_register",
		Value:    base64.StdEncoding.EncodeToString([]byte(email)),
		Expires:  time.Now().Add(30 * time.Minute),
		HTTPOnly: true,
		Path:     "api/v1/auth/register/verify",
	})

	return c.Status(200).JSON(response.CommonResponse{
		Message: "success, please check your email for otp",
	})
}

func (h *AuthHandler) VerifyRegister(c fiber.Ctx) error {
	req := new(request.VerifyRegister)

	if err := c.Bind().Body(req); err != nil {
		return err
	}

	logBody := *req
	logBody.Otp = "******"
	c.Locals("request_body", logBody)

	//ambil email dari cookie
	email, err := base64.StdEncoding.DecodeString(c.Cookies("pending_register"))
	if err != nil {
		return err
	}
	//menyimpan email ke request
	req.Email = string(email)

	//meneruskan yang di simpan di req ke service
	if err := h.service.VerifyRegister(c, req); err != nil {
		return err
	}

	utils.ClearCookie(c, "pending_register", "api/v1/auth/register/verify")
	return c.Status(200).JSON(response.CommonResponse{
		Message: "verify register success",
	})
}

func (h *AuthHandler) Login(c fiber.Ctx) error {
	req := new(request.LoginRequest)

	if err := c.Bind().Body(req); err != nil {
		return err
	}

	logBody := *req
	logBody.Password = "******"
	c.Locals("request_body", logBody)

	user, err := h.service.Login(c, req)
	if err != nil {
		return err
	}
	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    user.AccessToken,
		Expires:  time.Now().Add(15 * time.Minute),
		HTTPOnly: true,
		Path:     "/",
	})

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    user.RefreshToken,
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
		Path:     "/",
	})

	return c.Status(200).JSON(response.CommonResponse{
		Message: "login success",
		Data:    user.Data,
	})
}

func (h *AuthHandler) Logout(c fiber.Ctx) error {
	refreshToken := c.Cookies("refresh_token")

	if err := h.service.Logout(c, refreshToken); err != nil {
		return err
	}
	utils.ClearCookie(c, "access_token", "/")
	utils.ClearCookie(c, "refresh_token", "/")
	return c.Status(200).JSON(response.CommonResponse{
		Message: "logout success",
	})
}

func (h *AuthHandler) RefreshToken(c fiber.Ctx) error {
	refreshToken := c.Cookies("refresh_token")

	res, err := h.service.RefreshToken(c, refreshToken)
	if err != nil {
		return err
	}
	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    res.AccessToken,
		Expires:  time.Now().Add(1 * time.Hour),
		HTTPOnly: true,
		Path:     "/",
	})
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    res.RefreshToken,
		Expires:  time.Now().Add(7 * time.Hour),
		HTTPOnly: true,
		Path:     "/",
	})
	return c.Status(200).JSON(response.CommonResponse{
		Message: "refresh token success",
	})
}
