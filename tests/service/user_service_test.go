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

type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) FindByID(id int) (*model.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepo) FindAll() ([]model.User, error) {
	args := m.Called()
	return args.Get(0).([]model.User), args.Error(1)
}

func (m *MockUserRepo) Update(user *model.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepo) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

type MockAuthRepo struct {
	mock.Mock
}

func (m *MockAuthRepo) Create(user *model.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockAuthRepo) FindByEmail(email string) (*model.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func newUserService() (*MockUserRepo, *MockAuthRepo, service.UserService) {
	mockUserRepo := new(MockUserRepo)
	mockAuthRepo := new(MockAuthRepo)
	svc := service.NewUserService(mockUserRepo, mockAuthRepo)
	return mockUserRepo, mockAuthRepo, svc
}

var testUser = &model.User{
	ID:       1,
	Username: "john",
	Email:    "john@example.com",
	Password: "hashedpassword",
	Role:     model.RoleUser,
}

func TestUserService_GetByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockUserRepo, _, svc := newUserService()

		mockUserRepo.On("FindByID", 1).Return(testUser, nil)

		user, err := svc.GetByID(1)
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, uint(1), user.ID)
		assert.Equal(t, "john", user.Username)
		assert.Equal(t, "john@example.com", user.Email)
	})

	t.Run("not found", func(t *testing.T) {
		mockUserRepo, _, svc := newUserService()

		mockUserRepo.On("FindByID", 999).Return(nil, fmt.Errorf("not found"))

		user, err := svc.GetByID(999)
		assert.ErrorIs(t, err, service.ErrUserNotFound)
		assert.Nil(t, user)
	})
}

func TestUserService_GetAll(t *testing.T) {
	t.Run("success with multiple users", func(t *testing.T) {
		mockUserRepo, _, svc := newUserService()

		users := []model.User{
			{ID: 1, Username: "john", Email: "john@example.com"},
			{ID: 2, Username: "jane", Email: "jane@example.com"},
		}
		mockUserRepo.On("FindAll").Return(users, nil)

		result, err := svc.GetAll()
		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, "john", result[0].Username)
		assert.Equal(t, "jane", result[1].Username)
	})

	t.Run("empty result", func(t *testing.T) {
		mockUserRepo, _, svc := newUserService()

		mockUserRepo.On("FindAll").Return([]model.User{}, nil)

		result, err := svc.GetAll()
		assert.NoError(t, err)
		assert.Empty(t, result)
	})

	t.Run("repository error", func(t *testing.T) {
		mockUserRepo, _, svc := newUserService()

		mockUserRepo.On("FindAll").Return([]model.User(nil), fmt.Errorf("db error"))

		result, err := svc.GetAll()
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "db error")
	})
}

