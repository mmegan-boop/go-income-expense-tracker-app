package service_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"go-income-expense-tracker-app/internal/dto"
	"go-income-expense-tracker-app/internal/model"
	"go-income-expense-tracker-app/internal/service"
)

type MockCategoryRepo struct {
	mock.Mock
}

func (m *MockCategoryRepo) Create(category *model.Category) error {
	args := m.Called(category)
	return args.Error(0)
}

func (m *MockCategoryRepo) FindByID(id int) (*model.Category, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Category), args.Error(1)
}

func (m *MockCategoryRepo) FindByName(name string) (*model.Category, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Category), args.Error(1)
}

func (m *MockCategoryRepo) FindAll() ([]model.Category, error) {
	args := m.Called()
	return args.Get(0).([]model.Category), args.Error(1)
}

func (m *MockCategoryRepo) Update(category *model.Category) error {
	args := m.Called(category)
	return args.Error(0)
}

func (m *MockCategoryRepo) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func newCategoryService() (*MockCategoryRepo, service.CategoryService) {
	mockRepo := new(MockCategoryRepo)
	svc := service.NewCategoryService(mockRepo)
	return mockRepo, svc
}

var testCat = &model.Category{ID: 1, Name: "Food"}

func TestCategoryService_Create(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo, svc := newCategoryService()

		mockRepo.On("FindByName", "Food").Return(nil, fmt.Errorf("not found"))
		mockRepo.On("Create", mock.AnythingOfType("*model.Category")).Return(nil)

		cat, err := svc.Create(dto.CategoryRequest{Name: "Food"})
		assert.NoError(t, err)
		assert.NotNil(t, cat)
		assert.Equal(t, "Food", cat.Name)
		mockRepo.AssertExpectations(t)
	})

	t.Run("name already exists", func(t *testing.T) {
		mockRepo, svc := newCategoryService()

		mockRepo.On("FindByName", "Food").Return(testCat, nil)

		cat, err := svc.Create(dto.CategoryRequest{Name: "Food"})
		assert.ErrorIs(t, err, service.ErrCategoryNameExists)
		assert.Nil(t, cat)
	})

	t.Run("repository create error", func(t *testing.T) {
		mockRepo, svc := newCategoryService()

		mockRepo.On("FindByName", "Food").Return(nil, fmt.Errorf("not found"))
		mockRepo.On("Create", mock.AnythingOfType("*model.Category")).
			Return(fmt.Errorf("db error"))

		cat, err := svc.Create(dto.CategoryRequest{Name: "Food"})
		assert.Error(t, err)
		assert.Nil(t, cat)
		assert.Contains(t, err.Error(), "db error")
	})
}

func TestCategoryService_GetByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo, svc := newCategoryService()

		mockRepo.On("FindByID", 1).Return(testCat, nil)

		cat, err := svc.GetByID(1)
		assert.NoError(t, err)
		assert.NotNil(t, cat)
		assert.Equal(t, uint(1), cat.ID)
		assert.Equal(t, "Food", cat.Name)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo, svc := newCategoryService()

		mockRepo.On("FindByID", 999).Return(nil, fmt.Errorf("not found"))

		cat, err := svc.GetByID(999)
		assert.ErrorIs(t, err, service.ErrCategoryNotFound)
		assert.Nil(t, cat)
	})
}

func TestCategoryService_GetAll(t *testing.T) {
	t.Run("success with multiple categories", func(t *testing.T) {
		mockRepo, svc := newCategoryService()

		categories := []model.Category{
			{ID: 1, Name: "Food"},
			{ID: 2, Name: "Transport"},
			{ID: 3, Name: "Salary"},
		}
		mockRepo.On("FindAll").Return(categories, nil)

		result, err := svc.GetAll()
		assert.NoError(t, err)
		assert.Len(t, result, 3)
		assert.Equal(t, "Food", result[0].Name)
		assert.Equal(t, "Salary", result[2].Name)
	})

	t.Run("empty result", func(t *testing.T) {
		mockRepo, svc := newCategoryService()

		mockRepo.On("FindAll").Return([]model.Category{}, nil)

		result, err := svc.GetAll()
		assert.NoError(t, err)
		assert.Empty(t, result)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo, svc := newCategoryService()

		mockRepo.On("FindAll").Return([]model.Category(nil), fmt.Errorf("db error"))

		result, err := svc.GetAll()
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "db error")
	})
}

