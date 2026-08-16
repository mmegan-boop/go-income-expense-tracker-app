package db

import (
	"fmt"
	"go-income-expense-tracker-app/internal/model"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBConfig struct {
	Username string
	Password string
	Database string
	Host     string
	Port     string
}

func (config *DBConfig) InitDB() *gorm.DB {
	var err error

	var dsn string = fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s",
		config.Host,
		config.Username,
		config.Password,
		config.Database,
		config.Port,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("error when connecting to the database: %s\n", err)
	}

	log.Println("connected to the database")

	return db
}

func MigrateDB(db *gorm.DB) {
	db.Exec(`DO $$ BEGIN
		CREATE TYPE role AS ENUM ('admin', 'user');
	EXCEPTION WHEN duplicate_object THEN null; END $$;`)
	db.Exec(`DO $$ BEGIN
		CREATE TYPE record_type AS ENUM ('income', 'expense');
	EXCEPTION WHEN duplicate_object THEN null; END $$;`)

	err := db.AutoMigrate(
		&model.User{},
		&model.Category{},
		&model.Record{},
	)

	if err != nil {
		log.Fatalf("database migration failed: %v\n", err)
	}

	log.Println("database migration succeed")
}

func CloseDB(db *gorm.DB) error {
	database, err := db.DB()

	if err != nil {
		log.Printf("error when getting the database instance: %v", err)
		return err
	}

	if err := database.Close(); err != nil {
		log.Printf("error when closing the database connection: %v", err)
		return err
	}

	log.Println("database connection is closed")

	return nil
}
