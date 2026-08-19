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

func TestUserController_GetProfile(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		e := newEcho()
		mockUser := new(MockUserService)
		ctrl := controller.NewUserController(mockUser)

		user := &model.User{ID: 1, Username: "john", Email: "john@example.com"}
		mockUser.On("GetByID", 1).Return(user, nil)

		req := createRequest(http.MethodGet, "/api/users/me", nil)
		ctx, rec := createContextWithJWT(e, req, 1, model.RoleUser)

		err := ctrl.GetProfile(ctx)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "user profile retrieved successfully")
	})

	t.Run("user not found", func(t *testing.T) {
		e := newEcho()
		mockUser := new(MockUserService)
		ctrl := controller.NewUserController(mockUser)

		mockUser.On("GetByID", 999).Return(nil, service.ErrUserNotFound)

		req := createRequest(http.MethodGet, "/api/users/me", nil)
		ctx, rec := createContextWithJWT(e, req, 999, model.RoleUser)

		err := ctrl.GetProfile(ctx)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Contains(t, rec.Body.String(), "user not found")
	})

	t.Run("generic error", func(t *testing.T) {
		e := newEcho()
		mockUser := new(MockUserService)
		ctrl := controller.NewUserController(mockUser)

		mockUser.On("GetByID", 1).Return(nil, fmt.Errorf("db error"))

		req := createRequest(http.MethodGet, "/api/users/me", nil)
		ctx, rec := createContextWithJWT(e, req, 1, model.RoleUser)

		err := ctrl.GetProfile(ctx)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), "failed to get user")
	})
}

func TestUserController_UpdateProfile(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		e := newEcho()
		mockUser := new(MockUserService)
		ctrl := controller.NewUserController(mockUser)

		user := &model.User{ID: 1, Username: "johnny", Email: "john@example.com"}
		mockUser.On("UpdateProfile", 1, dto.UpdateUserRequest{Username: "johnny"}).Return(user, nil)

		req := createRequest(http.MethodPut, "/api/users/me", dto.UpdateUserRequest{Username: "johnny"})
		ctx, rec := createContextWithJWT(e, req, 1, model.RoleUser)

		err := ctrl.UpdateProfile(ctx)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "user profile updated successfully")
	})

	t.Run("invalid request body", func(t *testing.T) {
		e := newEcho()
		mockUser := new(MockUserService)
		ctrl := controller.NewUserController(mockUser)

		req := httptest.NewRequest(http.MethodPut, "/api/users/me", bytes.NewBufferString("invalid json"))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		ctx, rec := createContextWithJWT(e, req, 1, model.RoleUser)

		err := ctrl.UpdateProfile(ctx)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "invalid request body")
	})

	t.Run("user not found", func(t *testing.T) {
		e := newEcho()
		mockUser := new(MockUserService)
		ctrl := controller.NewUserController(mockUser)

		mockUser.On("UpdateProfile", 999, dto.UpdateUserRequest{Username: "new"}).Return(nil, service.ErrUserNotFound)

		req := createRequest(http.MethodPut, "/api/users/me", dto.UpdateUserRequest{Username: "new"})
		ctx, rec := createContextWithJWT(e, req, 999, model.RoleUser)

		err := ctrl.UpdateProfile(ctx)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Contains(t, rec.Body.String(), "user not found")
	})

	t.Run("email already exists", func(t *testing.T) {
		e := newEcho()
		mockUser := new(MockUserService)
		ctrl := controller.NewUserController(mockUser)

		mockUser.On("UpdateProfile", 1, dto.UpdateUserRequest{Email: "taken@example.com"}).
			Return(nil, service.ErrEmailExists)

		req := createRequest(http.MethodPut, "/api/users/me", dto.UpdateUserRequest{Email: "taken@example.com"})
		ctx, rec := createContextWithJWT(e, req, 1, model.RoleUser)

		err := ctrl.UpdateProfile(ctx)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusConflict, rec.Code)
		assert.Contains(t, rec.Body.String(), "email already registered")
	})

	t.Run("generic error", func(t *testing.T) {
		e := newEcho()
		mockUser := new(MockUserService)
		ctrl := controller.NewUserController(mockUser)

		mockUser.On("UpdateProfile", 1, dto.UpdateUserRequest{Username: "new"}).
			Return(nil, fmt.Errorf("db error"))

		req := createRequest(http.MethodPut, "/api/users/me", dto.UpdateUserRequest{Username: "new"})
		ctx, rec := createContextWithJWT(e, req, 1, model.RoleUser)

		err := ctrl.UpdateProfile(ctx)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), "failed to update user")
	})
}

