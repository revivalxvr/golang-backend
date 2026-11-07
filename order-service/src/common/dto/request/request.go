package request

import "order-service/src/common/model"

type CreateOrder struct {
	Order         *model.Order          `json:"order" validate:"required"`
	ProductOrders []*model.ProductOrder `json:"product_orders" validate:"required,dive"`
}

type MidtransTxNotif struct {
	OrderId           string `json:"order_id"`
	StatusCode        string `json:"status_code"`
	GrossAmount       string `json:"gross_amount"`
	TransactionStatus string `json:"transaction_status"`
	SingratureKey     string `json:"signature_key"`
	PaymentType       string `json:"payment_type"`
	FraudStatus       string `json:"fraud_status"`
}
