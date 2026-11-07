package repository

import (
	"gorm.io/gorm"
)

type Transactioner interface {
	Commit() error
	Rollback() error

	OrderRepository() Order
	ProductOrderRepository() ProductOrder
}

type transactioner struct {
	db *gorm.DB
}

func (t *transactioner) Commit() error {
	return t.db.Commit().Error
}

func (t *transactioner) Rollback() error {
	return t.db.Rollback().Error
}

func (t *transactioner) OrderRepository() Order {
	return NewOrder(t.db)
}

func (t *transactioner) ProductOrderRepository() ProductOrder {
	return NewProductOrder(t.db)
}
