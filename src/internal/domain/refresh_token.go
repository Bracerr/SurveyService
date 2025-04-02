package domain

import "time"

type RefreshToken struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

type RefreshTokenRepository interface {
	Create(token *RefreshToken) error
	GetByToken(token string) (*RefreshToken, error)
	UpdateToken(userID int, token string, expiresAt time.Time) error
}
