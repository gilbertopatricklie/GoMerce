package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
	ErrRefreshTokenReuse   = errors.New("refresh token reuse detected")
)

type RefreshTokenRepository struct {
	db *pgxpool.Pool
}

func NewRefreshTokenRepository(db *pgxpool.Pool) *RefreshTokenRepository {
	return &RefreshTokenRepository{
		db: db,
	}
}

type RefreshTokenData struct {
	UserID   int64
	Role     string
	FamilyID uuid.UUID
}

func (r *RefreshTokenRepository) Create(ctx context.Context, userID int64, tokenHash string, expiresAt time.Time) error {
	query := `INSERT INTO refresh_tokens (
			  user_id, token_hash, family_id, expires_at)
		      VALUES ($1, $2, $3, $4)`

	_, err := r.db.Exec(ctx, query, userID, tokenHash, uuid.New(), expiresAt)
	if err != nil {
		return fmt.Errorf("create refresh token: %w", err)
	}

	return nil
}


func (r *RefreshTokenRepository) Rotate(ctx context.Context, oldTokenHash string, newTokenHash string, newExpiresAt time.Time) (*RefreshTokenData, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin refresh transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	query := `SELECT rt.user_id, rt.family_id, rt.expires_at, rt.revoked_at, u.role, u.is_active
			  FROM refresh_tokens rt JOIN users u ON u.id = rt.user_id
		      WHERE rt.token_hash = $1 FOR UPDATE OF rt`

	var data RefreshTokenData
	var expiresAt time.Time
	var revokedAt *time.Time
	var isActive bool

	err = tx.QueryRow(ctx, query, oldTokenHash).Scan(
		&data.UserID,
		&data.FamilyID,
		&expiresAt,
		&revokedAt,
		&data.Role,
		&isActive,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrInvalidRefreshToken
	}
	if err != nil {
		return nil, fmt.Errorf("find refresh token: %w", err)
	}


	if revokedAt != nil {
		_, err := tx.Exec(ctx, `
			UPDATE refresh_tokens
			SET revoked_at = COALESCE(revoked_at, NOW())
			WHERE family_id = $1
		`, data.FamilyID)
		if err != nil {
			return nil, fmt.Errorf("revoke token family: %w", err)
		}

		if err := tx.Commit(ctx); err != nil {
			return nil, fmt.Errorf("commit token family revocation: %w", err)
		}

		return nil, ErrRefreshTokenReuse
	}

	if !expiresAt.After(time.Now()) || !isActive {
		return nil, ErrInvalidRefreshToken
	}

	_, err = tx.Exec(ctx, `
		UPDATE refresh_tokens
		SET revoked_at = NOW()
		WHERE token_hash = $1
	`, oldTokenHash)
	if err != nil {
		return nil, fmt.Errorf("revoke old refresh token: %w", err)
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO refresh_tokens (
			user_id, token_hash, family_id, expires_at
		)
		VALUES ($1, $2, $3, $4)
	`, data.UserID, newTokenHash, data.FamilyID, newExpiresAt)
	if err != nil {
		return nil, fmt.Errorf("create rotated refresh token: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit refresh rotation: %w", err)
	}

	return &data, nil
}