package service

import (
	"log"
	"product-service/src/common/dto/request"
	"product-service/src/common/errors"
	"product-service/src/common/model"
	v "product-service/src/common/pkg/validator"
	"product-service/src/common/utils"
	"product-service/src/factory"
	"product-service/src/repository"
	"product-service/src/restful/client"

	"github.com/google/uuid"

	"github.com/gofiber/fiber/v3"
)

type Product interface {
	Create(c fiber.Ctx, req *request.CreateProduct) (string, error)
	FindMany(c fiber.Ctx, req *request.FindManyProduct) ([]*model.Product, error)
	Update(c fiber.Ctx, id string, req *request.UpdateProduct) error
	Delete(c fiber.Ctx, id string) error
}

type productImpl struct {
	productRepository repository.ProductRepository
	imagekitclinet    *client.ImageKit
}

func NewProduct(f *factory.Factory) Product {
	return &productImpl{
		productRepository: f.ProductRepository,
		imagekitclinet:    f.ImageKitClient,
	}
}

func (s *productImpl) Create(c fiber.Ctx, req *request.CreateProduct) (string, error) {
	fileName, ok := c.Locals("fileName").(string)
	if !ok || fileName == "" {
		return "", &errors.Response{HttpCode: 400, Message: "Invalid Image"}
	}

	path := "./tmp/" + fileName
	res, err := s.imagekitclinet.UploadFile(c.Context(), path, fileName)
	if err != nil {
		return "", err
	}

	utils.DeleteFile(path)

	req.ImageId = res.FileID
	req.Image = res.URL

	if err := v.Validate.Struct(req); err != nil {
		return "", err
	}
	if err := s.productRepository.Create(c.Context(), &model.Product{
		ID:         uuid.New(),
		Name:       req.Name,
		Price:      req.Price,
		Image:      req.Image,
		ImageId:    req.ImageId,
		Sku:        req.Sku,
		Stock:      req.Stock,
		CategoryID: req.CategoryID,
		Badge:      req.Badge,
		Length:     req.Length,
		Width:      req.Width,
		Height:     req.Height,
		Deskripsi:  req.Deskripsi,
	}); err != nil {
		return "", err
	}
	return res.FileID, nil
}
func (s *productImpl) Update(c fiber.Ctx, id string, req *request.UpdateProduct) error {
	// 1. Cari produk lama berdasarkan ID (pastikan produknya ada)
	// Sesuaikan pemanggilan repository Anda untuk FindByID
	var product model.Product
	if err := s.productRepository.FindByID(c.Context(), id, &product); err != nil {
		return &errors.Response{HttpCode: 404, Message: "Product not found"}
	}

	// 2. Cek apakah ada file gambar BARU yang diunggah melalui middleware temporary image
	fileName, hasNewImage := c.Locals("fileName").(string)

	var newFileID, newImageURL string
	if hasNewImage && fileName != "" {
		path := "./tmp/" + fileName

		// Upload gambar baru ke ImageKit
		res, err := s.imagekitclinet.UploadFile(c.Context(), path, fileName)
		if err != nil {
			return err
		}

		// Bersihkan file temporary di lokal
		utils.DeleteFile(path)

		newFileID = res.FileID
		newImageURL = res.URL

		// 3. Hapus gambar LAMA dari ImageKit jika produk lama memiliki ImageId
		if product.ImageId != "" {
			_ = s.imagekitclinet.DeleteFile(c.Context(), product.ImageId)
			// Error di-ignore atau di-log saja agar proses update database tidak terhenti total
		}

		// Update properti gambar di struct produk
		product.ImageId = newFileID
		product.Image = newImageURL
	}

	// 4. Update data teks lainnya jika dikirim
	if req.Name != "" {
		product.Name = req.Name
	}
	if req.Price > 0 {
		product.Price = req.Price
	}
	if req.Sku != "" {
		product.Sku = req.Sku
	}
	if req.Stock > 0 {
		product.Stock = req.Stock
	}
	if req.CategoryID > 0 {
		product.CategoryID = req.CategoryID
	}
	if req.Deskripsi != "" {
		product.Deskripsi = req.Deskripsi
	}
	if req.Badge != "" {
		product.Badge = req.Badge
	}
	if req.Length > 0 {
		product.Length = req.Length
	}
	if req.Width > 0 {
		product.Width = req.Width
	}
	if req.Height > 0 {
		product.Height = req.Height
	}

	// Validasi data produk yang sudah diperbarui
	if err := v.Validate.Struct(&product); err != nil {
		return err
	}

	// 5. Simpan perubahan ke repository/database
	if err := s.productRepository.Update(c.Context(), &product); err != nil {
		// (Opsional) Jika update DB gagal dan sempat upload gambar baru, Anda bisa hapus newFileID di sini
		return err
	}

	return nil
}

func (s *productImpl) FindMany(c fiber.Ctx, req *request.FindManyProduct) ([]*model.Product, error) {
	if err := v.Validate.Struct(req); err != nil {
		return nil, err
	}
	res, err := s.productRepository.FindMany(c.Context(), req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
func (s *productImpl) Delete(c fiber.Ctx, id string) error {
	// 1. Cari produk berdasarkan ID (pastikan produknya ada)
	var product model.Product
	if err := s.productRepository.FindByID(c.Context(), id, &product); err != nil {
		return &errors.Response{HttpCode: 404, Message: "Product not found"}
	}

	// 2. Hapus data produk dari database
	if err := s.productRepository.Delete(c.Context(), id); err != nil {
		return err
	}

	// 3. Hapus gambar dari ImageKit jika ImageId tersedia
	if product.ImageId != "" {
		if err := s.imagekitclinet.DeleteFile(c.Context(), product.ImageId); err != nil {
			// Error dari ImageKit bisa di-log saja agar tidak menggagalkan response sukses ke client
			log.Println(err)
		}
	}

	return nil
}
