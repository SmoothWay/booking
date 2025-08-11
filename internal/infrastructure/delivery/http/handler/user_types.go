package handler

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type GetUsersResponse struct {
	Pageable
	Users []User `json:"users"`
}

type CreateUserRequest struct {
	Name        string `json:"name" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	PhoneNumber string `json:"phone" validate:"required,min=10,max=15"`
	Role        string `json:"role" validate:"required,oneof=admin user"`
	Password    string `json:"password" validate:"required,min=8"`
}
