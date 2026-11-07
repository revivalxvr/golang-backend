package request

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
