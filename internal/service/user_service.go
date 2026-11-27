package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"gomerce/internal/model"
	"gomerce/internal/repository"
	"gomerce/internal/security"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrAccountInactive    = errors.New("account is inactive")
)

type UserService struct {
	repo *repository.UserRepository
	refreshRepo *repository.RefreshTokenRepository
}

func NewUserService(repo *repository.UserRepository, refreshRepo *repository.RefreshTokenRepository) *UserService {
	return &UserService{
		repo: repo,
		refreshRepo: refreshRepo,
	}
}

//Signing up method
func (s *UserService) SignUp(ctx context.Context, name, email, password string) (*model.User, error) {
	name = strings.TrimSpace(name)
	email = strings.ToLower(strings.TrimSpace(email))

	if name == "" || email == "" || password == "" {
		return nil, errors.New("name, email, and password are required")
	}

	if len(password) > 72 {
		return nil, errors.New("password exceeds 72 bytes")
	}

	hash, err := security.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	user := &model.User{
		Name:         name,
		Email:        email,
		PasswordHash: hash,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, ErrEmailAlreadyExists
		}

		return nil, err
	}

	return user, nil
}

//Signing in method
func (s *UserService) SignIn(ctx context.Context,email, password string) (*model.User, error) {
	email = strings.ToLower(strings.TrimSpace(email))

	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrInvalidCredentials
		}

		return nil, err
	}

	if err := security.CheckPassword(password, user.PasswordHash); err != nil {
		return nil, ErrInvalidCredentials
	}

	if !user.IsActive {
		return nil, ErrAccountInactive
	}

	return user, nil
}

func (s *UserService) CreateRefreshToken(ctx context.Context, userID int64, ttl time.Duration) (string, error) {
	token, err := security.GenerateRefreshToken()
	if err != nil {
		return "", err
	}

	expiresAt := time.Now().Add(ttl)
	tokenHash := security.HashRefreshToken(token)

	if err := s.refreshRepo.Create(
		ctx,
		userID,
		tokenHash,
		expiresAt,
	); err != nil {
		return "", err
	}

	return token, nil
}

func (s *UserService) RotateRefreshToken(ctx context.Context, oldToken string, ttl time.Duration,) (int64, string, string, error) {
	newToken, err := security.GenerateRefreshToken()
	if err != nil {
		return 0, "", "", err
	}

	newExpiresAt := time.Now().Add(ttl)

	data, err := s.refreshRepo.Rotate(
		ctx,
		security.HashRefreshToken(oldToken),
		security.HashRefreshToken(newToken),
		newExpiresAt,
	)
	if err != nil {
		return 0, "", "", err
	}

	return data.UserID, data.Role, newToken, nil
}