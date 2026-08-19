package main

import (
	"go-income-expense-tracker-app/internal/constant"
	"go-income-expense-tracker-app/internal/db"
	"go-income-expense-tracker-app/internal/model"
	"go-income-expense-tracker-app/internal/repository"
	"go-income-expense-tracker-app/internal/service"
	"go-income-expense-tracker-app/internal/utils"
	"log"

	"github.com/go-playground/validator/v10"
)

func main() {
	dbConfig := db.DBConfig{
		Username: utils.GetConfig(constant.DB_USERNAME),
		Password: utils.GetConfig(constant.DB_PASSWORD),
		Database: utils.GetConfig(constant.DB_NAME),
		Host:     utils.GetConfig(constant.DB_HOST),
		Port:     utils.GetConfig(constant.DB_PORT),
	}

	database := dbConfig.InitDB()

	// Run database migrations when needed.
	// db.MigrateDB(database)

	username := utils.GetConfig(constant.ADMIN_USERNAME)
	email := utils.GetConfig(constant.ADMIN_EMAIL)
	password := utils.GetConfig(constant.ADMIN_PASSWORD)

	if username == "" || email == "" || password == "" {
		log.Fatal("ADMIN_USERNAME, ADMIN_EMAIL, and ADMIN_PASSWORD must be set in .env")
	}

	v := validator.New()
	if err := v.Var(email, "email"); err != nil {
		log.Fatal(service.ErrInvalidEmailFormat)
	}

	if len(password) < 6 {
		log.Fatal(service.ErrPasswordTooShort)
	}

	authRepository := repository.NewAuthRepository(database)

	if _, err := authRepository.FindByEmail(email); err == nil {
		log.Fatal(service.ErrEmailExists)
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		log.Fatal(err)
	}

	admin := &model.User{
		Username: username,
		Email:    email,
		Password: hashedPassword,
		Role:     model.RoleAdmin,
	}

	if err := authRepository.Create(admin); err != nil {
		log.Fatalf("failed to create admin: %v\n", err)
	}

	log.Println("admin created successfully")
}
