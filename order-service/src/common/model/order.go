package model

import (
	"time"
)

type Order struct {
	ID            string    `gorm:"primaryKey" json:"id"` //kita gunakan nano id untuk id orderan(transaksi) tingkat keamanan setara UUID v4
	UserID        string    `json:"user_id"`
	GrossAmount   int       `json:"gross_amount"`
	Status        string    `json:"status"`
	PaymentMethod string    `json:"payment_method"`
	PaymentUrl    string    `json:"payment_url"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
