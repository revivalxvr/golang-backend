package repository

import (
	"context"
	"order-service/src/common/model"

	"gorm.io/gorm"
)

type Order interface {
	Create(ctx context.Context, order *model.Order) error
	Update(ctx context.Context, whereClause string, args []any, data map[string]any) error
	Find(ctx context.Context, whereClause string, args []any) (*model.Order, error)
	FindByUserID(ctx context.Context, userID string) ([]model.OrderModel, error)
}

type orderImpl struct {
	db *gorm.DB
}

func NewOrder(db *gorm.DB) Order {
	return &orderImpl{
		db: db,
	}
}

func (r *orderImpl) Create(ctx context.Context, order *model.Order) error {
	return r.db.WithContext(ctx).Create(order).Error
}

func (r *orderImpl) Update(ctx context.Context, whereClause string, args []any, data map[string]any) error {
	return r.db.WithContext(ctx).Model(model.Order{}).Where(whereClause, args...).Updates(data).Error
}
func (r *orderImpl) Find(ctx context.Context, whereClause string, args []any) (*model.Order, error) {
	var order model.Order
	if err := r.db.WithContext(ctx).Where(whereClause, args...).First(&order).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &order, nil
}
func (r *orderImpl) FindByUserID(ctx context.Context, userID string) ([]model.OrderModel, error) {
	var orders []model.OrderModel

	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Preload("ProductOrders").
		Order("created_at DESC").
		Find(&orders).Error

	if err != nil {
		return nil, err
	}

	return orders, nil
}
