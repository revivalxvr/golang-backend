package model

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name       string    `json:"name"`
	Price      float64   `json:"price"`
	Sku        string    `json:"sku"`
	ImageId    string    `json:"image_id"`
	Image      string    `json:"image"`
	Stock      int       `json:"stock"`
	CategoryID int       `json:"category_id"` // Menggunakan int biasa karena NOT NULL
	Category   *Category `gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"category,omitempty"`
	Length     float32   `json:"length"`
	Width      float32   `json:"width"`
	Height     float32   `json:"height"`
	Deskripsi  string    `json:"deskripsi"`
	Badge      string    `json:"badge"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type Category struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"type:varchar(255);not null" json:"name"`
	Slug      string    `gorm:"type:varchar(255);not null" json:"slug"`
	Deskripsi string    `gorm:"type:varchar(255);not null" json:"deskripsi"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
