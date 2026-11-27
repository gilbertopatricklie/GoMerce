package repository

import (
	"context"
	"fmt"

	"gomerce/internal/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	query := `INSERT INTO users (name, email, password_hash)
			  VALUES($1, $2, $3) RETURNING id, role, is_active, created_at, updated_at`

	err := r.db.QueryRow(ctx, query, user.Name, user.Email, user.PasswordHash).Scan(
		&user.ID, &user.Role, &user.IsActive, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}

	return nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `SELECT id, name, email, password_hash, role, is_active, created_at, updated_at
	   		  FROM users WHERE email = $1`

	var user model.User

	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows{
			return nil, pgx.ErrNoRows
		}

		return nil, fmt.Errorf("get user by email: %w", err)
	}

	return &user, nil
}