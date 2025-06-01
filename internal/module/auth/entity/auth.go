package entity

type VerifyCode struct {
	Email string
	Code  string
}

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

type UserSesion struct {
	UserID       string
	RefreshToken string
}
