package repository

import (
	"context"
	"product-service/src/common/dto/request"
	"product-service/src/common/model"
	"product-service/src/common/utils"
	"strings"

	"gorm.io/gorm"
)

type CategoryRepository interface {
	FindManyCategory(ctx context.Context, req *request.FindManyCategory) ([]*model.Category, error)
	FindCategoryByID(ctx context.Context, categoryID string, category *model.Category) error
	CreateCategory(ctx context.Context, req *model.Category) error
	UpdateCategory(ctx context.Context, category *model.Category) error
	DeleteCategory(ctx context.Context, category *model.Category) error
}
type categoryImpl struct {
	db *gorm.DB
}

func NewCategory(db *gorm.DB) CategoryRepository {
	return &categoryImpl{
		db: db,
	}
}

func (r *categoryImpl) FindManyCategory(ctx context.Context, req *request.FindManyCategory) ([]*model.Category, error) {
	limit, offset := utils.CreateLimitAndOffset(req.Page)

	var args []any

	// Berikan spasi di akhir string awal
	query := "SELECT * FROM categories "

	if req.Search != "" {
		// Gunakan to_tsquery(?::text) atau plainto_tsquery(?) agar PostgreSQL tahu tipe datanya
		// Dan pastikan ada spasi di awal/akhir string
		query += "WHERE to_tsvector('indonesian', name) @@ to_tsquery('indonesian', ?::text) "
		args = append(args, strings.Join(strings.Fields(req.Search), " & "))
	}

	// Tambahkan spasi sebelum LIMIT
	query += "LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	var res []*model.Category

	err := r.db.WithContext(ctx).Raw(query, args...).Scan(&res).Error
	return res, err
}

func (r *categoryImpl) FindCategoryByID(ctx context.Context, categoryID string, category *model.Category) error {
	return r.db.WithContext(ctx).Where("id = ?", categoryID).First(category).Error
}
func (r *categoryImpl) CreateCategory(ctx context.Context, req *model.Category) error {
	return r.db.WithContext(ctx).Create(req).Error
}

func (r *categoryImpl) UpdateCategory(ctx context.Context, category *model.Category) error {
	return r.db.WithContext(ctx).Save(category).Error
}

func (r *categoryImpl) DeleteCategory(ctx context.Context, category *model.Category) error {
	return r.db.WithContext(ctx).Delete(category).Error
}
