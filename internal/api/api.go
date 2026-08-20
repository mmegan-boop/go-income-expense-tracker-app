package api

import (
	"go-income-expense-tracker-app/internal/controller"
	appmiddleware "go-income-expense-tracker-app/internal/middleware"
	"go-income-expense-tracker-app/internal/repository"
	"go-income-expense-tracker-app/internal/router"
	"go-income-expense-tracker-app/internal/service"
	appvalidator "go-income-expense-tracker-app/internal/validator"

	echojwt "github.com/labstack/echo-jwt/v5"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"gorm.io/gorm"
)

// NewEcho creates and configures the Echo HTTP server,
// including middleware, repositories, services, controllers, and registers API routes via the router package.
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
	userService := service.NewUserService(userRepository, authRepository)
	categoryService := service.NewCategoryService(categoryRepository)
	recordService := service.NewRecordService(recordRepository, categoryRepository)

	// Initialize controllers responsible for handling HTTP requests.
	authController := controller.NewAuthController(authService)
	userController := controller.NewUserController(userService)
	categoryController := controller.NewCategoryController(categoryService)
	recordController := controller.NewRecordController(recordService)

	// Create the base API route group.
	apiGroup := e.Group("/api")

	// Create a protected route group that requires JWT authentication.
	// Validate token and stores it in Echo's context
	protectedGroup := apiGroup.Group("", echojwt.WithConfig(jwtConfig.Init()), appmiddleware.VerifyToken)

	// Register routes for each domain.
	router.RegisterAuthRoutes(apiGroup, protectedGroup, authController)
	router.RegisterUserRoutes(protectedGroup, userController)
	router.RegisterCategoryRoutes(protectedGroup, categoryController)
	router.RegisterRecordRoutes(protectedGroup, recordController)

	return e
}
