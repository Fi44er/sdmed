package user_dto

type UserDTO struct {
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	Patronymic  string `json:"patronymic"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
	Password    string `json:"password"`
}

type UserResponse struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Surname     string   `json:"surname"`
	Patronymic  string   `json:"patronymic"`
	Email       string   `json:"email"`
	PhoneNumber string   `json:"phone_number"`
	Role        []string `json:"role"`
}
