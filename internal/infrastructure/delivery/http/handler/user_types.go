package handler

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type GetUsersResponse struct {
	Pageable
	Users []User `json:"users"`
}
