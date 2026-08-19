package controller_test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/assert"

	"go-income-expense-tracker-app/internal/controller"
	"go-income-expense-tracker-app/internal/dto"
	"go-income-expense-tracker-app/internal/model"
	"go-income-expense-tracker-app/internal/service"
)

func TestCategoryController_Create(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		e := newEcho()
		mockCat := new(MockCategoryService)
		ctrl := controller.NewCategoryController(mockCat)

		cat := &model.Category{ID: 1, Name: "Food"}
		mockCat.On("Create", dto.CategoryRequest{Name: "Food"}).Return(cat, nil)

		req := createRequest(http.MethodPost, "/api/categories", dto.CategoryRequest{Name: "Food"})
		ctx, rec := createContextWithJWT(e, req, 1, model.RoleAdmin)

		err := ctrl.Create(ctx)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Contains(t, rec.Body.String(), "category created successfully")
	})

	t.Run("invalid request body", func(t *testing.T) {
		e := newEcho()
		mockCat := new(MockCategoryService)
		ctrl := controller.NewCategoryController(mockCat)

		req := httptest.NewRequest(http.MethodPost, "/api/categories", bytes.NewBufferString("invalid"))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		ctx, rec := createContextWithJWT(e, req, 1, model.RoleAdmin)

		err := ctrl.Create(ctx)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "invalid request body")
	})

	t.Run("name already exists", func(t *testing.T) {
		e := newEcho()
		mockCat := new(MockCategoryService)
		ctrl := controller.NewCategoryController(mockCat)

		mockCat.On("Create", dto.CategoryRequest{Name: "Food"}).Return(nil, service.ErrCategoryNameExists)

		req := createRequest(http.MethodPost, "/api/categories", dto.CategoryRequest{Name: "Food"})
		ctx, rec := createContextWithJWT(e, req, 1, model.RoleAdmin)

		err := ctrl.Create(ctx)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusConflict, rec.Code)
		assert.Contains(t, rec.Body.String(), "category name already exists")
	})

	t.Run("generic error", func(t *testing.T) {
		e := newEcho()
		mockCat := new(MockCategoryService)
		ctrl := controller.NewCategoryController(mockCat)

		mockCat.On("Create", dto.CategoryRequest{Name: "Food"}).Return(nil, fmt.Errorf("db error"))

		req := createRequest(http.MethodPost, "/api/categories", dto.CategoryRequest{Name: "Food"})
		ctx, rec := createContextWithJWT(e, req, 1, model.RoleAdmin)

		err := ctrl.Create(ctx)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), "failed to create category")
	})
}

func TestCategoryController_GetByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		e := newEcho()
		mockCat := new(MockCategoryService)
		ctrl := controller.NewCategoryController(mockCat)

		cat := &model.Category{ID: 1, Name: "Food"}
		mockCat.On("GetByID", 1).Return(cat, nil)

		req := createRequest(http.MethodGet, "/api/categories/1", nil)
		ctx, rec := createContextWithPathParamsAndJWT(e, req, echo.PathValues{
			{Name: "id", Value: "1"},
		}, 1, model.RoleUser)

		err := ctrl.GetByID(ctx)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "category retrieved successfully")
	})

	t.Run("invalid category id", func(t *testing.T) {
		e := newEcho()
		mockCat := new(MockCategoryService)
		ctrl := controller.NewCategoryController(mockCat)

		req := createRequest(http.MethodGet, "/api/categories/abc", nil)
		ctx, rec := createContextWithPathParamsAndJWT(e, req, echo.PathValues{
			{Name: "id", Value: "abc"},
		}, 1, model.RoleUser)

		err := ctrl.GetByID(ctx)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "invalid category id")
	})

	t.Run("not found", func(t *testing.T) {
		e := newEcho()
		mockCat := new(MockCategoryService)
		ctrl := controller.NewCategoryController(mockCat)

		mockCat.On("GetByID", 999).Return(nil, service.ErrCategoryNotFound)

		req := createRequest(http.MethodGet, "/api/categories/999", nil)
		ctx, rec := createContextWithPathParamsAndJWT(e, req, echo.PathValues{
			{Name: "id", Value: "999"},
		}, 1, model.RoleUser)

		err := ctrl.GetByID(ctx)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Contains(t, rec.Body.String(), "category not found")
	})

	t.Run("generic error", func(t *testing.T) {
		e := newEcho()
		mockCat := new(MockCategoryService)
		ctrl := controller.NewCategoryController(mockCat)

		mockCat.On("GetByID", 1).Return(nil, fmt.Errorf("db error"))

		req := createRequest(http.MethodGet, "/api/categories/1", nil)
		ctx, rec := createContextWithPathParamsAndJWT(e, req, echo.PathValues{
			{Name: "id", Value: "1"},
		}, 1, model.RoleUser)

		err := ctrl.GetByID(ctx)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), "failed to get category")
	})
}

