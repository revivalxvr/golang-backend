package handler

import (
	"user-service/src/common/dto/request"
	"user-service/src/common/dto/response"
	v "user-service/src/common/pkg/validator"
	"user-service/src/factory"
	"user-service/src/restful/service"

	"github.com/gofiber/fiber/v3"
)

type Profile struct {
	service service.ProfileService
}

// function ini akan di panggil di server/restful.go
func NewProfile(f *factory.Factory) *Profile {
	return &Profile{
		service: service.NewProfileService(f),
	}
}

func (h *Profile) GetProfile(c fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	user, err := h.service.GetProfile(c.Context(), userID)
	if err != nil {
		return err
	}
	return c.Status(200).JSON(response.CommonResponse{
		Message: "success get profile",
		Data:    user,
	})
}
func (h *Profile) UpdateProfileAdmin(c fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok {
		return fiber.NewError(
			fiber.StatusUnauthorized,
			"invalid user ID",
		)
	}
	var req request.UpdateProfileRequest
	if err := c.Bind().Body(&req); err != nil {
		return err
	}
	// Jalankan validasi struct
	if err := v.Validate.Struct(req); err != nil {
		return err
	}
	res, err := h.service.UpdateProfileAdmin(c.Context(), userID, req)
	if err != nil {
		return err
	}
	return c.Status(200).JSON(response.CommonResponse{
		Message: "success update profile",
		Data:    res,
	})
}
func (h *Profile) UpdateProfileUser(c fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok {
		return fiber.NewError(
			fiber.StatusUnauthorized,
			"invalid user ID",
		)
	}
	var req request.UpdateProfileRequest
	if err := c.Bind().Body(&req); err != nil {
		return err
	}
	// Jalankan validasi struct
	if err := v.Validate.Struct(req); err != nil {
		return err
	}
	res, err := h.service.UpdateProfileUser(c.Context(), userID, req)
	if err != nil {
		return err
	}
	return c.Status(200).JSON(response.CommonResponse{
		Message: "success update profile",
		Data:    res,
	})
}

// h.UpdatePasswordAdmin)
func (h *Profile) UpdatePasswordAdmin(c fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok {
		return fiber.NewError(
			fiber.StatusUnauthorized,
			"invalid user ID",
		)
	}
	var req request.ChangePasswordRequest
	if err := c.Bind().Body(&req); err != nil {
		return err
	}
	// Jalankan validasi struct
	if err := v.Validate.Struct(req); err != nil {
		return err
	}
	res, err := h.service.UpdatePasswordAdmin(c.Context(), userID, req)
	if err != nil {
		return err
	}
	return c.Status(200).JSON(response.CommonResponse{
		Message: "success update password",
		Data: fiber.Map{
			"user_id":    res.UserID,
			"name":       res.User.Name,
			"updated_at": res.UpdatedAt,
		},
	})
}

// h.UpdatePasswordUser)
func (h *Profile) UpdatePasswordUser(c fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok {
		return fiber.NewError(
			fiber.StatusUnauthorized,
			"invalid user ID",
		)
	}
	var req request.ChangePasswordRequest
	if err := c.Bind().Body(&req); err != nil {
		return err
	}
	// Jalankan validasi struct
	if err := v.Validate.Struct(req); err != nil {
		return err
	}
	res, err := h.service.UpdatePasswordUser(c.Context(), userID, req)
	if err != nil {
		return err
	}
	return c.Status(200).JSON(response.CommonResponse{
		Message: "success update password",
		Data: fiber.Map{
			"user_id":    res.UserID,
			"name":       res.User.Name,
			"updated_at": res.UpdatedAt,
		},
	})
}
