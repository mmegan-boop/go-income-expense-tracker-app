package api

import (
	"go-income-expense-tracker-app/internal/controller"
	appmiddleware "go-income-expense-tracker-app/internal/middleware"
	"go-income-expense-tracker-app/internal/repository"
	"go-income-expense-tracker-app/internal/service"

	echojwt "github.com/labstack/echo-jwt/v5"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"gorm.io/gorm"
)

func NewEcho(db *gorm.DB, jwtConfig appmiddleware.JWTConfig) *echo.Echo {
	e := echo.New()

	e.Use(middleware.Recover())

	authRepository := repository.NewAuthRepository(db)
	userRepository := repository.NewUserRepository(db)
	categoryRepository := repository.NewCategoryRepository(db)
	recordRepository := repository.NewRecordRepository(db)

	authService := service.NewAuthService(authRepository, jwtConfig)
	userService := service.NewUserService(userRepository)
	categoryService := service.NewCategoryService(categoryRepository)
	recordService := service.NewRecordService(recordRepository)

	authController := controller.NewAuthController(authService)
	userController := controller.NewUserController(userService)
	categoryController := controller.NewCategoryController(categoryService)
	recordController := controller.NewRecordController(recordService)

	apiGroup := e.Group("/api")

	apiGroup.POST("/auth/register", authController.Register)
	apiGroup.POST("/auth/login", authController.Login)

	protectedGroup := apiGroup.Group("", echojwt.WithConfig(jwtConfig.Init()))

	protectedGroup.GET("/users/me", userController.GetProfile)
	protectedGroup.PUT("/users/me", userController.UpdateProfile)
	protectedGroup.GET("/users", userController.GetAll)
	protectedGroup.DELETE("/users/:id", userController.Delete)

	protectedGroup.POST("/categories", categoryController.Create)
	protectedGroup.GET("/categories", categoryController.GetAll)
	protectedGroup.GET("/categories/:id", categoryController.GetByID)
	protectedGroup.PUT("/categories/:id", categoryController.Update)
	protectedGroup.DELETE("/categories/:id", categoryController.Delete)

	protectedGroup.POST("/records", recordController.Create)
	protectedGroup.GET("/records", recordController.GetAll)
	protectedGroup.GET("/records/:id", recordController.GetByID)
	protectedGroup.PUT("/records/:id", recordController.Update)
	protectedGroup.DELETE("/records/:id", recordController.Delete)

	return e
}
