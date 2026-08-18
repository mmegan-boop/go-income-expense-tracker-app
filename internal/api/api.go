package api

import (
	"go-income-expense-tracker-app/internal/controller"
	appmiddleware "go-income-expense-tracker-app/internal/middleware"
	"go-income-expense-tracker-app/internal/repository"
	"go-income-expense-tracker-app/internal/service"
	appvalidator "go-income-expense-tracker-app/internal/validator"

	echojwt "github.com/labstack/echo-jwt/v5"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"gorm.io/gorm"
)

// NewEcho creates and configures the Echo HTTP server,
// including middleware, repositories, services, controllers, and API routes.
func NewEcho(db *gorm.DB, jwtConfig appmiddleware.JWTConfig) *echo.Echo {
	e := echo.New()

	// Register middleware for recovering from unexpected panics.
	e.Use(middleware.Recover())

	// Register custom validator for request validation.
	e.Validator = appvalidator.New()

	// Initialize repositories for database access.
	authRepository := repository.NewAuthRepository(db)
	userRepository := repository.NewUserRepository(db)
	categoryRepository := repository.NewCategoryRepository(db)
	recordRepository := repository.NewRecordRepository(db)

	// Initialize services containing application's business logic.
	authService := service.NewAuthService(authRepository, jwtConfig)
	userService := service.NewUserService(userRepository)
	categoryService := service.NewCategoryService(categoryRepository)
	recordService := service.NewRecordService(recordRepository, categoryRepository)

	// Initialize controllers responsible for handling HTTP requests.
	authController := controller.NewAuthController(authService)
	userController := controller.NewUserController(userService)
	categoryController := controller.NewCategoryController(categoryService)
	recordController := controller.NewRecordController(recordService)

	// Create the base API route group.
	apiGroup := e.Group("/api")

	// Register public authentication endpoints.
	apiGroup.POST("/auth/register", authController.Register)
	apiGroup.POST("/auth/login", authController.Login)

	// Create a protected route group that requires JWT authentication.
	// Validate token and stores it in Echo's context
	protectedGroup := apiGroup.Group("", echojwt.WithConfig(jwtConfig.Init()), appmiddleware.VerifyToken)

	// Register authentication endpoints.
	protectedGroup.POST("/auth/logout", authController.Logout)

	// Register user endpoints.
	protectedGroup.GET("/users/me", userController.GetProfile)
	protectedGroup.PUT("/users/me", userController.UpdateProfile)
	protectedGroup.GET("/users", userController.GetAll)
	protectedGroup.DELETE("/users/:id", userController.Delete)

	// Register categories endpoints.
	protectedGroup.POST("/categories", categoryController.Create)
	protectedGroup.GET("/categories", categoryController.GetAll)
	protectedGroup.GET("/categories/:id", categoryController.GetByID)
	protectedGroup.PUT("/categories/:id", categoryController.Update)
	protectedGroup.DELETE("/categories/:id", categoryController.Delete)

	// Register financial records endpoints.
	protectedGroup.POST("/records", recordController.Create)
	protectedGroup.GET("/records", recordController.GetAll)
	protectedGroup.GET("/records/:id", recordController.GetByID)
	protectedGroup.PUT("/records/:id", recordController.Update)
	protectedGroup.DELETE("/records/:id", recordController.Delete)

	return e
}
