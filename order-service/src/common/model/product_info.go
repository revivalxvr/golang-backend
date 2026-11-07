package model

import (
	"time"
)

type OrderModel struct {
	ID            string              `gorm:"primaryKey;type:varchar(36);column:id"`
	UserID        string              `gorm:"type:varchar(36);column:user_id"`
	GrossAmount   int                 `gorm:"column:gross_amount"`
	Status        string              `gorm:"type:varchar(255);column:status"`
	PaymentMethod string              `gorm:"type:varchar(255);column:payment_method"`
	PaymentUrl    string              `gorm:"type:varchar(500);column:payment_url"`
	CreatedAt     time.Time           `gorm:"column:created_at"`
	UpdatedAt     time.Time           `gorm:"column:updated_at"`
	ProductOrders []ProductOrderModel `gorm:"foreignKey:OrderID;references:ID"`
}

func (OrderModel) TableName() string {
	return "orders"
}

type ProductOrderModel struct {
	ID        string    `gorm:"primaryKey;type:uuid;column:id"`
	OrderID   string    `gorm:"type:varchar(36);column:order_id"`
	ProductID string    `gorm:"type:uuid;column:product_id"`
	Quantity  int       `gorm:"column:quantity"`
	Price     int       `gorm:"column:price"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (ProductOrderModel) TableName() string {
	return "product_orders"
}
