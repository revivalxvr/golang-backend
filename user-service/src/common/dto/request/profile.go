package request

type UpdateProfileRequest struct {
	Name  *string `json:"name" validate:"omitempty,min=3"`
	Email *string `json:"email" validate:"omitempty,email"`
}

// omitempty: Artinya opsional (boleh tidak dikirim / nil). Tapi jika dikirim, jalankan validasi setelahnya.

// min:3: Jika nama dikirim, minimal 3 karakter.

// email: Jika email dikirim, formatnya harus berupa email yang valid. (Secara otomatis menolak string kosong "" atau format ngawur seperti "abc").

type ChangePasswordRequest struct {
	OldPassword     string `json:"old_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=NewPassword"`
}
