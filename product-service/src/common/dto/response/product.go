package response

import "time"

// Di dalam package response (misal: internal/delivery/http/response/product.go)
type ProductResponse struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Price      float64   `json:"price"`
	Sku        string    `json:"sku"`
	ImageId    string    `json:"image_id"`
	Image      string    `json:"image"`
	Stock      int       `json:"stock"`
	CategoryID int       `json:"category_id"`
	Category   string    `json:"category"` // Diubah menjadi string untuk menampung nama category
	Length     float32   `json:"length"`
	Width      float32   `json:"width"`
	Height     float32   `json:"height"`
	Deskripsi  string    `json:"deskripsi"`
	Badge      string    `json:"badge"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
