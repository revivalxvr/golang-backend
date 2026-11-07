package model

import (
	"time"

	"github.com/google/uuid"
)

type ProductOrder struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	OrderID   string    `json:"order_id"`
	ProductID string    `json:"product_id"`
	Quantity  int       `json:"quantity"`
	Price     int       `json:"price"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
