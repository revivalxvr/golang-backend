package response

type CommonResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}
