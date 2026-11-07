package service

import (
	"product-service/src/common/dto/request"
	"product-service/src/common/errors"
	"product-service/src/common/model"
	v "product-service/src/common/pkg/validator"
	"product-service/src/factory"
	"product-service/src/repository"

	"github.com/gofiber/fiber/v3"
)

type Category interface {
	CreateCategory(c fiber.Ctx, req *request.CreateCategory) (*model.Category, error)
	FindCategory(c fiber.Ctx, req *request.FindManyCategory) ([]*model.Category, error)
	UpdateCategory(c fiber.Ctx, id string, req *request.UpdateCategory) error
	DeleteCategory(c fiber.Ctx, id string) error
}

type categoryImpl struct {
	categoryRepository repository.CategoryRepository
}

func NewCategoryService(f *factory.Factory) Category {
	return &categoryImpl{
		categoryRepository: f.CategoryRepository,
	}
}

func (s *categoryImpl) CreateCategory(c fiber.Ctx, req *request.CreateCategory) (*model.Category, error) {
	ctx := c.Context()
	if err := v.Validate.Struct(req); err != nil {
		return nil, err
	}

	// Inisialisasi model kategori
	category := &model.Category{
		Name:      req.Name,
		Slug:      req.Slug,
		Deskripsi: req.Deskripsi,
		IsActive:  req.IsActive,
	}

	// Panggil repository untuk create data ke database
	if err := s.categoryRepository.CreateCategory(ctx, category); err != nil {
		return nil, err
	}

	// GORM otomatis mengisi category.ID setelah berhasil create, kembalikan objeknya
	return category, nil
}
func (s *categoryImpl) FindCategory(c fiber.Ctx, req *request.FindManyCategory) ([]*model.Category, error) {
	if err := v.Validate.Struct(req); err != nil {
		return nil, err
	}
	res, err := s.categoryRepository.FindManyCategory(c.Context(), req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
func (s *categoryImpl) UpdateCategory(c fiber.Ctx, id string, req *request.UpdateCategory) error {
	var category model.Category
	if err := s.categoryRepository.FindCategoryByID(c.Context(), id, &category); err != nil {
		return &errors.Response{HttpCode: 404, Message: "Category not found"}
	}

	if req.Name != "" {
		category.Name = req.Name
	}
	if req.Slug != "" {
		category.Slug = req.Slug
	}
	if req.Deskripsi != "" {
		category.Deskripsi = req.Deskripsi
	}
	// 3. Update IsActive hanya jika dikirim oleh frontend (tidak nil)
	if req.IsActive != nil {
		category.IsActive = *req.IsActive // Dereference pointer untuk mengambil nilai aslinya (true/false)
	}
	// 4. Simpan ke database
	if err := s.categoryRepository.UpdateCategory(c.Context(), &category); err != nil {
		return err
	}

	return nil

}
func (s *categoryImpl) DeleteCategory(c fiber.Ctx, id string) error {
	var category model.Category
	if err := s.categoryRepository.FindCategoryByID(c.Context(), id, &category); err != nil {
		return &errors.Response{HttpCode: 404, Message: "Category not found"}
	}
	if err := s.categoryRepository.DeleteCategory(c.Context(), &category); err != nil {
		return err
	}
	return nil
}
