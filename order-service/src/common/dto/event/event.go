package event

type OrderCreated struct {
	EventID    string   `json:"event_id"`
	OrderID    string   `json:"order_id"`
	UserID     string   `json:"user_id"`
	ProductIDs []string `json:"product_ids"`
}
