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

func TestAuthController_Register(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		e := newEcho()
		mockAuth := new(MockAuthService)
		ctrl := controller.NewAuthController(mockAuth)

		user := &model.User{ID: 1, Username: "john", Email: "john@example.com", Role: model.RoleUser}
		mockAuth.On("Register", dto.RegisterRequest{
			Username: "john",
			Email:    "john@example.com",
			Password: "password123",
		}).Return(user, nil)

		req := createRequest(http.MethodPost, "/api/auth/register", dto.RegisterRequest{
			Username: "john",
			Email:    "john@example.com",
			Password: "password123",
		})
		ctx, rec := createContextWithJWT(e, req, 0, "")

		err := ctrl.Register(ctx)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Contains(t, rec.Body.String(), "user registered successfully")
	})

	t.Run("invalid request body", func(t *testing.T) {
		e := newEcho()
		mockAuth := new(MockAuthService)
		ctrl := controller.NewAuthController(mockAuth)

		req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBufferString("invalid json"))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		ctx, rec := createContextWithJWT(e, req, 0, "")

		err := ctrl.Register(ctx)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "invalid request body")
	})

	t.Run("email already exists", func(t *testing.T) {
		e := newEcho()
		mockAuth := new(MockAuthService)
		ctrl := controller.NewAuthController(mockAuth)

		mockAuth.On("Register", dto.RegisterRequest{
			Username: "john",
			Email:    "john@example.com",
			Password: "password123",
		}).Return(nil, service.ErrEmailExists)

		req := createRequest(http.MethodPost, "/api/auth/register", dto.RegisterRequest{
			Username: "john",
			Email:    "john@example.com",
			Password: "password123",
		})
		ctx, rec := createContextWithJWT(e, req, 0, "")

		err := ctrl.Register(ctx)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusConflict, rec.Code)
		assert.Contains(t, rec.Body.String(), "email already registered")
	})

	t.Run("password too short (validation)", func(t *testing.T) {
		e := newEcho()
		mockAuth := new(MockAuthService)
		ctrl := controller.NewAuthController(mockAuth)

		req := createRequest(http.MethodPost, "/api/auth/register", dto.RegisterRequest{
			Username: "john",
			Email:    "john@example.com",
			Password: "123",
		})
		ctx, rec := createContextWithJWT(e, req, 0, "")

		err := ctrl.Register(ctx)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "min")
	})

	t.Run("password too short (service)", func(t *testing.T) {
		e := newEcho()
		mockAuth := new(MockAuthService)
		ctrl := controller.NewAuthController(mockAuth)

		mockAuth.On("Register", dto.RegisterRequest{
			Username: "john",
			Email:    "john@example.com",
			Password: "abcdef",
		}).Return(nil, service.ErrPasswordTooShort)

		req := createRequest(http.MethodPost, "/api/auth/register", dto.RegisterRequest{
			Username: "john",
			Email:    "john@example.com",
			Password: "abcdef",
		})
		ctx, rec := createContextWithJWT(e, req, 0, "")

		err := ctrl.Register(ctx)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "password must be at least 6 characters")
	})

	t.Run("generic error", func(t *testing.T) {
		e := newEcho()
		mockAuth := new(MockAuthService)
		ctrl := controller.NewAuthController(mockAuth)

		mockAuth.On("Register", dto.RegisterRequest{
			Username: "john",
			Email:    "john@example.com",
			Password: "password123",
		}).Return(nil, fmt.Errorf("db error"))

		req := createRequest(http.MethodPost, "/api/auth/register", dto.RegisterRequest{
			Username: "john",
			Email:    "john@example.com",
			Password: "password123",
		})
		ctx, rec := createContextWithJWT(e, req, 0, "")

		err := ctrl.Register(ctx)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), "failed to register user")
	})
}

func TestAuthController_Login(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		e := newEcho()
		mockAuth := new(MockAuthService)
		ctrl := controller.NewAuthController(mockAuth)

		mockAuth.On("Login", dto.LoginRequest{
			Email:    "john@example.com",
			Password: "password123",
		}).Return("jwt-token-abc", nil)

		req := createRequest(http.MethodPost, "/api/auth/login", dto.LoginRequest{
			Email:    "john@example.com",
			Password: "password123",
		})
		ctx, rec := createContextWithJWT(e, req, 0, "")

		err := ctrl.Login(ctx)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "login successful")
		assert.Contains(t, rec.Body.String(), "jwt-token-abc")
	})

	t.Run("invalid credentials", func(t *testing.T) {
		e := newEcho()
		mockAuth := new(MockAuthService)
		ctrl := controller.NewAuthController(mockAuth)

		mockAuth.On("Login", dto.LoginRequest{
			Email:    "wrong@example.com",
			Password: "wrongpassword",
		}).Return("", service.ErrInvalidCredentials)

		req := createRequest(http.MethodPost, "/api/auth/login", dto.LoginRequest{
			Email:    "wrong@example.com",
			Password: "wrongpassword",
		})
		ctx, rec := createContextWithJWT(e, req, 0, "")

		err := ctrl.Login(ctx)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.Contains(t, rec.Body.String(), "invalid email or password")
	})

	t.Run("generic error", func(t *testing.T) {
		e := newEcho()
		mockAuth := new(MockAuthService)
		ctrl := controller.NewAuthController(mockAuth)

		mockAuth.On("Login", dto.LoginRequest{
			Email:    "john@example.com",
			Password: "password123",
		}).Return("", fmt.Errorf("unexpected"))

		req := createRequest(http.MethodPost, "/api/auth/login", dto.LoginRequest{
			Email:    "john@example.com",
			Password: "password123",
		})
		ctx, rec := createContextWithJWT(e, req, 0, "")

		err := ctrl.Login(ctx)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), "failed to login")
	})
}

func TestAuthController_Logout(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		e := newEcho()
		mockAuth := new(MockAuthService)
		ctrl := controller.NewAuthController(mockAuth)

		req := createRequest(http.MethodPost, "/api/auth/logout", nil)
		ctx, rec := createContextWithJWT(e, req, 1, model.RoleUser)

		err := ctrl.Logout(ctx)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "logged out successfully")
	})
}
