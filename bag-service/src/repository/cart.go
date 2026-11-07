package repository

import (
	"bag-service/src/common/model"
	"context"

	"github.com/google/uuid"

	"gorm.io/gorm"
)

type CartRepository interface {
	AddCart(ctx context.Context, cart *model.Cart) error
	GetCartItemsByUserID(ctx context.Context, userID string) ([]model.Cart, error)
	DeleteCheckoutItems(ctx context.Context, userID uuid.UUID, productIDs []uuid.UUID) error
	DeleteByUserIDAndProductID(ctx context.Context, userID string, productID string) error
}

type cartImpl struct {
	db *gorm.DB
}

func NewCart(db *gorm.DB) CartRepository {
	return &cartImpl{
		db: db,
	}
}

func (s *cartImpl) AddCart(ctx context.Context, cart *model.Cart) error {
	var existingCart model.Cart
	err := s.db.WithContext(ctx).
		Where("user_id = ? AND product_id = ?", cart.UserID, cart.ProductID).
		First(&existingCart).Error

	if err == nil {
		// Jika sudah ada, tambahkan quantity-nya
		existingCart.Quantity += cart.Quantity
		return s.db.WithContext(ctx).Save(&existingCart).Error
	} else {
		// Jika belum ada, buat record baru
		return s.db.WithContext(ctx).Create(cart).Error
	}
}
func (r *cartImpl) GetCartItemsByUserID(ctx context.Context, userID string) ([]model.Cart, error) {
	var carts []model.Cart
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&carts).Error
	return carts, err
}

// function delete untuk checkout (order-service)
func (r *cartImpl) DeleteCheckoutItems(ctx context.Context, userID uuid.UUID, productIDs []uuid.UUID) error {
	if len(productIDs) == 0 {
		return nil
	}
	err := r.db.WithContext(ctx).Where(
		"user_id = ? AND product_id IN ?",
		userID,
		productIDs,
	).
		Delete(&model.Cart{}).
		Error
	return err
}
func (r *cartImpl) DeleteByUserIDAndProductID(ctx context.Context, userID string, productID string) error {
	// Eksekusi query DELETE WHERE user_id = ? AND product_id = ?
	return r.db.WithContext(ctx).
		Where("user_id = ? AND product_id = ?", userID, productID).
		Delete(&model.Cart{}).Error
}
