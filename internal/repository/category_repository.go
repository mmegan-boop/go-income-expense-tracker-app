package repository

import (
	"go-income-expense-tracker-app/internal/model"

	"gorm.io/gorm"
)

type CategoryRepository interface {
	Create(category *model.Category) error
	FindByID(id int) (*model.Category, error)
	FindByName(name string) (*model.Category, error)
	FindAll() ([]model.Category, error)
	Update(category *model.Category) error
	Delete(id int) error
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) Create(category *model.Category) error {
	return r.db.Create(category).Error
}

func (r *categoryRepository) FindByID(id int) (*model.Category, error) {
	var category model.Category

	if err := r.db.First(&category, id).Error; err != nil {
		return nil, err
	}

	return &category, nil
}

func (r *categoryRepository) FindByName(name string) (*model.Category, error) {
	var category model.Category

	if err := r.db.Where("name = ?", name).First(&category).Error; err != nil {
		return nil, err
	}

	return &category, nil
}

func (r *categoryRepository) FindAll() ([]model.Category, error) {
	var categories []model.Category

	if err := r.db.Find(&categories).Error; err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *categoryRepository) Update(category *model.Category) error {
	return r.db.Save(category).Error
}

func (r *categoryRepository) Delete(id int) error {
	return r.db.Delete(&model.Category{}, id).Error
}
