package model

import (
	"time"

	"gorm.io/gorm"
)

// gorm = ORM (Object-Relational Mapper) for Go
type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Username  string         `json:"username"`
	Email     string         `json:"email" gorm:"unique"`
	Password  string         `json:"-"` // json:"-"" = dont expose password through api json
	Role      Role           `json:"role" gorm:"type:role"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)
