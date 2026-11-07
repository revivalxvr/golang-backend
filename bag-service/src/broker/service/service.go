package service

import (
	"bag-service/src/common/dto/request"
	v "bag-service/src/common/pkg/validator"
	"bag-service/src/factory"
	"bag-service/src/repository"
	"context"
	"encoding/json"
	"fmt"
)

type RabbitMQ struct {
	cartRepository repository.CartRepository
}

func NewRabbitMQ(f *factory.Factory) *RabbitMQ {
	return &RabbitMQ{
		cartRepository: f.Cartrepository,
	}
}

func (s *RabbitMQ) RemoveCheckedOutItems(ctx context.Context, data []byte) error {
	var req request.OrderCreated
	if err := json.Unmarshal(data, &req); err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	// 2. Validasi payload.
	if err := v.Validate.Struct(req); err != nil {
		return fmt.Errorf(
			"invalid OrderCreated event: %w",
			err,
		)
	}
	// 3. Hapus item yang termasuk dalam order.
	if err := s.cartRepository.DeleteCheckoutItems(
		ctx,
		req.UserID,
		req.ProductIDs,
	); err != nil {
		return fmt.Errorf(
			"failed to delete checkout cart items for order %s: %w",
			req.OrderID,
			err,
		)
	}
	return nil
}
