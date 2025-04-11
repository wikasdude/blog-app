package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/wikasdude/blog-backend/model"
	service "github.com/wikasdude/blog-backend/services"

	"github.com/wikasdude/blog-backend/utils"
)

type CategoryController struct {
	Service *service.CategoryService
}

func NewCategoryController(s *service.CategoryService) *CategoryController {
	return &CategoryController{Service: s}
}

// CreateCategory godoc
// @Summary      Create a new category
// @Description  Takes a JSON body and creates a new category.
// @Tags         categories
// @Accept       json
// @Produce      json
// @Param        category  body      model.CategoryRequest true  "Category to create"
// @Success      201       {object}  utils.APIResponse
// @Failure      400       {object}  utils.APIResponse
// @Failure      500       {object}  utils.APIResponse
// @Router       /api/categories [post]
func (c *CategoryController) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var category model.Category
	var err error
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		utils.SendError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	if category.Name == "" {
		utils.SendError(w, http.StatusBadRequest, "Category name cannot be empty", err)
		return
	}

	err = c.Service.CreateCategory(&category)
	if err != nil {

		utils.SendError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to create category %s", err), err)
		return
	}
	utils.SendSuccess(w, http.StatusCreated, "Category created successfully", category)
}

func (c *CategoryController) GetAllCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := c.Service.GetAllCategories()
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Failed to fetch categories", err)
		return
	}
	utils.SendSuccess(w, http.StatusOK, "Categories fetched successfully", categories)
}
