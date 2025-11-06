package entity

import (
	"time"

	"github.com/google/uuid"
)

type FileStatus string

const (
	FileStatusTemporary FileStatus = "temporary"
	FileStatusPermanent FileStatus = "permanent"
)

type File struct {
	ID        string
	Name      string
	Data      []byte
	OwnerID   *string
	OwnerType *string
	Status    FileStatus
	ExpiresAt *time.Time
	CreatedAt time.Time
}

func (f *File) IsExpired() bool {
	if f.Status != FileStatusTemporary || f.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*f.ExpiresAt)
}

func (f *File) MarkAsPermanent(ownerID, ownerType string) {
	f.Status = FileStatusPermanent
	f.OwnerID = &ownerID
	f.OwnerType = &ownerType
	f.ExpiresAt = nil
}

func (f *File) GenerateName() error {
	f.Name = uuid.New().String()
	return nil
}
