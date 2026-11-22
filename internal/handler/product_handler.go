package handler

import (
	"errors"

	"gomerce/internal/response"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

//Request Data Transfer Object (DTO) or a struct that will be used to receive or send data within API
//This struct make sure client can only send field that is allowed, not allowed example: (id, created at)
type CreateProductRequest struct {
	Name 	string	`json:"name" binding:"required"`
	Price	int64	`json:"price" binding:"gt=0"`
	Stock	int		`json:"stock" binding:"gte=0"`
}

type ProductHandler struct{}

func NewProductHandler() *ProductHandler {
	return &ProductHandler{}
}

//Function for client to create request, with validator so client couldn't send wrong data
func (h *ProductHandler) Create(c *gin.Context) {
	var request CreateProductRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		var validationErrors validator.ValidationErrors
		//checks if error from validation
		if errors.As(err, &validationErrors) {
			details := make(map[string]string)
			//if error from validation, checks which field is error
			for _, fieldErr := range validationErrors {
				field := fieldErr.Field()

				switch fieldErr.Tag() {
                case "required":
                    details[field] = "This field is required"
                case "gt":
                    details[field] = "Must be greater than zero"
                case "gte":
                    details[field] = "Must be zero or greater"
                }
			}
			//if the error comes from validation, this will be showed
			response.Fail(c, 400, "VALIDATION_ERROR", "Request Validation Failed", details)

			return
		}
		//if the error doesnt comes from validation, assume json is invalid (typos/anything else)
		response.Fail(c, 400, "INVALID_JSON", "Invalid Request Body", nil)
		
		return
	}
	//if success or nothing went wrong returns 201 status created
	response.Success(c, 201, request)
}