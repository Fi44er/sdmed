package dto

type UploadFiles struct {
	Name string `json:"name"`
	Data []byte `json:"data"`
}
