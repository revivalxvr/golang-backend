package request

import (
	"github.com/google/uuid"
)

type OrderCreated struct {
	EventID    string      `json:"event_id" validate:"required"`
	OrderID    string      `json:"order_id" validate:"required"`
	UserID     uuid.UUID   `json:"user_id" validate:"required"`
	ProductIDs []uuid.UUID `json:"product_ids" validate:"required,min=1,dive,required"`
}
