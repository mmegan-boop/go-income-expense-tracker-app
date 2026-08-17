package service

import (
	"errors"
	"go-income-expense-tracker-app/internal/dto"
	"go-income-expense-tracker-app/internal/model"
	"go-income-expense-tracker-app/internal/repository"
)

var (
	ErrCategoryNotFound = errors.New("category not found")
)

type CategoryService interface {
	Create(req dto.CategoryRequest) (*model.Category, error)
	GetByID(id int) (*model.Category, error)
	GetAll() ([]model.Category, error)
	Update(id int, req dto.CategoryRequest) (*model.Category, error)
	Delete(id int) error
}

type categoryService struct {
	categoryRepository repository.CategoryRepository
}

func NewCategoryService(categoryRepository repository.CategoryRepository) CategoryService {
	return &categoryService{categoryRepository: categoryRepository}
}

func (s *categoryService) Create(req dto.CategoryRequest) (*model.Category, error) {
	category := &model.Category{
		Name: req.Name,
	}

	if err := s.categoryRepository.Create(category); err != nil {
		return nil, err
	}

	return category, nil
}

func (s *categoryService) GetByID(id int) (*model.Category, error) {
	category, err := s.categoryRepository.FindByID(id)
	if err != nil {
		return nil, ErrCategoryNotFound
	}

	return category, nil
}

func (s *categoryService) GetAll() ([]model.Category, error) {
	return s.categoryRepository.FindAll()
}

func (s *categoryService) Update(id int, req dto.CategoryRequest) (*model.Category, error) {
	category, err := s.categoryRepository.FindByID(id)
	if err != nil {
		return nil, ErrCategoryNotFound
	}

	category.Name = req.Name

	if err := s.categoryRepository.Update(category); err != nil {
		return nil, err
	}

	return category, nil
}

func (s *categoryService) Delete(id int) error {
	if _, err := s.categoryRepository.FindByID(id); err != nil {
		return ErrCategoryNotFound
	}

	return s.categoryRepository.Delete(id)
}
