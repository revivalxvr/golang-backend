package handler

import (
	"product-service/src/common/dto/request"
	"product-service/src/common/dto/response"
	"product-service/src/factory"
	"product-service/src/restful/service"

	"github.com/gofiber/fiber/v3"
)

type Category struct {
	service service.Category
}

func NewCategoryHandler(f *factory.Factory) *Category {
	return &Category{
		service: service.NewCategoryService(f),
	}
}
func (h *Category) CreateCategory(c fiber.Ctx) error {
	req := new(request.CreateCategory)

	if err := c.Bind().Body(req); err != nil {
		return err
	}

	// Memanggil service dan mendapatkan objek kategori yang baru dibuat
	categoryData, err := h.service.CreateCategory(c, req)
	if err != nil {
		return err
	}

	// Mengembalikan response dengan data kategori terisi
	return c.Status(201).JSON(response.CommonResponse{
		Message: "Category created successfully",
		Data:    categoryData, // <-- Berisi data objek kategori (ID, Name, CreatedAt, dll)
	})
}
func (h *Category) FindCategory(c fiber.Ctx) error {
	req := new(request.FindManyCategory)

	if err := c.Bind().Query(req); err != nil {
		return err
	}
	res, err := h.service.FindCategory(c, req)
	if err != nil {
		return err
	}
	return c.Status(200).JSON(response.CommonResponse{
		Message: "Get Category Success",
		Data:    res,
	})
}

func (h *Category) UpdateCategory(c fiber.Ctx) error {
	categoryId := c.Params("id")
	req := new(request.UpdateCategory)

	if err := c.Bind().Body(req); err != nil {
		return err
	}

	// Panggil service update
	if err := h.service.UpdateCategory(c, categoryId, req); err != nil {
		return err
	}

	return c.Status(200).JSON(response.CommonResponse{
		Message: "Category updated successfully",
		Data:    nil,
	})
}

// DeleteCategory
func (h *Category) DeleteCategory(c fiber.Ctx) error {
	categoryId := c.Params("id")
	// Panggil service delete
	if err := h.service.DeleteCategory(c, categoryId); err != nil {
		return err
	}

	return c.Status(200).JSON(response.CommonResponse{
		Message: "Category deleted successfully",
		Data:    nil,
	})
}
