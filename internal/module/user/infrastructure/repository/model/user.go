package user_model

type User struct {
	ID           string `gorm:"primaryKey;type:uuid;default:uuid_generate_v4();"`
	Email        string `gorm:"type:varchar(100);unique;not null"`
	PasswordHash string `gorm:"type:varchar(255);not null"`
	Name         string `gorm:"type:varchar(255);not null"`
	Surname      string `gorm:"type:varchar(255);"`
	Patronymic   string `gorm:"type:varchar(255);"`
	PhoneNumber  string `gorm:"type:varchar(255);"`
	Roles        []Role `gorm:"many2many:user_module.user_roles;"`
}

func (User) TableName() string {
	return "user_module.users"
}