func TestCategoryController_GetAll(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		e := newEcho()
		mockCat := new(MockCategoryService)
		ctrl := controller.NewCategoryController(mockCat)

		categories := []model.Category{
			{ID: 1, Name: "Food"},
			{ID: 2, Name: "Transport"},
		}
		mockCat.On("GetAll").Return(categories, nil)

		req := createRequest(http.MethodGet, "/api/categories", nil)
		ctx, rec := createContextWithJWT(e, req, 1, model.RoleUser)

		err := ctrl.GetAll(ctx)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "categories retrieved successfully")
	})

	t.Run("generic error", func(t *testing.T) {
		e := newEcho()
		mockCat := new(MockCategoryService)
		ctrl := controller.NewCategoryController(mockCat)

		mockCat.On("GetAll").Return([]model.Category{}, fmt.Errorf("db error"))

		req := createRequest(http.MethodGet, "/api/categories", nil)
		ctx, rec := createContextWithJWT(e, req, 1, model.RoleUser)

		err := ctrl.GetAll(ctx)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), "failed to get categories")
	})
}

func TestCategoryController_Update(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		e := newEcho()
		mockCat := new(MockCategoryService)
		ctrl := controller.NewCategoryController(mockCat)

		cat := &model.Category{ID: 1, Name: "Groceries"}
		mockCat.On("Update", 1, dto.CategoryRequest{Name: "Groceries"}).Return(cat, nil)

		req := createRequest(http.MethodPut, "/api/categories/1", dto.CategoryRequest{Name: "Groceries"})
		ctx, rec := createContextWithPathParamsAndJWT(e, req, echo.PathValues{
			{Name: "id", Value: "1"},
		}, 1, model.RoleAdmin)

		err := ctrl.Update(ctx)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "category updated successfully")
	})

	t.Run("invalid category id", func(t *testing.T) {
		e := newEcho()
		mockCat := new(MockCategoryService)
		ctrl := controller.NewCategoryController(mockCat)

		req := createRequest(http.MethodPut, "/api/categories/abc", dto.CategoryRequest{Name: "Food"})
		ctx, rec := createContextWithPathParamsAndJWT(e, req, echo.PathValues{
			{Name: "id", Value: "abc"},
		}, 1, model.RoleAdmin)

		err := ctrl.Update(ctx)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "invalid category id")
	})

	t.Run("invalid request body", func(t *testing.T) {
		e := newEcho()
		mockCat := new(MockCategoryService)
		ctrl := controller.NewCategoryController(mockCat)

		req := httptest.NewRequest(http.MethodPut, "/api/categories/1", bytes.NewBufferString("invalid"))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		ctx, rec := createContextWithPathParamsAndJWT(e, req, echo.PathValues{
			{Name: "id", Value: "1"},
		}, 1, model.RoleAdmin)

		err := ctrl.Update(ctx)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "invalid request body")
	})

	t.Run("not found", func(t *testing.T) {
		e := newEcho()
		mockCat := new(MockCategoryService)
		ctrl := controller.NewCategoryController(mockCat)

		mockCat.On("Update", 999, dto.CategoryRequest{Name: "Food"}).Return(nil, service.ErrCategoryNotFound)

		req := createRequest(http.MethodPut, "/api/categories/999", dto.CategoryRequest{Name: "Food"})
		ctx, rec := createContextWithPathParamsAndJWT(e, req, echo.PathValues{
			{Name: "id", Value: "999"},
		}, 1, model.RoleAdmin)

		err := ctrl.Update(ctx)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Contains(t, rec.Body.String(), "category not found")
	})

	t.Run("name already exists", func(t *testing.T) {
		e := newEcho()
		mockCat := new(MockCategoryService)
		ctrl := controller.NewCategoryController(mockCat)

		mockCat.On("Update", 1, dto.CategoryRequest{Name: "Transport"}).Return(nil, service.ErrCategoryNameExists)

		req := createRequest(http.MethodPut, "/api/categories/1", dto.CategoryRequest{Name: "Transport"})
		ctx, rec := createContextWithPathParamsAndJWT(e, req, echo.PathValues{
			{Name: "id", Value: "1"},
		}, 1, model.RoleAdmin)

		err := ctrl.Update(ctx)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusConflict, rec.Code)
		assert.Contains(t, rec.Body.String(), "category name already exists")
	})

	t.Run("generic error", func(t *testing.T) {
		e := newEcho()
		mockCat := new(MockCategoryService)
		ctrl := controller.NewCategoryController(mockCat)

		mockCat.On("Update", 1, dto.CategoryRequest{Name: "Food"}).Return(nil, fmt.Errorf("db error"))

		req := createRequest(http.MethodPut, "/api/categories/1", dto.CategoryRequest{Name: "Food"})
		ctx, rec := createContextWithPathParamsAndJWT(e, req, echo.PathValues{
			{Name: "id", Value: "1"},
		}, 1, model.RoleAdmin)

		err := ctrl.Update(ctx)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), "failed to update category")
	})
}

