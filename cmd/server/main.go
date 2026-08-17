package main

import (
	"fmt"
	"go-income-expense-tracker-app/internal/api"
	"go-income-expense-tracker-app/internal/config"
	"go-income-expense-tracker-app/internal/constant"
	"go-income-expense-tracker-app/internal/db"
	"go-income-expense-tracker-app/internal/utils"
)

// main initializes the application configuration, database connection,
// JWT configuration, HTTP server, and starts the application.
func main() {
	// Load database configuration from environment variables.
	dbConfig := db.DBConfig{
		Username: utils.GetConfig(constant.DB_USERNAME),
		Password: utils.GetConfig(constant.DB_PASSWORD),
		Database: utils.GetConfig(constant.DB_NAME),
		Host:     utils.GetConfig(constant.DB_HOST),
		Port:     utils.GetConfig(constant.DB_PORT),
	}

	// Initialize the database repository and JWT configuration.
	repository := dbConfig.InitDB()
	jwtConfig := config.LoadJWTConfig()

	// Run database migrations when needed.
	// db.MigrateDB(repository)

	// Initialize the Echo HTTP server with the database and JWT configuration.
	e := api.NewEcho(repository, jwtConfig)

	// Get the application port from the environment configuration.
	appPort := fmt.Sprintf(":%s", utils.GetConfig(constant.PORT))

	// Start the HTTP server and log any startup error.
	if err := e.Start(appPort); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
