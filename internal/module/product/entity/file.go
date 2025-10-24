package product_entity

type File struct {
	ID        string
	Name      string
	OwnerID   string
	OwnerType string
	Data      []byte
}
