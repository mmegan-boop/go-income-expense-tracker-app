package router

import (
	"go-income-expense-tracker-app/internal/controller"

	"github.com/labstack/echo/v5"
)

func RegisterAuthRoutes(apiGroup, protectedGroup *echo.Group, authController *controller.AuthController) {
	apiGroup.POST("/auth/register", authController.Register)
	apiGroup.POST("/auth/login", authController.Login)
	protectedGroup.POST("/auth/logout", authController.Logout)
}
