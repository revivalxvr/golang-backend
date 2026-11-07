package repository

import (
	"context"
	"order-service/src/common/model"

	"gorm.io/gorm"
)

type ProductOrder interface {
	Create(ctx context.Context, product []*model.ProductOrder) error
	Find(ctx context.Context, oderId string) ([]*model.ProductOrder, error)
}

type productOrderImpl struct {
	db *gorm.DB
}

func NewProductOrder(db *gorm.DB) ProductOrder {
	return &productOrderImpl{
		db: db,
	}
}

func (r *productOrderImpl) Create(ctx context.Context, product []*model.ProductOrder) error {
	err := r.db.WithContext(ctx).Create(product).Error
	return err
}

func (r *productOrderImpl) Find(ctx context.Context, oderId string) ([]*model.ProductOrder, error) {
	var product []*model.ProductOrder
	err := r.db.WithContext(ctx).Where("order_id = ?", oderId).Find(&product).Error
	return product, err
}
