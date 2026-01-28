package auth_entity

type CodeType string

const (
	CodeTypeVerify CodeType = "verify"
	CodeTypeReset  CodeType = "reset"
)

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

type Code struct {
	Email string
	Code  string
	Type  CodeType
}

type TokenDetails struct {
	Token     *string
	TokenUUID string
	UserID    string
	ExpiresIn *int64
	DeviceID  string
}
