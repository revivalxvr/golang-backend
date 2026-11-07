package service

import (
	"bag-service/src/common/dto/request"
	"bag-service/src/common/model"
	v "bag-service/src/common/pkg/validator"
	"bag-service/src/factory"
	grpcclinet "bag-service/src/grpc/client"
	"bag-service/src/repository"

	"github.com/gofiber/fiber/v3"

	pb "github.com/revivalxvr/protobuf/protogen/cart"
)

type CartItemResponse struct {
	ProductID   string  `json:"product_id"`
	ProductName string  `json:"product_name"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
	SubTotal    float64 `json:"sub_total"`
	Stock       int32   `json:"stock"`
	Image       string  `json:"image"`
}

type CartService interface {
	AddCart(c fiber.Ctx, req *request.AddCart) (*model.Cart, error)
	GetCart(c fiber.Ctx, userID string) ([]CartItemResponse, error)
	DeleteCartItem(c fiber.Ctx, userID string, productID string) error
}
type cartImpl struct {
	cartRepository repository.CartRepository
	cartClinet     grpcclinet.Cart
}

func NewCartService(f *factory.Factory) CartService {
	return &cartImpl{
		cartRepository: f.Cartrepository,
		cartClinet:     f.CartClient,
	}
}

func (s *cartImpl) AddCart(c fiber.Ctx, req *request.AddCart) (*model.Cart, error) {
	ctx := c.Context()

	// 1. Validasi input request
	if err := v.Validate.Struct(req); err != nil {
		return nil, err
	}

	// 2. Mapping dari request ke model.Cart (representasi database)
	cart := &model.Cart{
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
		UserID:    req.UserID,
	}

	// 3. Panggil repository (hanya mengembalikan 1 nilai yaitu error)
	err := s.cartRepository.AddCart(ctx, cart)
	if err != nil {
		return nil, err
	}

	// 4. Kembalikan cart (karena GORM otomatis mengisi data/ID setelah Create sukses) dan nil untuk error
	return cart, nil
}
func (s *cartImpl) GetCart(c fiber.Ctx, userID string) ([]CartItemResponse, error) {
	ctx := c.Context()

	// 1. Ambil data keranjang lokal milik user dari database
	localCarts, err := s.cartRepository.GetCartItemsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if len(localCarts) == 0 {
		return []CartItemResponse{}, nil
	}

	// 2. Kumpulkan semua product_id untuk request gRPC secara massal (bulk)
	var productIDs []string
	qtyMap := make(map[string]int)
	for _, item := range localCarts {
		productIDStr := item.ProductID // <-- Tambahkan .String() di sini jika menggunakan UUID
		productIDs = append(productIDs, productIDStr)
		qtyMap[productIDStr] = item.Quantity
	}

	// 3. Panggil gRPC client untuk mendapatkan harga dan stok terbaru dari product-service
	grpcReq := &pb.GetProductsForCartReq{
		ProductIds: productIDs,
	}
	grpcRes, err := s.cartClinet.GetProductsForCart(ctx, grpcReq)
	if err != nil {
		return nil, err
	}

	// Mapping response gRPC ke dalam map agar mudah di-lookup
	productMap := make(map[string]*pb.ProductCartInfo)
	for _, p := range grpcRes.Products {
		productMap[p.ProductId] = p
	}

	// 4. Gabungkan data lokal (quantity) dengan data gRPC (nama, harga, stok, image) -> Data Enrichment
	var response []CartItemResponse
	for _, id := range productIDs {
		qty := qtyMap[id]
		prodInfo, exists := productMap[id]

		var name string
		var price float64
		var stock int32
		var image string

		if exists {
			name = prodInfo.Name
			price = prodInfo.Price
			stock = prodInfo.Stock
			image = prodInfo.Image
		}

		response = append(response, CartItemResponse{
			ProductID:   id,
			ProductName: name,
			Price:       price,
			Quantity:    qty,
			SubTotal:    price * float64(qty),
			Stock:       stock,
			Image:       image,
		})
	}

	return response, nil
}

// DeleteCartItem
func (s *cartImpl) DeleteCartItem(c fiber.Ctx, userID string, productID string) error {
	// Hapus item berdasarkan user_id dan product_id agar aman (milik user yang bersangkutan saja)
	return s.cartRepository.DeleteByUserIDAndProductID(c, userID, productID)
}
