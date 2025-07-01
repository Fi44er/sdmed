package dto

import "github.com/Fi44er/sdmed/internal/module/file/entity"

type UploadFiles struct {
	File entity.File
	Data []byte `json:"data"`
}
