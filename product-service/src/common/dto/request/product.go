package request

type CreateProduct struct {
	Name       string  `json:"name" validate:"required"`
	Price      float64 `json:"price" validate:"required"`
	Sku        string  `json:"sku"`
	ImageId    string  `json:"-" validate:"required"`
	Image      string  `json:"-" validate:"required"`
	Stock      int     `json:"stock" validate:"required"`
	CategoryID int     `json:"category_id" form:"category_id" validate:"required,min=1"`
	Length     float32 `json:"length"`
	Width      float32 `json:"width"`
	Height     float32 `json:"height"`
	Deskripsi  string  `json:"deskripsi"`
	Badge      string  `json:"badge"`
}

type UpdateProduct struct {
	Name       string  `form:"name" json:"name"`
	Price      float64 `form:"price" json:"price"`
	Sku        string  `form:"sku" json:"sku"`
	Stock      int     `form:"stock" json:"stock"`
	CategoryID int     `form:"category_id" json:"category_id" validate:"min=1"`
	Length     float32 `form:"length" json:"length"`
	Width      float32 `form:"width" json:"width"`
	Height     float32 `form:"height" json:"height"`
	Deskripsi  string  `form:"deskripsi" json:"deskripsi"`
	Badge      string  `form:"badge" json:"badge"`
	// Image & ImageId tidak wajib di struct request karena ditangani via file multipart terpisah
}
type FindManyProduct struct {
	Page   int    `json:"page" validate:"min=1"`
	Search string `json:"search"`
}
type FindManyCategory struct {
	Page   int    `json:"page" validate:"min=1"`
	Search string `json:"search"`
}
type CreateCategory struct {
	Name      string `gorm:"type:varchar(255);not null" json:"name"`
	Slug      string `gorm:"type:varchar(255);not null" json:"slug"`
	Deskripsi string `gorm:"type:varchar(255);not null" json:"deskripsi"`
	IsActive  bool   `json:"is_active"`
}
type UpdateCategory struct {
	Name      string `gorm:"type:varchar(255);not null" json:"name"`
	Slug      string `gorm:"type:varchar(255);not null" json:"slug"`
	Deskripsi string `gorm:"type:varchar(255);not null" json:"deskripsi"`
	IsActive  *bool  `json:"is_active"`
}
