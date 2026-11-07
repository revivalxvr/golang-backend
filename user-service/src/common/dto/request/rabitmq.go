package request

type RabbitMqMessage struct {
	RequestID any `json:"request_id"`
	UserID    any `json:"user_id"`
	Message   any `json:"message"`
}
