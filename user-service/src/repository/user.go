package repository

import (
	"context"
	"errors"
	"user-service/src/common/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// jgn lupa tambahkan User di factory
type User interface {
	Find(ctx context.Context, whereClause string, args []any) (*model.User, error)
	Create(ctx context.Context, user *model.User, password string) error
	FindByID(userID string) (*model.User, error)
	FindByEmailAndNotID(email string, notID string) (*model.User, error)
	Update(user *model.User) error
}

type userImpl struct {
	db *gorm.DB
}

func NewUser(db *gorm.DB) User {
	return &userImpl{
		db: db,
	}
}

func (r *userImpl) Find(ctx context.Context, whereClause string, args []any) (*model.User, error) {
	var user model.User

	if err := r.db.WithContext(ctx).Where(whereClause, args...).First(&user).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, nil
	}

	if user.ID == uuid.Nil {
		return nil, nil
	}
	return &user, nil
}

func (r *userImpl) Create(ctx context.Context, user *model.User, password string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return err
		}
		err := tx.Create(&model.Credential{
			UserID:   user.ID,
			Password: password,
		}).Error
		return err
	})
}

// FindByID
func (r *userImpl) FindByID(userID string) (*model.User, error) {
	var user model.User
	if err := r.db.First(&user, "id = ?", userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *userImpl) FindByEmailAndNotID(email string, notID string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("email = ? AND id != ?", email, notID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Tidak ditemukan berarti aman (email belum dipakai)
		}
		return nil, err // Error lain saat query
	}
	return &user, nil
}

// if err := s.UserRepository.Update(user);
func (r *userImpl) Update(user *model.User) error {
	if err := r.db.Save(user).Error; err != nil {
		return err
	}
	return nil
}
