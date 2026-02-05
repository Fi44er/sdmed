package auth_entity

type CodeType string

const (
	CodeTypeVerify CodeType = "verify"
	CodeTypeReset  CodeType = "reset"
)

type Code struct {
	Email string
	Code  string
	Type  CodeType
}
