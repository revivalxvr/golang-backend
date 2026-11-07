package response

import "time"

type MidtransTx struct {
	Token       string `json:"token"`
	RedirectURL string `json:"redirect_url"`
}

type ProductItemResponse struct {
	ID           string    `json:"id"`
	ProductID    string    `json:"product_id"`
	ProductName  string    `json:"product_name"`
	ProductImage string    `json:"product_image"`
	Quantity     int       `json:"quantity"`
	Price        int       `json:"price"`
	CreatedAt    time.Time `json:"created_at"`
}

// OrderWithProducts merepresentasikan 1 data order utuh beserta daftar item produk di dalamnya
type OrderWithProducts struct {
	ID            string                `json:"id"`
	UserID        string                `json:"user_id"`
	GrossAmount   int                   `json:"gross_amount"`
	Status        string                `json:"status"`
	PaymentMethod string                `json:"payment_method"`
	PaymentUrl    string                `json:"payment_url"`
	CreatedAt     time.Time             `json:"created_at"`
	UpdatedAt     time.Time             `json:"updated_at"`
	Items         []ProductItemResponse `json:"items"`
}
