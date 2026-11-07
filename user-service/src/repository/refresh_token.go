package repository

import (
	"context"
	"user-service/src/common/model"

	"gorm.io/gorm"
)

// jangan lupa tambahkan RefreshToken di factory
type RefreshToken interface {
	Create(c context.Context, data *model.RefreshToken) error
	Delete(c context.Context, whereClause string, args []any) error
	Find(c context.Context, refreshToken string) (*model.RefreshTokenWithRole, error)
	Update(c context.Context, data *model.RefreshToken, oldToken string) error
}

type refreshTokenImpl struct {
	db *gorm.DB
}

func NewRefreshToken(db *gorm.DB) RefreshToken {
	return &refreshTokenImpl{
		db: db,
	}
}

func (r *refreshTokenImpl) Create(c context.Context, data *model.RefreshToken) error {
	return r.db.WithContext(c).Create(data).Error
}

func (r *refreshTokenImpl) Delete(c context.Context, whereClause string, args []any) error {
	return r.db.WithContext(c).Where(whereClause, args...).Delete(&model.RefreshToken{}).Error
}

func (r *refreshTokenImpl) Find(c context.Context, refreshToken string) (*model.RefreshTokenWithRole, error) {
	res := new(model.RefreshTokenWithRole)

	if err := r.db.WithContext(c).
		Table("refresh_tokens AS rt").
		Select("rt.id, rt.token, rt.user_id, u.role"). // Tambahkan rt.id jika diperlukan
		Joins("JOIN users AS u ON rt.user_id = u.id").
		Where("rt.token = ?", refreshToken).
		First(res).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	if res.Token == "" {
		return nil, nil
	}
	return res, nil
}

func (r *refreshTokenImpl) Update(c context.Context, data *model.RefreshToken, oldToken string) error {
	// Melakukan UPDATE langsung pada record dengan token lama
	return r.db.WithContext(c).
		Model(&model.RefreshToken{}).
		Where("token = ?", oldToken).
		Updates(map[string]any{
			"token":      data.Token,
			"updated_at": data.UpdatedAt,
		}).Error
}
