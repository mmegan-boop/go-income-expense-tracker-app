package router

import (
	"go-income-expense-tracker-app/internal/controller"
	"go-income-expense-tracker-app/internal/middleware"
	"go-income-expense-tracker-app/internal/model"

	"github.com/labstack/echo/v5"
)

func RegisterCategoryRoutes(protectedGroup *echo.Group, categoryController *controller.CategoryController) {
	protectedGroup.POST("/categories", categoryController.Create, middleware.RequireRole(model.RoleAdmin))
	protectedGroup.GET("/categories", categoryController.GetAll)
	protectedGroup.GET("/categories/:id", categoryController.GetByID)
	protectedGroup.PUT("/categories/:id", categoryController.Update, middleware.RequireRole(model.RoleAdmin))
	protectedGroup.DELETE("/categories/:id", categoryController.Delete, middleware.RequireRole(model.RoleAdmin))
}
