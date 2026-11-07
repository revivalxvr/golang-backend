package request

import "google.golang.org/grpc/codes"

type GrpcLog struct {
	Timestamp  string     `json:"@timestamp"`
	RequestID  *string    `json:"request_id"`
	UserID     *string    `json:"user_id"`
	AppID      *string    `json:"app_id"`
	Method     string     `json:"method"`
	Body       any        `json:"body"`
	StatusCode codes.Code `json:"status_code"`
	Latency    float64    `json:"latency"`
	Error      string     `json:"error,omitempty"`
}

type RestFullLog struct {
	RequestID  string  `json:"request_id"`
	Timestamp  string  `json:"@timestamp"`
	Method     string  `json:"method"`
	Url        string  `json:"url"`
	Body       any     `json:"body"`
	StatusCode int     `json:"status_code"`
	ClientIP   string  `json:"client_ip"`
	UserAgent  string  `json:"user_agent"`
	UserID     *string `json:"user_id"`
	Latency    float64 `json:"latency"`
	Error      string  `json:"error"`
}
