package dto

import "time"

type UpdateUserInput struct {
	Email    *string `json:"email,omitempty"`
	FullName *string `json:"full_name,omitempty"`
}

type UserResponse struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	FullName  string    `json:"full_name"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
