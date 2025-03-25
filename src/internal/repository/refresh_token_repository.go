package repository

import (
	"context"
	"errors"
	"time"

	"survey-project/src/internal/apperrors"
	"survey-project/src/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RefreshTokenRepository struct {
	pool *pgxpool.Pool
}

func NewRefreshTokenRepository(pool *pgxpool.Pool) *RefreshTokenRepository {
	return &RefreshTokenRepository{pool: pool}
}

func (r *RefreshTokenRepository) Create(token *domain.RefreshToken) error {
	query := `
		INSERT INTO refresh_tokens (user_id, token, expires_at)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	err := r.pool.QueryRow(
		context.Background(),
		query,
		token.UserID,
		token.Token,
		token.ExpiresAt,
	).Scan(&token.ID)

	return err
}

func (r *RefreshTokenRepository) GetByToken(token string) (*domain.RefreshToken, error) {
    refreshToken := &domain.RefreshToken{}

    query := `
        SELECT id, user_id, token, expires_at, created_at
        FROM refresh_tokens
        WHERE token = $1
    `

    err := r.pool.QueryRow(
        context.Background(),
        query,
        token,
    ).Scan(
        &refreshToken.ID,
        &refreshToken.UserID,
        &refreshToken.Token,
        &refreshToken.ExpiresAt,
        &refreshToken.CreatedAt,
    )

    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, apperrors.ErrInvalidToken
        }
        return nil, err
    }

    return refreshToken, nil
}

func (r *RefreshTokenRepository) UpdateToken(userID int, token string, expiresAt time.Time) error {
    query := `
        INSERT INTO refresh_tokens (user_id, token, expires_at)
        VALUES ($1, $2, $3)
        ON CONFLICT (user_id) DO UPDATE
        SET token = $2, expires_at = $3
    `

    _, err := r.pool.Exec(
        context.Background(),
        query,
        userID,
        token,
        expiresAt,
    )

    return err
}