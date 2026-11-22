package response

import (
	"github.com/gin-gonic/gin"
)

//An API response in form of struct, so that the response is more structured and organized
type APIResponse struct {
	Success 	bool		`json:"success"`
	Data		any			`json:"data,omitempty"`
	Error		*APIError	`json:"error,omitempty"`
}
//An API error in form of struct, so that the response is more structured and organized
type APIError struct{
	Code	string				`json:"code"`
	Message	string				`json:"message"`
	Details map[string]string	`json:"details,omitempty"`
}

//Function if the response successful
func Success(c *gin.Context, status int, data any) {
	c.JSON(status, APIResponse{Success: true, Data: data})
}

//Function if the response unsuccessful
func Fail(c *gin.Context, status int, code string, message string, details map[string]string) {
	c.JSON(status, APIResponse{
		Success: false,
		Error: &APIError{Code: code, Message: message, Details: details},
	})
}