package request

type MidtransTransactionDetails struct {
	OrderId     string `json:"order_id"`
	GrossAmount int    `json:"gross_amount"`
}
type MidtransTransactionMetadata struct {
	OriginalOrderId string `json:"original_order_id"`
}

type MidtransTransaction struct {
	TransactionDetails MidtransTransactionDetails  `json:"transaction_details"`
	Metadata           MidtransTransactionMetadata `json:"metadata"`
}

// --- Tambahan untuk Snap Token ---
type SnapCustomerDetails struct {
	FirstName string `json:"first_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
}
type SnapItemDetail struct {
	ID    string `json:"id"`
	Price int64  `json:"price"`
	Qty   int    `json:"quantity"`
	Name  string `json:"name"`
}
type SnapTransactionRequest struct {
	TransactionDetails struct {
		OrderId     string `json:"order_id"`
		GrossAmount int64  `json:"gross_amount"`
	} `json:"transaction_details"`
	CustomerDetails *SnapCustomerDetails `json:"customer_details,omitempty"`
	ItemDetails     []SnapItemDetail     `json:"item_details,omitempty"`
	// Tambahkan metadata di level root agar terbaca oleh Midtrans Snap
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}
