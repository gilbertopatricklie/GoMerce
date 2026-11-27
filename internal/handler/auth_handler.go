package handler

import (
	"errors"
	"net/http"
	"net/mail"
	"time"

	"gomerce/internal/model"
	"gomerce/internal/repository"
	"gomerce/internal/response"
	"gomerce/internal/security"
	"gomerce/internal/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service *service.UserService
	tokenManager *security.TokenManager
	refreshTTL time.Duration
}

func NewAuthHandler(userService *service.UserService, tokenManager *security.TokenManager, refreshTTL time.Duration) *AuthHandler {
	return &AuthHandler{
		service: userService,
		tokenManager: tokenManager,
		refreshTTL: refreshTTL,
	}
}

type SignUpRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type SignInRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type SignInResponse struct {
	AccessToken string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
	TokenType   string      `json:"token_type"`
	ExpiresIn   int64       `json:"expires_in"`
	User        *model.User `json:"user"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func validEmail(email string) bool {
	parsed, err := mail.ParseAddress(email)
	return err == nil && parsed.Address == email
}

//Method to sign up
func (h *AuthHandler) SignUp(c *gin.Context) {
	var req SignUpRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "VALIDATION_ERROR", "Invalid request body", nil)
		return
	}

	if !validEmail(req.Email) {
		response.Fail(c, 400, "INVALID_EMAIL", "Invalid email address", nil)
		return
	}

	if len(req.Password) < 8 || len(req.Password) > 72 {
		response.Fail(c, 400, "INVALID_PASSWORD", "Password must be between 8 and 72 bytes", nil)
		return
	}

	user, err := h.service.SignUp(
		c.Request.Context(),
		req.Name,
		req.Email,
		req.Password,
	)

	if err != nil {
		if errors.Is(err, service.ErrEmailAlreadyExists) {
			response.Fail(c, 409, "EMAIL_EXISTS", "Email already registered", nil)
			return
		}

		response.Fail(c, 500, "INTERNAL_SERVER_ERROR", "Failed to register user", nil)
		return
	}

	response.Success(c, 201, user)
}

//Method to sign in
func (h *AuthHandler) SignIn(c *gin.Context) {
	var req SignInRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "VALIDATION_ERROR", "Invalid request body", nil)
		return
	}

	if !validEmail(req.Email) {
		response.Fail(c, 400, "INVALID_EMAIL", "Invalid email address", nil)
		return
	}

	user, err := h.service.SignIn(
		c.Request.Context(),
		req.Email,
		req.Password,
	)

	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			response.Fail(c, 401, "INVALID_CREDENTIALS", "Invalid email or password", nil)
			return
		}

		if errors.Is(err, service.ErrAccountInactive) {
			response.Fail(c, 403, "ACCOUNT_INACTIVE", "Account is inactive", nil)
			return
		}

		response.Fail(c, 500, "INTERNAL_SERVER_ERROR", "Failed to sign in", nil)
		return
	}

	accessToken, expiresAt, err := h.tokenManager.GenerateAccessToken(user.ID,user.Role)

	if err != nil {
		response.Fail(c, 500, "INTERNAL_ERROR", "Failed to generate access token", nil)
		return
	}

	refreshToken, err := h.service.CreateRefreshToken(c.Request.Context(), user.ID, h.refreshTTL,)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to generate refresh token", nil)
		return
	}

	result := SignInResponse{
		AccessToken: accessToken,
		RefreshToken: refreshToken,
		TokenType:   "Bearer",
		ExpiresIn:   int64(time.Until(expiresAt).Seconds()),
		User:        user,
	}

	response.Success(c, 200, result)
}


func (h *AuthHandler) Refresh(c *gin.Context) {
	var req RefreshRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "VALIDATION_ERROR", "Refresh token is required", nil)
		return
	}

	userID, role, newRefreshToken, err := h.service.RotateRefreshToken(
		c.Request.Context(),
		req.RefreshToken,
		h.refreshTTL,
	)

	if err != nil {
		if errors.Is(err, repository.ErrInvalidRefreshToken) ||
			errors.Is(err, repository.ErrRefreshTokenReuse) {
			response.Fail(c, http.StatusUnauthorized, "INVALID_REFRESH_TOKEN", "Refresh token is invalid or expired", nil)
			return
		}

		response.Fail(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to refresh token", nil)
		return
	}

	accessToken, expiresAt, err := h.tokenManager.GenerateAccessToken(
		userID,
		role,
	)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to generate access token", nil)
		return
	}

	result := SignInResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(time.Until(expiresAt).Seconds()),
		User: &model.User{
			ID:   userID,
			Role: role,
		},
	}

	response.Success(c, http.StatusOK, result)
}