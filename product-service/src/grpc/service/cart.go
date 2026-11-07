package service

import (
	"context"
	"product-service/src/factory"
	"product-service/src/repository"

	pb "github.com/revivalxvr/protobuf/protogen/cart"
)

type CartService interface {
	GetProductsForCart(ctx context.Context, req *pb.GetProductsForCartReq) (*pb.GetProductsForCartRes, error)
}
type cartImpl struct {
	cartrepository repository.CartRepository
}

func NewCartService(f *factory.Factory) CartService {
	return &cartImpl{
		cartrepository: f.Cartrepository,
	}
}

func (s *cartImpl) GetProductsForCart(ctx context.Context, req *pb.GetProductsForCartReq) (*pb.GetProductsForCartRes, error) {
	// 1. Panggil repository untuk mendapatkan data produk berdasarkan slice product_ids
	products, err := s.cartrepository.GetProductsByIDs(ctx, req.ProductIds)
	if err != nil {
		return nil, err
	}

	// 2. Mapping dari model database ke format protobuf response (`pb.ProductCartInfo`)
	var pbProducts []*pb.ProductCartInfo
	for _, p := range products {
		pbProducts = append(pbProducts, &pb.ProductCartInfo{
			ProductId:   p.ID.String(),  // Ubah ke .String() jika ID produk bertipe UUID
			Name:        p.Name,         // Sesuaikan dengan nama field di model Product Anda
			Price:       p.Price,        // Sesuaikan dengan field harga
			Stock:       int32(p.Stock), // Sesuaikan dengan field stok
			IsAvailable: p.Stock > 0,
			Image:       p.Image,
		})
	}

	// 3. Kembalikan response sesuai dengan gRPC struct
	return &pb.GetProductsForCartRes{
		Products: pbProducts,
	}, nil
}
