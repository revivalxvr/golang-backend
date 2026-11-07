package repository

import (
	"context"
	"user-service/src/common/model"

	"gorm.io/gorm"
)

type UserBlock interface {
	Find(ctx context.Context, userId string) (*model.UserBlock, error)
}
type userBlockImpl struct {
	db *gorm.DB
}

func NewUserBlocked(db *gorm.DB) UserBlock {
	return &userBlockImpl{
		db: db,
	}
}

func (r *userBlockImpl) Find(ctx context.Context, userId string) (*model.UserBlock, error) {
	var userBlock model.UserBlock
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND is_active = ?", userId, true).
		First(&userBlock).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
	}
	if userBlock.ID == 0 {
		return nil, nil
	}
	return &userBlock, nil
}
