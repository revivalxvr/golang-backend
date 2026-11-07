package request

type RegisterRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=20"`
}

type SendOtpReq struct {
	Email string `json:"email" validate:"required,email"`
	Otp   string `json:"otp" validate:"required,min=6,max=6"`
}

type VerifyRegister struct {
	Email string `json:"email"`
	Otp   string `json:"otp"`
}
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=20"`
}
