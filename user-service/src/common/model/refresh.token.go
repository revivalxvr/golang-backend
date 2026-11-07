package model

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	ID        uuid.UUID `gorm:"primaryKey;type:uuid" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	Token     string    `gorm:"type:varchar(512);not null" json:"token"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
type RefreshTokenWithRole struct {
	ID     uuid.UUID `gorm:"primaryKey;type:uuid" json:"id"`
	UserID uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	Token  string
	Role   string
}
