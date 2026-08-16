package model

import (
	"time"

	"gorm.io/gorm"
)

type Record struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	UserID      uint           `json:"user_id"`
	CategoryID  uint           `json:"category_id"`
	RecordType  RecordType     `json:"record_type" gorm:"type:record_type"`
	Amount      float64        `json:"amount" gorm:"not null"`
	Description string         `json:"description"`
	RecordDate  time.Time      `json:"record_date"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

type RecordType string

const (
	RecordTypeIncome  RecordType = "income"
	RecordTypeExpense RecordType = "expense"
)
