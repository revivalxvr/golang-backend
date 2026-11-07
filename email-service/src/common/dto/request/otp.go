package request

type SendOtp struct {
	Email string `json:"email" validate:"required,email"`
	Otp   string `json:"otp" validate:"required,min=6,max=6"`
	Year  int    `validate:"omitempty"`
}
