package domain

import "time"

type UserRole string

const (
	RoleUser  UserRole = "user"
	RoleAdmin UserRole = "admin"
)

type User struct {
	ID           int       `json:"id"`
	Email        string    `json:"email"`
	Password     string    `json:"-"`
	FullName     string    `json:"full_name"`
	PasswordHash string    `json:"-"`
	Role         UserRole  `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type UserRepository interface {
	Create(user *User) error
	GetByID(id int) (*User, error)
	GetByEmail(email string) (*User, error)
	GetAll() ([]*User, error)
	Update(user *User) error
	Delete(id int) error
}
