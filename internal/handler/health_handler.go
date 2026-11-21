package handler

import (
	"github.com/gin-gonic/gin"
)

type HealthHandler struct {}
//constructor
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}
//HealthHandler struct method named Check, the function is to check wether the app is running ok
func (h *HealthHandler) Check(c *gin.Context) {
	c.JSON(200, gin.H{"message": "GoMerce is Running", "status": "ok"})
}