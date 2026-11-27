package handler

import (
	"errors"
	"strconv"
	"strings"

	"gomerce/internal/model"
	"gomerce/internal/response"
	"gomerce/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type CategoryHandler struct {
	service *service.CategoryService
}

func NewCategoryHandler(categoryService *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		service: categoryService,
	}
}

//Request DTO
type CreateCategoryRequest struct {
	Name string `json:"name" binding:"required"`
}

//Method to create category
func (h *CategoryHandler) Create(c *gin.Context) {
	var req CreateCategoryRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "VALIDATION_ERROR", "Invalid Request Body", nil)
		return
	}

	req.Name = strings.TrimSpace(req.Name)

	if req.Name == "" {
		response.Fail(c, 400, "VALIDATION_ERROR", "Category Name Cannot be Empty", nil)
		return
	}

	category := &model.Category{
		Name: req.Name,
	}

	if err := h.service.Create(c.Request.Context(), category); err != nil {
		response.Fail(c, 500, "INTERNAL_SERVER_ERROR", "Failed to Create Category, perhaps already exist?", nil)
		return
	}

	response.Success(c, 201, category)
}

//Method to get all category
func (h *CategoryHandler) GetAll(c *gin.Context) {
	categories, err := h.service.GetAll(c.Request.Context())

	if err != nil {
		response.Fail(c, 500, "INTERNAL_SERVER_ERROR", "Failed to Retrieve Categories", nil)
		return
	}

	response.Success(c, 200, categories)
}

func isDuplicateKeyError(err error) bool {
	var pgErr *pgconn.PgError

	return errors.As(err, &pgErr) && pgErr.Code == "23505" //Unique violation code
}

//Method to update categories
func (h *CategoryHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil || id <= 0 {
		response.Fail(c, 400, "INVALID_ID", "Invalid Category ID", nil)
		return
	}

	var req CreateCategoryRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "VALIDATION_ERROR", "Invalid Request Body", nil)
		return
	}

	req.Name = strings.TrimSpace(req.Name)

	if req.Name == "" {
		response.Fail(c, 400, "VALIDATION_ERROR", "Category Name Cannot be Empty", nil)
		return
	}

	category := &model.Category{
		ID: id,
		Name: req.Name,
	}

	if err := h.service.Update(c.Request.Context(), category); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			response.Fail(c, 404, "NOT_FOUND", "Category Not Found", nil)
			return
		}

		if isDuplicateKeyError(err) {
			response.Fail(c, 409, "CATEGORY_EXIST", "Category Already Exists", nil)
			return
		}

		response.Fail(c, 500, "INTERNAL_SERVER_ERROR", "Failed to Update Category", nil)

		return
	}

	response.Success(c, 200, category)
}

//Method to delete category
func (h *CategoryHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil || id <= 0 {
		response.Fail(c, 400, "INVALID_ID", "Invalid Category ID", nil)
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			response.Fail(c, 404, "NOT_FOUND", "Category not Found", nil)
			return
		}

		response.Fail(c, 500, "INTERNAL_SERVER_ERROR", "Failed to Delete Category", nil)
		return
	}

	c.Status(204)
}