func TestCategoryService_Update(t *testing.T) {
	t.Run("success with different name", func(t *testing.T) {
		mockRepo, svc := newCategoryService()

		mockRepo.On("FindByID", 1).Return(testCat, nil)
		mockRepo.On("FindByName", "Groceries").Return(nil, fmt.Errorf("not found"))
		mockRepo.On("Update", mock.AnythingOfType("*model.Category")).Return(nil)

		cat, err := svc.Update(1, dto.CategoryRequest{Name: "Groceries"})
		assert.NoError(t, err)
		assert.NotNil(t, cat)
		assert.Equal(t, "Groceries", cat.Name)
		mockRepo.AssertExpectations(t)
	})

	t.Run("success with same name (self update)", func(t *testing.T) {
		mockRepo, svc := newCategoryService()

		mockRepo.On("FindByID", 1).Return(testCat, nil)
		mockRepo.On("FindByName", "Food").Return(testCat, nil)
		mockRepo.On("Update", mock.AnythingOfType("*model.Category")).Return(nil)

		cat, err := svc.Update(1, dto.CategoryRequest{Name: "Food"})
		assert.NoError(t, err)
		assert.NotNil(t, cat)
		assert.Equal(t, "Food", cat.Name)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo, svc := newCategoryService()

		mockRepo.On("FindByID", 999).Return(nil, fmt.Errorf("not found"))

		cat, err := svc.Update(999, dto.CategoryRequest{Name: "Food"})
		assert.ErrorIs(t, err, service.ErrCategoryNotFound)
		assert.Nil(t, cat)
	})

	t.Run("name already exists by another category", func(t *testing.T) {
		mockRepo, svc := newCategoryService()

		otherCat := &model.Category{ID: 2, Name: "Transport"}
		mockRepo.On("FindByID", 1).Return(testCat, nil)
		mockRepo.On("FindByName", "Transport").Return(otherCat, nil)

		cat, err := svc.Update(1, dto.CategoryRequest{Name: "Transport"})
		assert.ErrorIs(t, err, service.ErrCategoryNameExists)
		assert.Nil(t, cat)
	})

	t.Run("repository update error", func(t *testing.T) {
		mockRepo, svc := newCategoryService()

		mockRepo.On("FindByID", 1).Return(testCat, nil)
		mockRepo.On("FindByName", "Groceries").Return(nil, fmt.Errorf("not found"))
		mockRepo.On("Update", mock.AnythingOfType("*model.Category")).
			Return(fmt.Errorf("update failed"))

		cat, err := svc.Update(1, dto.CategoryRequest{Name: "Groceries"})
		assert.Error(t, err)
		assert.Nil(t, cat)
		assert.Contains(t, err.Error(), "update failed")
	})
}

func TestCategoryService_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo, svc := newCategoryService()

		mockRepo.On("FindByID", 1).Return(testCat, nil)
		mockRepo.On("Delete", 1).Return(nil)

		err := svc.Delete(1)
		assert.NoError(t, err)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo, svc := newCategoryService()

		mockRepo.On("FindByID", 999).Return(nil, fmt.Errorf("not found"))

		err := svc.Delete(999)
		assert.ErrorIs(t, err, service.ErrCategoryNotFound)
	})

	t.Run("repository delete error", func(t *testing.T) {
		mockRepo, svc := newCategoryService()

		mockRepo.On("FindByID", 1).Return(testCat, nil)
		mockRepo.On("Delete", 1).Return(fmt.Errorf("delete failed"))

		err := svc.Delete(1)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "delete failed")
	})
}
