package request

type AddCart struct {
	UserID    string `json:"user_id"`
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}