func TestUserController_GetAll(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		e := newEcho()
		mockUser := new(MockUserService)
		ctrl := controller.NewUserController(mockUser)

		users := []model.User{
			{ID: 1, Username: "john"},
			{ID: 2, Username: "jane"},
		}
		mockUser.On("GetAll").Return(users, nil)

		req := createRequest(http.MethodGet, "/api/users", nil)
		ctx, rec := createContextWithJWT(e, req, 1, model.RoleAdmin)

		err := ctrl.GetAll(ctx)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "users retrieved successfully")
	})

	t.Run("generic error", func(t *testing.T) {
		e := newEcho()
		mockUser := new(MockUserService)
		ctrl := controller.NewUserController(mockUser)

		mockUser.On("GetAll").Return([]model.User{}, fmt.Errorf("db error"))

		req := createRequest(http.MethodGet, "/api/users", nil)
		ctx, rec := createContextWithJWT(e, req, 1, model.RoleAdmin)

		err := ctrl.GetAll(ctx)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), "failed to get users")
	})
}

func TestUserController_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		e := newEcho()
		mockUser := new(MockUserService)
		ctrl := controller.NewUserController(mockUser)

		mockUser.On("Delete", 2).Return(nil)

		req := createRequest(http.MethodDelete, "/api/users/2", nil)
		ctx, rec := createContextWithPathParamsAndJWT(e, req, echo.PathValues{
			{Name: "id", Value: "2"},
		}, 1, model.RoleAdmin)

		err := ctrl.Delete(ctx)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "user deleted successfully")
	})

	t.Run("invalid user id", func(t *testing.T) {
		e := newEcho()
		mockUser := new(MockUserService)
		ctrl := controller.NewUserController(mockUser)

		req := createRequest(http.MethodDelete, "/api/users/abc", nil)
		ctx, rec := createContextWithPathParamsAndJWT(e, req, echo.PathValues{
			{Name: "id", Value: "abc"},
		}, 1, model.RoleAdmin)

		err := ctrl.Delete(ctx)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "invalid user id")
	})

	t.Run("user not found", func(t *testing.T) {
		e := newEcho()
		mockUser := new(MockUserService)
		ctrl := controller.NewUserController(mockUser)

		mockUser.On("Delete", 999).Return(service.ErrUserNotFound)

		req := createRequest(http.MethodDelete, "/api/users/999", nil)
		ctx, rec := createContextWithPathParamsAndJWT(e, req, echo.PathValues{
			{Name: "id", Value: "999"},
		}, 1, model.RoleAdmin)

		err := ctrl.Delete(ctx)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Contains(t, rec.Body.String(), "user not found")
	})

	t.Run("generic error", func(t *testing.T) {
		e := newEcho()
		mockUser := new(MockUserService)
		ctrl := controller.NewUserController(mockUser)

		mockUser.On("Delete", 1).Return(fmt.Errorf("db error"))

		req := createRequest(http.MethodDelete, "/api/users/1", nil)
		ctx, rec := createContextWithPathParamsAndJWT(e, req, echo.PathValues{
			{Name: "id", Value: "1"},
		}, 1, model.RoleAdmin)

		err := ctrl.Delete(ctx)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), "failed to delete user")
	})
}
