package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// This function returns gin.HandlerFunc type so other package can use it by router.Use() later
// This function will be used as custom/self-made middleware
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := uuid.NewString()
		//set at c / *gin.Context in order to make it possible that other middleware get the same RequestID
		c.Set("request_id", requestID)
		//send identifier to client, for error reporting/identification
		c.Header("X-Request-ID", requestID)
		c.Next()
	}
}