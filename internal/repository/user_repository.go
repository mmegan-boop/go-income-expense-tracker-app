package repository

import (
	"go-income-expense-tracker-app/internal/model"

	"gorm.io/gorm"
)

type UserRepository interface {
	FindByID(id int) (*model.User, error)
	FindAll() ([]model.User, error)
	Update(user *model.User) error
	Delete(id int) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindByID(id int) (*model.User, error) {
	var user model.User

	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) FindAll() ([]model.User, error) {
	var users []model.User

	if err := r.db.Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (r *userRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) Delete(id int) error {
	return r.db.Delete(&model.User{}, id).Error
}
