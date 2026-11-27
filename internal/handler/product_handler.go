package handler

import (
	"errors"
	"reflect"
	"strconv"
	"strings"

	"gomerce/internal/model"
	"gomerce/internal/response"
	"gomerce/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"
)

//Request Data Transfer Object (DTO) or a struct that will be used to receive or send data within API
//This struct make sure client can only send field that is allowed, not allowed example: (id, created at)
type CreateProductRequest struct {
	Name 	string	`json:"name" binding:"required"`
	Price	int64	`json:"price" binding:"gt=0"`
	Stock	int		`json:"stock" binding:"gte=0"`
}

//Same as above, just for updating.
type UpdateProductRequest struct {
	Name 	string	`json:"name" binding:"required"`
	Price	int64	`json:"price" binding:"gt=0"`
	Stock	int		`json:"stock" binding:"gte=0"`
}

type ProductHandler struct{
	service *service.ProductService
}

func NewProductHandler(productService *service.ProductService) *ProductHandler {
	return &ProductHandler{
		service: productService,
	}
}

//Function for client to create request, with validator so client couldn't send wrong data
func (h *ProductHandler) Create(c *gin.Context) {
	var request CreateProductRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		var validationErrors validator.ValidationErrors
		//checks if error from validation
		if errors.As(err, &validationErrors) {
			//This line do : read validation error, take field name from JSON tag, choose error message, return err
			details := validationDetails(request, validationErrors)
			
			response.Fail(c, 400, "VALIDATION_ERROR", "Request Validation Failed", details)

			return
		}
		//if the error doesnt comes from validation, assume json is invalid (typos/anything else)
		response.Fail(c, 400, "INVALID_JSON", "Invalid Request Body", nil)
		
		return
	}
	//Checks if name that is given is empty space or just space
	request.Name = strings.TrimSpace(request.Name)

	if request.Name == "" {
		response.Fail(c, 400, "VALIDATION_ERROR", "Request Validation Failed",
					  map[string]string{"name":"This field is required"})

		return
	}

	product := &model.Product{
		Name: request.Name,
		Price: request.Price,
		Stock: request.Stock,
	}

	if err := h.service.Create(c.Request.Context(), product); err != nil {
		response.Fail(c, 500, "INTERNAL_SERVER_ERROR", "Failed To Create Product, perhaps already exist?", nil)

		return
	}
	
	//if success or nothing went wrong returns 201 status created
	response.Success(c, 201, product)
}

//Helper function, converts  Go Struct field like name to json format. i.e (Name -> name or Price -> price)
func getJSONFieldName(request any, fieldName string) string {
	requestType := reflect.TypeOf(request)

	field, ok := requestType.FieldByName(fieldName)

	if !ok {
		return fieldName
	}

	jsonTag := field.Tag.Get("json")

	if jsonTag == "" {
		return fieldName
	} 

	return strings.Split(jsonTag, ",")[0]
}

//If error comes from validation, this function gives the details based on what went wrong
func validationDetails(request any, validationErrors validator.ValidationErrors) map[string]string {
	details := make(map[string]string)

	for _, fieldErr := range validationErrors {
		field := getJSONFieldName(request, fieldErr.Field())

		switch fieldErr.Tag() {
        case "required":
            details[field] = "This field is required"
        case "gt":
            details[field] = "Must be greater than zero"
        case "gte":
            details[field] = "Must be zero or greater"
        }
	}

	return details
}

func (h *ProductHandler) GetAll(c *gin.Context) {
	products, err := h.service.GetAll(c.Request.Context())

	if err != nil {
		response.Fail(c, 500, "INTENAL_SERVER_ERROR", "Failed To Retrieve Products", nil)
		
		return
	}

	response.Success(c, 200, products)
}

func (h *ProductHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil || id <= 0 {
		response.Fail(c, 400, "INVALID_PRODUCT_ID", "Product ID Must be A Positive Integer", nil)

		return
	}

	product, err := h.service.GetByID(c.Request.Context(), id)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			response.Fail(c, 404, "PRODUCT_NOT_FOUND", "Product Not Found", nil)

			return
		}

		response.Fail(c, 500, "INTERNAL_SERVER_ERROR", "Failed To Retrieve Product", nil)

		return
	}

	response.Success(c, 200, product)
}

func (h *ProductHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil || id <= 0 {
		response.Fail(c, 400, "INVALID_PRODUCT_ID", "Product ID Must be A Positive Integer", nil)

		return
	}

	var request UpdateProductRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		var validationErrors validator.ValidationErrors

		if errors.As(err, &validationErrors) {
			details := validationDetails(request, validationErrors)

			response.Fail(c, 400, "VALIDATION_ERROR", "Request Validation Failed", details)

			return
		}

		response.Fail(c, 400, "INVALID_JSON", "Invalid Body Request", nil)

		return
	}

	request.Name = strings.TrimSpace(request.Name)

	if request.Name == "" {
		response.Fail(c, 400, "VALIDATION_ERROR", "Request Validation Error", map[string]string{"name":"Field required"})

		return
	}

	product := &model.Product{
		ID:    id,
        Name:  request.Name,
        Price: request.Price,
        Stock: request.Stock,
	}

	if err := h.service.Update(c.Request.Context(), product); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			response.Fail(c, 404, "PRODUCT_NOT_FOUND", "Product Not Found", nil)

			return
		}

		response.Fail(c, 500, "INTERNAL_SERVER_ERROR", "Failed To Update Product", nil)

		return
	}

	response.Success(c, 200, gin.H{"message": "Product updated successfully"})
}

func (h *ProductHandler) Delete(c *gin.Context){
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil || id <= 0 {
		response.Fail(c, 400, "INVALID_PRODUCT_ID", "Product ID Must be A Positive Integer", nil)

		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			response.Fail(c, 404, "PRODUCT_NOT_FOUND", "Product Not Found", nil)

			return
		}

		response.Fail(c, 500, "INTERNAL_SERVER_ERROR", "Failed To Delete Product", nil)

		return
	}

	c.Status(204)
}