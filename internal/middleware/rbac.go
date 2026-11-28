package middleware

import (
	"gomerce/internal/response"

	"github.com/gin-gonic/gin"
)

func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole := c.GetString("user_role")

		if userRole == "" {
			response.Fail(c, 401, "UNAUTHORIZED", "Authentication is required", nil)
			return
		}

		for _, role := range allowedRoles {
			if userRole == role {
				c.Next()
				return
			}
		}

		response.Fail(c, 403, "FORBIDDEN", "Not permitted to access this resource", nil)
		c.Abort()
	}
}