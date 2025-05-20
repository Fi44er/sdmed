package dto

type LoginDTO struct {
	Email     string `json:"email" validate:"required"`
	Password  string `json:"password" validate:"required"`
	UserAgent string `json:"user_agent"`
}

type LogoutDTO struct {
	RefreshToken    string `json:"refresh_token" validate:"required"`
	AccessTokenUUID string `json:"access_token_uuid" validate:"required"`
}

type VerifyCodeDTO struct {
	Email string `json:"email" validate:"required,email"`
	Code  string `json:"code" validate:"required"`
}

type RefreshTokenDTO struct {
	RefreshToken string `json:"refresh_token"`
	UserAgent    string `json:"user_agent"`
}

type RegisterDTO struct {
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required"`
	FIO         string `json:"fio" validate:"required"`
	PhoneNumber string `json:"phone_number" validate:"required"`
}

type SendCodeDTO struct {
	Email string `json:"email" validate:"required,email"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
