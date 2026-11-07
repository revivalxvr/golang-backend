package handler

import (
	"bag-service/src/common/dto/request"
	"bag-service/src/common/dto/response"
	"bag-service/src/factory"
	service "bag-service/src/restful/service"

	"github.com/gofiber/fiber/v3"
)

type Cart struct {
	service service.CartService
}

func NewAddCart(f *factory.Factory) *Cart {
	return &Cart{
		service: service.NewCartService(f),
	}
}

func (h *Cart) AddCart(c fiber.Ctx) error {
	userID := c.Locals("user_id").(string) //ambil dari context Fiber diset oleh middleware autentikasi (seperti JWT)
	req := new(request.AddCart)

	if err := c.Bind().Body(req); err != nil {
		return err
	}

	req.UserID = userID

	res, err := h.service.AddCart(c, req)
	if err != nil {
		return err
	}
	return c.Status(200).JSON(response.CommonResponse{
		Message: "success add cart",
		Data:    res,
	})
}
func (h *Cart) GetCart(c fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	res, err := h.service.GetCart(c, userID)
	if err != nil {
		return err
	}
	return c.Status(200).JSON(response.CommonResponse{
		Message: "success get cart",
		Data:    res,
	})
}
func (h *Cart) DeleteCartItem(c fiber.Ctx) error {
	// Ambil user_id dari middleware JWT / Locals
	userID := c.Locals("user_id").(string)

	// Ambil product_id dari parameter URL (contoh: /api/v1/cart/21f20ca3-...)
	productID := c.Params("product_id")

	if productID == "" {
		return c.Status(400).JSON(fiber.Map{
			"message": "product_id is required",
		})
	}

	// Panggil service untuk menghapus data di database
	err := h.service.DeleteCartItem(c, userID, productID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "failed to delete cart item",
			"error":   err.Error(),
		})
	}

	return c.Status(200).JSON(response.CommonResponse{
		Message: "success delete item from cart",
		Data:    nil,
	})
}
