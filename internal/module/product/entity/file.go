package product_entity

type File struct {
	ID        string
	Name      string
	Data      []byte
	OwnerID   *string
	OwnerType *string
}
