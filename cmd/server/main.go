package main

import (
	"go-income-expense-tracker-app/internal/constant"
	"go-income-expense-tracker-app/internal/db"
	"go-income-expense-tracker-app/internal/utils"
)

func main() {
	// fmt.Println("Hello, Go!")
	dbConfig := db.DBConfig{
		Username: utils.GetConfig(constant.DB_USERNAME),
		Password: utils.GetConfig(constant.DB_PASSWORD),
		Database: utils.GetConfig(constant.DB_NAME),
		Host:     utils.GetConfig(constant.DB_HOST),
		Port:     utils.GetConfig(constant.DB_PORT),
	}

	var (
		repository = dbConfig.InitDB()
		// e          = api.NewEcho(repository, jwtConfig)
	)

	db.MigrateDB(repository)

	// appPort := fmt.Sprintf(":%s", utils.GetConfig(constant.PORT))

	// if err := e.Start(appPort); err != nil {
	// 	e.Logger.Error("failed to start server", "error", err)
	// }
}
