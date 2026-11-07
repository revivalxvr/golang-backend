package response

type UserBlockResponse struct {
	IsBlocked bool   `json:"is_blocked"`
	Reason    string `json:"reason"`
}
