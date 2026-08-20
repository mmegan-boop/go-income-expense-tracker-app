package router

import (
	"go-income-expense-tracker-app/internal/controller"
	"go-income-expense-tracker-app/internal/middleware"
	"go-income-expense-tracker-app/internal/model"

	"github.com/labstack/echo/v5"
)

func RegisterUserRoutes(protectedGroup *echo.Group, userController *controller.UserController) {
	protectedGroup.GET("/users/me", userController.GetProfile)
	protectedGroup.PUT("/users/me", userController.UpdateProfile)
	protectedGroup.GET("/users", userController.GetAll, middleware.RequireRole(model.RoleAdmin))
	protectedGroup.DELETE("/users/:id", userController.Delete, middleware.RequireRole(model.RoleAdmin))
}
