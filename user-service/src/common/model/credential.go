package model

import (
	"time"

	"github.com/google/uuid"
)

type Credential struct {
	UserID    uuid.UUID `gorm:"primaryKey;column:user_id" json:"user_id"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Tambahkan relasi ke tabel User agar bisa di-preload
	User User `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
}

// Beritahu GORM nama tabelnya jika tidak plural otomatis
func (Credential) TableName() string {
	return "credentials"
}

type CredentialWithRole struct {
	UserID   uuid.UUID
	Password string
	Role     string
}
