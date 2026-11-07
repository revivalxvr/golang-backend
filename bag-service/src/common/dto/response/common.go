package response

type CommonResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}
type CartItemResponse struct {
	ProductID   string  `json:"product_id"`
	ProductName string  `json:"product_name"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
	SubTotal    float64 `json:"sub_total"`
	Stock       int32   `json:"stock"`
}
