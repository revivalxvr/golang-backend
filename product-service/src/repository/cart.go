package repository

import (
	"context"
	"product-service/src/common/model"

	"gorm.io/gorm"
)

type CartRepository interface {
	GetProductsByIDs(ctx context.Context, productIDs []string) ([]model.Product, error)
}

type cartImpl struct {
	db *gorm.DB
}

func NewCart(db *gorm.DB) CartRepository {
	return &cartImpl{
		db: db,
	}
}

func (r *cartImpl) GetProductsByIDs(ctx context.Context, productIDs []string) ([]model.Product, error) {
	var products []model.Product

	// Melakukan query massal dengan WHERE id IN (...)
	err := r.db.WithContext(ctx).Where("id IN ?", productIDs).Find(&products).Error
	if err != nil {
		return nil, err
	}

	return products, nil
}
