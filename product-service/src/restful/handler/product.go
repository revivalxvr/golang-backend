package handler

import (
	"product-service/src/common/dto/request"
	"product-service/src/common/dto/response"
	"product-service/src/factory"
	"product-service/src/restful/client"
	"product-service/src/restful/service"

	"product-service/src/common/log"

	"github.com/gofiber/fiber/v3"
)

type Product struct {
	service service.Product
	ik      *client.ImageKit
}

func NewProductHandler(f *factory.Factory) *Product {
	return &Product{
		service: service.NewProduct(f),
		ik:      f.ImageKitClient,
	}
}

func (h *Product) Create(c fiber.Ctx) error {
	req := new(request.CreateProduct)

	if err := c.Bind().Body(req); err != nil {
		return err
	}

	fileID, err := h.service.Create(c, req)
	if err != nil {
		if fileID != "" {
			if err := h.ik.DeleteFile(c.Context(), fileID); err != nil {
				log.Logger.Error(err.Error())
			}
		}
		return err
	}

	return c.Status(201).JSON(response.CommonResponse{
		Message: "berhasil buat produk",
		Data:    fileID,
	})
}
func (h *Product) Update(c fiber.Ctx) error {
	productId := c.Params("id")
	req := new(request.UpdateProduct)
	if err := c.Bind().Body(req); err != nil {
		return err
	}

	// Panggil service update
	if err := h.service.Update(c, productId, req); err != nil {
		return err
	}

	return c.Status(200).JSON(response.CommonResponse{
		Message: "berhasil update produk",
		Data:    nil,
	})
}
func (h *Product) Delete(c fiber.Ctx) error {
	productId := c.Params("id")

	// Panggil service delete
	if err := h.service.Delete(c, productId); err != nil {
		return err
	}

	return c.Status(200).JSON(response.CommonResponse{
		Message: "berhasil hapus produk",
		Data:    nil,
	})
}

func (h *Product) FindMany(c fiber.Ctx) error {
	req := new(request.FindManyProduct)

	if err := c.Bind().Query(req); err != nil {
		return err
	}

	res, err := h.service.FindMany(c, req)
	if err != nil {
		return err
	}

	return c.Status(200).JSON(response.CommonResponse{
		Message: "Get Product Success",
		Data:    res,
	})

}
