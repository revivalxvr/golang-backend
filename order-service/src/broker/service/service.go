package service

import (
	"context"
	"encoding/json"
	"fmt"
	"order-service/src/common/dto/request"
	v "order-service/src/common/pkg/validator"
	"order-service/src/factory"
	"order-service/src/grpc/client"
	"order-service/src/repository"

	pb "github.com/revivalxvr/protobuf/protogen/product"
)

type RabbitMQ struct {
	orderRepository        repository.Order
	productOrderRepository repository.ProductOrder
	productClient          client.Product
}

func NewRabbitMQ(f *factory.Factory) *RabbitMQ {
	return &RabbitMQ{
		orderRepository:        f.OrderRepository,
		productOrderRepository: f.ProductOrderRepository,
		productClient:          f.ProductClient,
	}
}

func (s *RabbitMQ) HandlerMidtransTxNotif(ctx context.Context, data []byte) error {
	var req request.MidtransTxNotif
	if err := json.Unmarshal(data, &req); err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	// 1. Cek keberadaan Order di Database menggunakan Find
	existingOrder, err := s.orderRepository.Find(ctx, "id = ?", []any{req.OrderId})
	if err != nil {
		return fmt.Errorf("failed to fetch order: %w", err)
	}
	if existingOrder == nil {
		return fmt.Errorf("order with ID %s not found", req.OrderId)
	}

	// 2. Cek Idempotency: Jika status sudah PAID, CANCELLED, atau FAILED, abaikan proses ulang
	if existingOrder.Status == "PAID" || existingOrder.Status == "CANCELLED" || existingOrder.Status == "FAILED" {
		// Return nil agar pesan di-ACK oleh RabbitMQ karena sudah pernah diproses sampai selesai
		return nil
	}

	// Evaluasi transaksi sukses
	isSuccess := (req.TransactionStatus == "settlement" || req.TransactionStatus == "capture") &&
		(req.FraudStatus == "" || req.FraudStatus == "accept")
	if isSuccess {
		if err := v.Validate.Struct(&req); err != nil {
			return fmt.Errorf("validation error: %w", err)
		}

		err := s.orderRepository.Update(ctx, "id = ?", []any{req.OrderId}, map[string]any{
			"status":         "PAID",
			"payment_method": req.PaymentType,
		})
		if err != nil {
			return fmt.Errorf("failed to update order status to PAID: %w", err)
		}
		return nil
	}

	// Evaluasi transaksi pending
	if req.TransactionStatus == "pending" {
		if existingOrder.Status == "PENDING" {
			return nil // Tidak perlu update DB jika memang sudah PENDING
		}

		if err := s.orderRepository.Update(ctx, "id = ?", []any{req.OrderId}, map[string]any{
			"status": "PENDING",
		}); err != nil {
			return fmt.Errorf("failed to update order status to PENDING: %w", err)
		}
		return nil
	}

	// Evaluasi transaksi gagal atau batal (butuh rollback stok)
	switch req.TransactionStatus {
	case "cancel", "deny", "expire":
		newStatus := "FAILED"
		if req.TransactionStatus == "cancel" {
			newStatus = "CANCELLED"
		}

		if err := s.rollbackStockAndFailOrder(ctx, req.OrderId, newStatus); err != nil {
			return fmt.Errorf("failed to handle %s transaction: %w", req.TransactionStatus, err)
		}
	}

	return nil
}

// Helper internal untuk menyatukan logika rollback stok dan update status
func (s *RabbitMQ) rollbackStockAndFailOrder(ctx context.Context, orderID string, targetStatus string) error {
	productOrders, err := s.productOrderRepository.Find(ctx, orderID)
	if err != nil {
		return fmt.Errorf("failed to find product order items: %w", err)
	}

	if len(productOrders) > 0 {
		reqItems := make([]*pb.ProductOrder, 0, len(productOrders))
		for _, product := range productOrders {
			reqItems = append(reqItems, &pb.ProductOrder{
				ProductId: product.ProductID,
				Quantity:  uint32(product.Quantity),
			})
		}

		if err := s.productClient.RollbackStock(ctx, &pb.RollbackStocksReq{ProductOrders: reqItems}); err != nil {
			return fmt.Errorf("failed to rollback stock via gRPC: %w", err)
		}
	}

	if err := s.orderRepository.Update(ctx, "id = ?", []any{orderID}, map[string]any{"status": targetStatus}); err != nil {
		return fmt.Errorf("failed to update order status to %s: %w", targetStatus, err)
	}

	return nil
}
