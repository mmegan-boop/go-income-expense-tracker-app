package controller

import (
	"errors"
	"go-income-expense-tracker-app/internal/dto"
	"go-income-expense-tracker-app/internal/model"
	"go-income-expense-tracker-app/internal/service"
	"net/http"

	"github.com/labstack/echo/v5"
)

type AuthController struct {
	authService service.AuthService
}

func NewAuthController(authService service.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

func (c *AuthController) Register(ctx *echo.Context) error {
	var req dto.RegisterRequest

	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, dto.Response[any]{
			Status:  http.StatusBadRequest,
			Message: "invalid request body",
		})
	}

	user, err := c.authService.Register(req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrEmailExists):
			return ctx.JSON(http.StatusConflict, dto.Response[any]{
				Status:  http.StatusConflict,
				Message: err.Error(),
			})
		case errors.Is(err, service.ErrPasswordTooShort):
			return ctx.JSON(http.StatusBadRequest, dto.Response[any]{
				Status:  http.StatusBadRequest,
				Message: err.Error(),
			})
		default:
			return ctx.JSON(http.StatusInternalServerError, dto.Response[any]{
				Status:  http.StatusInternalServerError,
				Message: "failed to register user",
			})
		}
	}

	return ctx.JSON(http.StatusCreated, dto.Response[*model.User]{
		Status:  http.StatusCreated,
		Message: "user registered successfully",
		Data:    user,
	})
}

func (c *AuthController) Login(ctx *echo.Context) error {
	var req dto.LoginRequest

	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, dto.Response[any]{
			Status:  http.StatusBadRequest,
			Message: "invalid request body",
		})
	}

	token, err := c.authService.Login(req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			return ctx.JSON(http.StatusUnauthorized, dto.Response[any]{
				Status:  http.StatusUnauthorized,
				Message: err.Error(),
			})
		}

		return ctx.JSON(http.StatusInternalServerError, dto.Response[any]{
			Status:  http.StatusInternalServerError,
			Message: "failed to login",
		})
	}

	return ctx.JSON(http.StatusOK, dto.Response[any]{
		Status:  http.StatusOK,
		Message: "login successful",
		Data: map[string]string{
			"token": token,
		},
	})
}

func (c *AuthController) Logout(ctx *echo.Context) error {
	return ctx.JSON(http.StatusOK, dto.Response[any]{
		Status:  http.StatusOK,
		Message: "logged out successfully",
	})
}
