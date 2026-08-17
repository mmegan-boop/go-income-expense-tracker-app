package controller

import (
	"errors"
	"go-income-expense-tracker-app/internal/dto"
	"go-income-expense-tracker-app/internal/middleware"
	"go-income-expense-tracker-app/internal/model"
	"go-income-expense-tracker-app/internal/service"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v5"
)

type UserController struct {
	userService service.UserService
}

func NewUserController(userService service.UserService) *UserController {
	return &UserController{userService: userService}
}

func (c *UserController) GetProfile(ctx *echo.Context) error {
	userID, err := middleware.GetUserID(ctx.Request().Context())
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, dto.Response[any]{
			Status:  http.StatusUnauthorized,
			Message: "invalid token",
		})
	}

	user, err := c.userService.GetByID(userID)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			return ctx.JSON(http.StatusNotFound, dto.Response[any]{
				Status:  http.StatusNotFound,
				Message: err.Error(),
			})
		}

		return ctx.JSON(http.StatusInternalServerError, dto.Response[any]{
			Status:  http.StatusInternalServerError,
			Message: "failed to get user",
		})
	}

	return ctx.JSON(http.StatusOK, dto.Response[*model.User]{
		Status:  http.StatusOK,
		Message: "user profile retrieved successfully",
		Data:    user,
	})
}

func (c *UserController) UpdateProfile(ctx *echo.Context) error {
	userID, err := middleware.GetUserID(ctx.Request().Context())
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, dto.Response[any]{
			Status:  http.StatusUnauthorized,
			Message: "invalid token",
		})
	}

	var req dto.UpdateUserRequest

	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, dto.Response[any]{
			Status:  http.StatusBadRequest,
			Message: "invalid request body",
		})
	}

	user, err := c.userService.Update(userID, req)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			return ctx.JSON(http.StatusNotFound, dto.Response[any]{
				Status:  http.StatusNotFound,
				Message: err.Error(),
			})
		}

		return ctx.JSON(http.StatusInternalServerError, dto.Response[any]{
			Status:  http.StatusInternalServerError,
			Message: "failed to update user",
		})
	}

	return ctx.JSON(http.StatusOK, dto.Response[*model.User]{
		Status:  http.StatusOK,
		Message: "user profile updated successfully",
		Data:    user,
	})
}

func (c *UserController) GetAll(ctx *echo.Context) error {
	users, err := c.userService.GetAll()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, dto.Response[any]{
			Status:  http.StatusInternalServerError,
			Message: "failed to get users",
		})
	}

	return ctx.JSON(http.StatusOK, dto.Response[[]model.User]{
		Status:  http.StatusOK,
		Message: "users retrieved successfully",
		Data:    users,
	})
}

func (c *UserController) Delete(ctx *echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, dto.Response[any]{
			Status:  http.StatusBadRequest,
			Message: "invalid user id",
		})
	}

	if err := c.userService.Delete(id); err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			return ctx.JSON(http.StatusNotFound, dto.Response[any]{
				Status:  http.StatusNotFound,
				Message: err.Error(),
			})
		}

		return ctx.JSON(http.StatusInternalServerError, dto.Response[any]{
			Status:  http.StatusInternalServerError,
			Message: "failed to delete user",
		})
	}

	return ctx.JSON(http.StatusOK, dto.Response[any]{
		Status:  http.StatusOK,
		Message: "user deleted successfully",
	})
}
