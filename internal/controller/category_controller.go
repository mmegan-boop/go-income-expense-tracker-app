package controller

import (
	"errors"
	"go-income-expense-tracker-app/internal/dto"
	"go-income-expense-tracker-app/internal/model"
	"go-income-expense-tracker-app/internal/service"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v5"
)

type CategoryController struct {
	categoryService service.CategoryService
}

func NewCategoryController(categoryService service.CategoryService) *CategoryController {
	return &CategoryController{categoryService: categoryService}
}

func (c *CategoryController) Create(ctx *echo.Context) error {
	// Create a request DTO
	var req dto.CategoryRequest

	// Bind the request body
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, dto.Response[any]{
			Status:  http.StatusBadRequest,
			Message: "invalid request body",
		})
	}

	// Validate the request
	if err := ctx.Validate(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, dto.Response[any]{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	// Call the service
	category, err := c.categoryService.Create(req)

	if err != nil {
		// Check for a duplicate category
		if errors.Is(err, service.ErrCategoryNameExists) {
			return ctx.JSON(http.StatusConflict, dto.Response[any]{
				Status:  http.StatusConflict,
				Message: err.Error(),
			})
		}

		// Handle unexpected errors
		return ctx.JSON(http.StatusInternalServerError, dto.Response[any]{
			Status:  http.StatusInternalServerError,
			Message: "failed to create category",
		})
	}

	return ctx.JSON(http.StatusCreated, dto.Response[*model.Category]{
		Status:  http.StatusCreated,
		Message: "category created successfully",
		Data:    category,
	})
}

func (c *CategoryController) GetByID(ctx *echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, dto.Response[any]{
			Status:  http.StatusBadRequest,
			Message: "invalid category id",
		})
	}

	category, err := c.categoryService.GetByID(id)
	if err != nil {
		if errors.Is(err, service.ErrCategoryNotFound) {
			return ctx.JSON(http.StatusNotFound, dto.Response[any]{
				Status:  http.StatusNotFound,
				Message: err.Error(),
			})
		}

		return ctx.JSON(http.StatusInternalServerError, dto.Response[any]{
			Status:  http.StatusInternalServerError,
			Message: "failed to get category",
		})
	}

	return ctx.JSON(http.StatusOK, dto.Response[*model.Category]{
		Status:  http.StatusOK,
		Message: "category retrieved successfully",
		Data:    category,
	})
}

func (c *CategoryController) GetAll(ctx *echo.Context) error {
	categories, err := c.categoryService.GetAll()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, dto.Response[any]{
			Status:  http.StatusInternalServerError,
			Message: "failed to get categories",
		})
	}

	return ctx.JSON(http.StatusOK, dto.Response[[]model.Category]{
		Status:  http.StatusOK,
		Message: "categories retrieved successfully",
		Data:    categories,
	})
}

func (c *CategoryController) Update(ctx *echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, dto.Response[any]{
			Status:  http.StatusBadRequest,
			Message: "invalid category id",
		})
	}

	var req dto.CategoryRequest

	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, dto.Response[any]{
			Status:  http.StatusBadRequest,
			Message: "invalid request body",
		})
	}

	if err := ctx.Validate(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, dto.Response[any]{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	category, err := c.categoryService.Update(id, req)
	if err != nil {
		if errors.Is(err, service.ErrCategoryNotFound) {
			return ctx.JSON(http.StatusNotFound, dto.Response[any]{
				Status:  http.StatusNotFound,
				Message: err.Error(),
			})
		}

		// Check for a duplicate category
		if errors.Is(err, service.ErrCategoryNameExists) {
			return ctx.JSON(http.StatusConflict, dto.Response[any]{
				Status:  http.StatusConflict,
				Message: err.Error(),
			})
		}

		return ctx.JSON(http.StatusInternalServerError, dto.Response[any]{
			Status:  http.StatusInternalServerError,
			Message: "failed to update category",
		})
	}

	return ctx.JSON(http.StatusOK, dto.Response[*model.Category]{
		Status:  http.StatusOK,
		Message: "category updated successfully",
		Data:    category,
	})
}

func (c *CategoryController) Delete(ctx *echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, dto.Response[any]{
			Status:  http.StatusBadRequest,
			Message: "invalid category id",
		})
	}

	if err := c.categoryService.Delete(id); err != nil {
		if errors.Is(err, service.ErrCategoryNotFound) {
			return ctx.JSON(http.StatusNotFound, dto.Response[any]{
				Status:  http.StatusNotFound,
				Message: err.Error(),
			})
		}

		return ctx.JSON(http.StatusInternalServerError, dto.Response[any]{
			Status:  http.StatusInternalServerError,
			Message: "failed to delete category",
		})
	}

	return ctx.JSON(http.StatusOK, dto.Response[any]{
		Status:  http.StatusOK,
		Message: "category deleted successfully",
	})
}
