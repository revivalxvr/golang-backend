package request

type RabbitMQLogRequest struct {
	Timestamp    string  `json:"@timestamp"`
	Exchange     string  `json:"exchange"`
	RoutingKey   string  `json:"routing_key"`
	Queue        string  `json:"queue"`
	ConsumerTag  string  `json:"consumer_tag"`
	Payload      any     `json:"payload"`
	ContentType  string  `json:"content_type"`
	DeliveryMode uint8   `json:"delivery_mode"`
	AppID        string  `json:"app_id"`
	Acked        bool    `json:"acked"`
	Latency      float64 `json:"latency"`
	RequestId    string  `json:"request_id,omitempty"`
	UserId       *string `json:"user_id,omitempty"`
	Error        *string `json:"error,omitempty"`
}
