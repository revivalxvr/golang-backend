package repository

import (
	"context"
	"errors"
	"user-service/src/common/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// jgn lupa tambahkan Credential di factory
type Credential interface {
	Find(ctx context.Context, whereClause string, args []any) (*model.CredentialWithRole, error)
	FindByUserID(userID string) (*model.Credential, error)
	Update(credential *model.Credential) error
}

type credentialImpl struct {
	db *gorm.DB
}

func NewCredential(db *gorm.DB) Credential {
	return &credentialImpl{
		db: db,
	}
}

func (r *credentialImpl) Find(ctx context.Context, whereClause string, args []any) (*model.CredentialWithRole, error) {
	var credential model.CredentialWithRole
	if err := r.db.WithContext(ctx).
		Table("credentials AS c").
		Select("c.*, r.role").
		Joins("JOIN users AS r ON c.user_id = r.id").
		Where(whereClause, args...).
		Scan(&credential).Error; err != nil {
		return nil, err
	}

	if credential.UserID == uuid.Nil {
		return nil, nil
	}
	return &credential, nil
}
func (r *credentialImpl) FindByUserID(userID string) (*model.Credential, error) {
	var credential model.Credential
	// Tambahkan .Preload("User") untuk mengambil data relasi ke tabel users
	if err := r.db.Preload("User").Where("user_id = ?", userID).First(&credential).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("credential not found")
		}
		return nil, err
	}
	return &credential, nil
}
func (r *credentialImpl) Update(credential *model.Credential) error {
	// GORM Save akan melakukan update berdasarkan primary key (user_id)
	return r.db.Save(credential).Error
}
