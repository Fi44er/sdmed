package entity

type VerifyCode struct {
	Email string
	Code  string
}

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

type UserSession struct {
	UserID       string
	RefreshToken string
}

type SendCode struct {
	Email string
	Code  string
}

type TokenDetails struct {
	Token     *string
	TokenUUID string
	UserID    string
	ExpiresIn *int64
}