func TestCategoryController_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		e := newEcho()
		mockCat := new(MockCategoryService)
		ctrl := controller.NewCategoryController(mockCat)

		mockCat.On("Delete", 1).Return(nil)

		req := createRequest(http.MethodDelete, "/api/categories/1", nil)
		ctx, rec := createContextWithPathParamsAndJWT(e, req, echo.PathValues{
			{Name: "id", Value: "1"},
		}, 1, model.RoleAdmin)

		err := ctrl.Delete(ctx)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "category deleted successfully")
	})

	t.Run("invalid category id", func(t *testing.T) {
		e := newEcho()
		mockCat := new(MockCategoryService)
		ctrl := controller.NewCategoryController(mockCat)

		req := createRequest(http.MethodDelete, "/api/categories/abc", nil)
		ctx, rec := createContextWithPathParamsAndJWT(e, req, echo.PathValues{
			{Name: "id", Value: "abc"},
		}, 1, model.RoleAdmin)

		err := ctrl.Delete(ctx)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "invalid category id")
	})

	t.Run("not found", func(t *testing.T) {
		e := newEcho()
		mockCat := new(MockCategoryService)
		ctrl := controller.NewCategoryController(mockCat)

		mockCat.On("Delete", 999).Return(service.ErrCategoryNotFound)

		req := createRequest(http.MethodDelete, "/api/categories/999", nil)
		ctx, rec := createContextWithPathParamsAndJWT(e, req, echo.PathValues{
			{Name: "id", Value: "999"},
		}, 1, model.RoleAdmin)

		err := ctrl.Delete(ctx)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Contains(t, rec.Body.String(), "category not found")
	})

	t.Run("generic error", func(t *testing.T) {
		e := newEcho()
		mockCat := new(MockCategoryService)
		ctrl := controller.NewCategoryController(mockCat)

		mockCat.On("Delete", 1).Return(fmt.Errorf("db error"))

		req := createRequest(http.MethodDelete, "/api/categories/1", nil)
		ctx, rec := createContextWithPathParamsAndJWT(e, req, echo.PathValues{
			{Name: "id", Value: "1"},
		}, 1, model.RoleAdmin)

		err := ctrl.Delete(ctx)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), "failed to delete category")
	})
}
