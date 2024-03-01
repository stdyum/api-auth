package dto

type SignUpRequestDTO struct {
	Email    string `json:"email"`
	Login    string `json:"login"`
	Password string `json:"password"`
	Picture  string `json:"picture"`
}

type LoginRequestDTO struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type ConfirmEmailByCodeRequestDTO struct {
	Code string `json:"code"`
}

type ResetPasswordRequestDTO struct {
	Login string `json:"login"`
	Email string `json:"email"`
}

type ResetPasswordByCodeRequestDTO struct {
	Code     string `json:"code"`
	Password string `json:"password"`
}

type UpdateTokenRequestDTO struct {
	Token string `json:"refresh"`
}
