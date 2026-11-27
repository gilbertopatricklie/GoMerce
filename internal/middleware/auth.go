package middleware

import (
	"strconv"
	"strings"

	"gomerce/internal/response"
	"gomerce/internal/security"

	"github.com/gin-gonic/gin"
)

func AuthRequired(tokenManager *security.TokenManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			response.Fail(c, 401, "UNAUTHORIZED", "Authorization header is required", nil)
			c.Abort()
			return
		}

		parts := strings.Fields(authHeader)

		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			response.Fail(c, 401, "UNAUTHORIZED", "Invalid authorization header format", nil)
			c.Abort()
			return
		}

		claims, err := tokenManager.ParseAccessToken(parts[1])
		if err != nil {
			response.Fail(c, 401, "INVALID_TOKEN", "Invalid or expired access token", nil)
			c.Abort()
			return
		}

		userID, err := strconv.ParseInt(claims.Subject, 10, 64)
		if err != nil || userID <= 0 {
			response.Fail(c, 401, "INVALID_TOKEN", "Invalid token subject", nil)
			c.Abort()
			return
		}

		if claims.Role != "customer" && claims.Role != "admin" {
			response.Fail(c, 401, "INVALID_TOKEN", "Invalid token role", nil)
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		c.Set("user_role", claims.Role)

		c.Next()
	}
}