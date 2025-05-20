package entity

type VerifyCode struct {
	Emai string
	Code string
}

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

type UserSesion struct {
	ID           string
	UserID       string // может быть nil
	UserAgent    string
	CreatedAt    int64
	LastActivity int64
	IsRevoked    bool
	Data         map[string]interface{}
}