func TestUserService_UpdateProfile(t *testing.T) {
	t.Run("success update username only", func(t *testing.T) {
		mockUserRepo, _, svc := newUserService()

		mockUserRepo.On("FindByID", 1).Return(testUser, nil)
		mockUserRepo.On("Update", mock.AnythingOfType("*model.User")).Return(nil)

		user, err := svc.UpdateProfile(1, dto.UpdateUserRequest{Username: "johnny"})
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, "johnny", user.Username)
		assert.Equal(t, "john@example.com", user.Email)
	})

	t.Run("success update email only", func(t *testing.T) {
		mockUserRepo, mockAuthRepo, svc := newUserService()

		mockUserRepo.On("FindByID", 1).Return(testUser, nil)
		mockAuthRepo.On("FindByEmail", "new@example.com").Return(nil, fmt.Errorf("not found"))
		mockUserRepo.On("Update", mock.AnythingOfType("*model.User")).Return(nil)

		user, err := svc.UpdateProfile(1, dto.UpdateUserRequest{Email: "new@example.com"})
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, "new@example.com", user.Email)
	})

	t.Run("success update password only", func(t *testing.T) {
		mockUserRepo, _, svc := newUserService()

		mockUserRepo.On("FindByID", 1).Return(testUser, nil)
		mockUserRepo.On("Update", mock.AnythingOfType("*model.User")).Return(nil)

		user, err := svc.UpdateProfile(1, dto.UpdateUserRequest{Password: "newpassword123"})
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.NotEqual(t, "hashedpassword", user.Password)
	})

	t.Run("success update all fields", func(t *testing.T) {
		mockUserRepo, mockAuthRepo, svc := newUserService()

		mockUserRepo.On("FindByID", 1).Return(testUser, nil)
		mockAuthRepo.On("FindByEmail", "updated@example.com").Return(nil, fmt.Errorf("not found"))
		mockUserRepo.On("Update", mock.AnythingOfType("*model.User")).Return(nil)

		user, err := svc.UpdateProfile(1, dto.UpdateUserRequest{
			Username: "updated",
			Email:    "updated@example.com",
			Password: "newpassword123",
		})
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, "updated", user.Username)
		assert.Equal(t, "updated@example.com", user.Email)
	})

	t.Run("not found", func(t *testing.T) {
		mockUserRepo, _, svc := newUserService()

		mockUserRepo.On("FindByID", 999).Return(nil, fmt.Errorf("not found"))

		user, err := svc.UpdateProfile(999, dto.UpdateUserRequest{Username: "newname"})
		assert.ErrorIs(t, err, service.ErrUserNotFound)
		assert.Nil(t, user)
	})

	t.Run("email already exists by another user", func(t *testing.T) {
		mockUserRepo, mockAuthRepo, svc := newUserService()

		otherUser := &model.User{ID: 2, Email: "jane@example.com"}
		mockUserRepo.On("FindByID", 1).Return(testUser, nil)
		mockAuthRepo.On("FindByEmail", "jane@example.com").Return(otherUser, nil)

		user, err := svc.UpdateProfile(1, dto.UpdateUserRequest{Email: "jane@example.com"})
		assert.ErrorIs(t, err, service.ErrEmailExists)
		assert.Nil(t, user)
	})

	t.Run("email already exists by same user (self update allowed)", func(t *testing.T) {
		mockUserRepo, mockAuthRepo, svc := newUserService()

		mockUserRepo.On("FindByID", 1).Return(testUser, nil)
		mockAuthRepo.On("FindByEmail", "john@example.com").Return(testUser, nil)
		mockUserRepo.On("Update", mock.AnythingOfType("*model.User")).Return(nil)

		user, err := svc.UpdateProfile(1, dto.UpdateUserRequest{Email: "john@example.com"})
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, "john@example.com", user.Email)
	})

	t.Run("invalid email format", func(t *testing.T) {
		mockUserRepo, _, svc := newUserService()

		mockUserRepo.On("FindByID", 1).Return(testUser, nil)

		user, err := svc.UpdateProfile(1, dto.UpdateUserRequest{Email: "not-an-email"})
		assert.ErrorIs(t, err, service.ErrInvalidEmailFormat)
		assert.Nil(t, user)
	})

	t.Run("password too short", func(t *testing.T) {
		mockUserRepo, _, svc := newUserService()

		mockUserRepo.On("FindByID", 1).Return(testUser, nil)

		user, err := svc.UpdateProfile(1, dto.UpdateUserRequest{Password: "123"})
		assert.ErrorIs(t, err, service.ErrPasswordTooShort)
		assert.Nil(t, user)
	})

	t.Run("repository update error", func(t *testing.T) {
		mockUserRepo, _, svc := newUserService()

		mockUserRepo.On("FindByID", 1).Return(testUser, nil)
		mockUserRepo.On("Update", mock.AnythingOfType("*model.User")).
			Return(fmt.Errorf("update failed"))

		user, err := svc.UpdateProfile(1, dto.UpdateUserRequest{Username: "newname"})
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "update failed")
	})
}

func TestUserService_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockUserRepo, _, svc := newUserService()

		mockUserRepo.On("FindByID", 1).Return(testUser, nil)
		mockUserRepo.On("Delete", 1).Return(nil)

		err := svc.Delete(1)
		assert.NoError(t, err)
	})

	t.Run("not found", func(t *testing.T) {
		mockUserRepo, _, svc := newUserService()

		mockUserRepo.On("FindByID", 999).Return(nil, fmt.Errorf("not found"))

		err := svc.Delete(999)
		assert.ErrorIs(t, err, service.ErrUserNotFound)
	})

	t.Run("repository delete error", func(t *testing.T) {
		mockUserRepo, _, svc := newUserService()

		mockUserRepo.On("FindByID", 1).Return(testUser, nil)
		mockUserRepo.On("Delete", 1).Return(fmt.Errorf("delete failed"))

		err := svc.Delete(1)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "delete failed")
	})
}
