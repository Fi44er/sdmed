package entity

import "github.com/google/uuid"

type File struct {
	ID        string
	Name      string
	OwnerID   string
	OwnerType string
	Data      []byte
}

func (f *File) GenerateName() error {
	f.Name = uuid.New().String()
	return nil
}
