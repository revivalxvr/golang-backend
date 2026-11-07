package request

type MidtransTx struct {
	OrderId           string             `json:"order_id"`
	StatusCode        string             `json:"status_code"`
	GrossAmount       string             `json:"gross_amount"`
	TransactionStatus string             `json:"transaction_status"`
	SingratureKey     string             `json:"signature_key"`
	PaymentType       string             `json:"payment_type"`
	FraudStatus       string             `json:"fraud_status"`
	Metadata          MidtransTxMetadata `json:"metadata"`
}

type MidtransTxMetadata struct {
	OriginalOrderId string `json:"original_order_id"`
}
