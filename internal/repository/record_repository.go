package repository

import (
	"go-income-expense-tracker-app/internal/model"

	"gorm.io/gorm"
)

type RecordRepository interface {
	Create(record *model.Record) error
	FindByID(id int) (*model.Record, error)
	FindAllByUserID(userID uint) ([]model.Record, error)
	Update(record *model.Record) error
	Delete(id int) error
}

type recordRepository struct {
	db *gorm.DB
}

func NewRecordRepository(db *gorm.DB) RecordRepository {
	return &recordRepository{db: db}
}

func (r *recordRepository) Create(record *model.Record) error {
	return r.db.Create(record).Error
}

func (r *recordRepository) FindByID(id int) (*model.Record, error) {
	var record model.Record

	if err := r.db.First(&record, id).Error; err != nil {
		return nil, err
	}

	return &record, nil
}

func (r *recordRepository) FindAllByUserID(userID uint) ([]model.Record, error) {
	var records []model.Record

	if err := r.db.Where("user_id = ?", userID).Find(&records).Error; err != nil {
		return nil, err
	}

	return records, nil
}

func (r *recordRepository) Update(record *model.Record) error {
	return r.db.Save(record).Error
}

func (r *recordRepository) Delete(id int) error {
	return r.db.Delete(&model.Record{}, id).Error
}
