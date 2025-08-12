package domain

import (
	"errors"
	"time"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
)

type User struct {
	ID          string
	Name        string
	Email       string
	PhoneNumber string
	Password    string
	Role        string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   time.Time
}

func (u *User) IsAdmin() bool {
	return u.Role == "admin"
}

func (u *User) IsUser() bool {
	return u.Role == "user"
}

func (u *User) IsGuest() bool {
	return u.Role == "guest"
}

func (u *User) IsActive() bool {
	return u.DeletedAt.IsZero()
}

func (u *User) IsDeleted() bool {
	return !u.DeletedAt.IsZero()
}
