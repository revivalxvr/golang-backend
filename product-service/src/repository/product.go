package repository

import (
	"context"
	"net/http"
	"product-service/src/common/dto/request"
	"product-service/src/common/errors"
	"product-service/src/common/model"
	"product-service/src/common/utils"
	"sort"
	"strings"

	productpb "github.com/revivalxvr/protobuf/protogen/product"
	"google.golang.org/grpc/codes"

	"gorm.io/gorm"
)

type ProductRepository interface {
	Create(ctx context.Context, req *model.Product) error
	CreateCategory(ctx context.Context, req *model.Category) error
	FindByID(ctx context.Context, productID string, product *model.Product) error
	FindMany(ctx context.Context, req *request.FindManyProduct) ([]*model.Product, error)
	FindManyCategory(ctx context.Context, req *request.FindManyCategory) ([]*model.Category, error)
	FindCategoryByID(ctx context.Context, categoryID string, category *model.Category) error
	Update(ctx context.Context, product *model.Product) error
	Delete(ctx context.Context, id string) error
	ReduceStocks(ctx context.Context, req []*productpb.ProductOrder) error
	RollbackStocks(ctx context.Context, req []*productpb.ProductOrder) error
}
type productImpl struct {
	db *gorm.DB
}

func NewProduct(db *gorm.DB) ProductRepository {
	return &productImpl{
		db: db,
	}
}

func (p *productImpl) Create(ctx context.Context, req *model.Product) error {
	return p.db.WithContext(ctx).Create(req).Error
}

func (p *productImpl) FindByID(ctx context.Context, productID string, product *model.Product) error {
	return p.db.WithContext(ctx).Where("id = ?", productID).First(product).Error
}
func (p *productImpl) FindMany(ctx context.Context, req *request.FindManyProduct) ([]*model.Product, error) {
	limit, offset := utils.CreateLimitAndOffset(req.Page)

	var res []*model.Product

	// 1. Gunakan Model() dan Preload("Category") agar GORM otomatis mengambil relasi tabel categories
	query := p.db.WithContext(ctx).Model(&model.Product{}).Preload("Category")

	// 2. Tangani logika pencarian (search)
	if req.Search != "" {
		searchQuery := strings.Join(strings.Fields(req.Search), " & ")
		query = query.Where("to_tsvector('indonesian', deskripsi) @@ to_tsquery('indonesian', ?::text)", searchQuery)
	}

	// 3. Eksekusi query dengan limit dan offset
	err := query.Limit(limit).Offset(offset).Find(&res).Error
	return res, err
}
func (p *productImpl) FindManyCategory(ctx context.Context, req *request.FindManyCategory) ([]*model.Category, error) {
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

	err := p.db.WithContext(ctx).Raw(query, args...).Scan(&res).Error
	return res, err
}

func (p *productImpl) FindCategoryByID(ctx context.Context, categoryID string, category *model.Category) error {
	return p.db.WithContext(ctx).Where("id = ?", categoryID).First(category).Error
}
func (p *productImpl) Update(ctx context.Context, product *model.Product) error {
	return p.db.WithContext(ctx).Save(product).Error
}
func (p *productImpl) Delete(ctx context.Context, id string) error {
	// Menghapus data produk berdasarkan ID
	return p.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Product{}).Error
}

func (p *productImpl) CreateCategory(ctx context.Context, req *model.Category) error {
	return p.db.WithContext(ctx).Create(req).Error
}

func (p *productImpl) ReduceStocks(ctx context.Context, req []*productpb.ProductOrder) error {
	// SORT slice berdasarkan ProductId secara Asending (A-Z) untuk mencegah Deadlock!
	sort.Slice(req, func(i, j int) bool {
		return req[i].ProductId < req[j].ProductId
	})
	//end of deadlock
	return p.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, order := range req {
			result := tx.Model(&model.Product{}).
				Where("id = ? AND stock >= ?", order.ProductId, order.Quantity).
				Update("stock", gorm.Expr("stock - ?", order.Quantity))

			if result.Error != nil {
				return result.Error
			}

			if result.RowsAffected == 0 {
				return &errors.Response{
					Message:  "product not found or insufficient stock",
					HttpCode: http.StatusBadRequest,
					GRPCCode: codes.FailedPrecondition,
				}
			}
		}
		return nil
	})
}

func (p *productImpl) RollbackStocks(ctx context.Context, req []*productpb.ProductOrder) error {
	//mencegah deadlock
	sort.Slice(req, func(i, j int) bool {
		return req[i].ProductId < req[j].ProductId
	})
	//end of deadlock
	return p.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, order := range req {
			result := tx.Model(&model.Product{}).
				Where("id = ?", order.ProductId).
				Update("stock", gorm.Expr("stock + ?", order.Quantity))

			if result.Error != nil {
				return result.Error
			}
		}
		return nil
	})
}
