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

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

func (r *UserRepository) Create(user *domain.User) error {
	query := `
		INSERT INTO users (email, password_hash, full_name, role)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	err := r.pool.QueryRow(
		context.Background(),
		query,
		user.Email,
		user.PasswordHash,
		user.FullName,
		user.Role,
	).Scan(&user.ID)

	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) GetByEmail(email string) (*domain.User, error) {
	user := &domain.User{}

	query := `
		SELECT id, email, password_hash, full_name, role, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	err := r.pool.QueryRow(
		context.Background(),
		query,
		email,
	).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.FullName,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) GetByID(id int) (*domain.User, error) {
	user := &domain.User{}

	query := `
		SELECT id, email, password_hash, full_name, role, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	err := r.pool.QueryRow(
		context.Background(),
		query,
		id,
	).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.FullName,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) Update(user *domain.User) error {
	query := `
		UPDATE users
		SET email = $1, full_name = $2, updated_at = $3
		WHERE id = $4
	`

	result, err := r.pool.Exec(
		context.Background(),
		query,
		user.Email,
		user.FullName,
		time.Now(),
		user.ID,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return apperrors.ErrUserNotFound
	}

	return nil
}

func (r *UserRepository) GetAll() ([]*domain.User, error) {
	query := `
		SELECT id, email, password_hash, full_name, role, created_at, updated_at
		FROM users
	`

	rows, err := r.pool.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		user := &domain.User{}
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.PasswordHash,
			&user.FullName,
			&user.Role,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *UserRepository) Delete(id int) error {
	query := `DELETE FROM users WHERE id = $1`
	result, err := r.pool.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return apperrors.ErrUserNotFound
	}

	return nil
}

func (r *UserRepository) UpdateFields(id int, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return nil
	}

	query := `
		UPDATE users
		SET email = COALESCE($1, email),
			full_name = COALESCE($2, full_name),
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $3
	`

	result, err := r.pool.Exec(
		context.Background(),
		query,
		updates["email"],
		updates["full_name"],
		id,
	)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return apperrors.ErrUserNotFound
	}

	return nil
}
