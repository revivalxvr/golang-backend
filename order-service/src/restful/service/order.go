package service

import (
	"context"
	"order-service/src/common/constant"
	"order-service/src/common/dto/request"
	"order-service/src/common/dto/response"
	v "order-service/src/common/pkg/validator"
	"order-service/src/factory"
	grpcclinet "order-service/src/grpc/client"
	"order-service/src/repository"
	"order-service/src/restful/client"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	pbi "github.com/revivalxvr/protobuf/protogen/cart"
	pb "github.com/revivalxvr/protobuf/protogen/product"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

type Order interface {
	CreateOrder(c fiber.Ctx, order *request.CreateOrder, userID string) (*response.MidtransTx, error)
	GetOrderByUser(c fiber.Ctx, userID string) ([]response.OrderWithProducts, error)
}

type orderImpl struct {
	midtransClient  client.Midtrans
	orderRepository repository.Order
	productClient   grpcclinet.Product
	cartClient      grpcclinet.Cart
	txRepository    repository.TxBeginner
}

func NewOrder(f *factory.Factory) Order {
	return &orderImpl{
		midtransClient:  f.MidtransClient,
		orderRepository: f.OrderRepository,
		productClient:   f.ProductClient,
		cartClient:      f.CartClient,
		txRepository:    f.TxRepository,
	}
}
func (s *orderImpl) CreateOrder(c fiber.Ctx, order *request.CreateOrder, userID string) (*response.MidtransTx, error) {
	ctx := c.Context()
	if err := v.Validate.Struct(order); err != nil {
		return nil, err
	}

	orderID, err := gonanoid.New()
	if err != nil {
		return nil, err
	}
	order.Order.ID = orderID
	order.Order.UserID = userID
	order.Order.Status = "PENDING" // Set default status order

	// Siapkan item details untuk Midtrans Snap
	var itemDetails []request.SnapItemDetail
	for _, p := range order.ProductOrders {
		itemDetails = append(itemDetails, request.SnapItemDetail{
			ID:    p.ProductID,
			Price: int64(p.Price),
			Qty:   p.Quantity,
			Name:  "Product " + p.ProductID, // Bisa disesuaikan jika nama product didapat
		})
	}

	midtransTxReq := &request.SnapTransactionRequest{
		TransactionDetails: struct {
			OrderId     string `json:"order_id"`
			GrossAmount int64  `json:"gross_amount"`
		}{
			OrderId:     order.Order.ID,
			GrossAmount: int64(order.Order.GrossAmount),
		},
		ItemDetails: itemDetails,
		// Kirim metadata agar Midtrans mengembalikannya pada webhook notification
		Metadata: map[string]interface{}{
			"original_order_id": order.Order.ID,
		},
	}

	// Panggil Midtrans Snap API
	res, err := s.midtransClient.CreateTransaction(ctx, midtransTxReq)
	if err != nil {
		return nil, err
	}

	// Simpan Payment URL / Token ke model order jika diperlukan
	order.Order.PaymentUrl = res.RedirectURL

	// Simpan order & product_orders ke database via transaction
	tx, err := s.txRepository.Begin()
	if err != nil {
		return nil, err
	}

	for _, p := range order.ProductOrders {
		p.ID = uuid.New()
		p.OrderID = order.Order.ID
	}

	if err := tx.OrderRepository().Create(ctx, order.Order); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.ProductOrderRepository().Create(ctx, order.ProductOrders); err != nil {
		tx.Rollback()
		return nil, err
	}

	var req []*pb.ProductOrder
	for _, v := range order.ProductOrders {
		req = append(req, &pb.ProductOrder{
			ProductId: v.ProductID,
			Quantity:  uint32(v.Quantity),
		})
	}

	ctx = context.WithValue(ctx, constant.RequestID, c.Locals("request_id"))
	ctx = context.WithValue(ctx, constant.UserID, c.Locals("user_id"))

	err = s.productClient.ReduceStocks(ctx, &pb.ReduceStocksReq{ProductOrders: req})
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return nil, err
	}

	return res, nil
}

// GetOrderByUser(c, parsedUserID)
func (s *orderImpl) GetOrderByUser(c fiber.Ctx, userID string) ([]response.OrderWithProducts, error) {
	ctx := c.Context()
	// 1. Ambil orders milik user beserta product_orders dari database lokal
	orders, err := s.orderRepository.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	var result []response.OrderWithProducts

	for _, order := range orders {
		var productIds []string
		// Kumpulkan semua product_id yang ada di order ini untuk di-batch gRPC request
		for _, item := range order.ProductOrders {
			productIds = append(productIds, item.ProductID)
		}

		// 2. Panggil gRPC Product Service untuk mengambil detail (nama, gambar, dll) secara massal
		productMap := make(map[string]*pbi.ProductCartInfo)
		if len(productIds) > 0 {
			grpcRes, err := s.cartClient.GetProductsForCart(ctx, &pbi.GetProductsForCartReq{
				ProductIds: productIds,
			})
			if err == nil && grpcRes != nil {
				for _, p := range grpcRes.Products {
					productMap[p.ProductId] = p
				}
			}
		}

		// 3. Mapping data ke struct response
		var itemsRes []response.ProductItemResponse
		for _, item := range order.ProductOrders {
			name := "Unknown Product"
			image := ""

			if prodInfo, exists := productMap[item.ProductID]; exists {
				name = prodInfo.Name
				image = prodInfo.Image
			}

			itemsRes = append(itemsRes, response.ProductItemResponse{
				ID:           item.ID,
				ProductID:    item.ProductID,
				ProductName:  name,
				ProductImage: image,
				Quantity:     item.Quantity,
				Price:        item.Price,
				CreatedAt:    item.CreatedAt,
			})
		}

		result = append(result, response.OrderWithProducts{
			ID:            order.ID,
			UserID:        order.UserID,
			GrossAmount:   order.GrossAmount,
			Status:        order.Status,
			PaymentMethod: order.PaymentMethod,
			PaymentUrl:    order.PaymentUrl,
			CreatedAt:     order.CreatedAt,
			UpdatedAt:     order.UpdatedAt,
			Items:         itemsRes,
		})
	}

	return result, nil
}
