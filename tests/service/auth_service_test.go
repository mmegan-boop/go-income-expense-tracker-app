package service_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"go-income-expense-tracker-app/internal/dto"
	"go-income-expense-tracker-app/internal/middleware"
	"go-income-expense-tracker-app/internal/model"
	"go-income-expense-tracker-app/internal/service"
	"go-income-expense-tracker-app/internal/utils"
)

type MockAuthRepoForAuth struct {
	mock.Mock
}

func (m *MockAuthRepoForAuth) Create(user *model.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockAuthRepoForAuth) FindByEmail(email string) (*model.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func newAuthService() (*MockAuthRepoForAuth, service.AuthService) {
	mockRepo := new(MockAuthRepoForAuth)
	jwtCfg := &middleware.JWTConfig{
		SecretKey:       "test-secret-key",
		ExpiresDuration: 60,
	}
	svc := service.NewAuthService(mockRepo, *jwtCfg)
	return mockRepo, svc
}

func TestAuthService_Register(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo, svc := newAuthService()

		mockRepo.On("FindByEmail", "john@example.com").Return(nil, fmt.Errorf("not found"))
		mockRepo.On("Create", mock.AnythingOfType("*model.User")).Return(nil)

		user, err := svc.Register(dto.RegisterRequest{
			Username: "john",
			Email:    "john@example.com",
			Password: "password123",
		})
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, "john", user.Username)
		assert.Equal(t, "john@example.com", user.Email)
		assert.Equal(t, model.RoleUser, user.Role)
		assert.NotEmpty(t, user.Password)
		assert.NotEqual(t, "password123", user.Password)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid email format", func(t *testing.T) {
		_, svc := newAuthService()

		user, err := svc.Register(dto.RegisterRequest{
			Username: "john",
			Email:    "not-an-email",
			Password: "password123",
		})
		assert.ErrorIs(t, err, service.ErrInvalidEmailFormat)
		assert.Nil(t, user)
	})

	t.Run("password too short", func(t *testing.T) {
		_, svc := newAuthService()

		user, err := svc.Register(dto.RegisterRequest{
			Username: "john",
			Email:    "john@example.com",
			Password: "123",
		})
		assert.ErrorIs(t, err, service.ErrPasswordTooShort)
		assert.Nil(t, user)
	})

	t.Run("email already exists", func(t *testing.T) {
		mockRepo, svc := newAuthService()

		existingUser := &model.User{ID: 1, Email: "john@example.com"}
		mockRepo.On("FindByEmail", "john@example.com").Return(existingUser, nil)

		user, err := svc.Register(dto.RegisterRequest{
			Username: "john",
			Email:    "john@example.com",
			Password: "password123",
		})
		assert.ErrorIs(t, err, service.ErrEmailExists)
		assert.Nil(t, user)
	})

	t.Run("repository create error", func(t *testing.T) {
		mockRepo, svc := newAuthService()

		mockRepo.On("FindByEmail", "john@example.com").Return(nil, fmt.Errorf("not found"))
		mockRepo.On("Create", mock.AnythingOfType("*model.User")).
			Return(fmt.Errorf("db error"))

		user, err := svc.Register(dto.RegisterRequest{
			Username: "john",
			Email:    "john@example.com",
			Password: "password123",
		})
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "db error")
	})
}

func TestAuthService_Login(t *testing.T) {
	hashedPassword, _ := utils.HashPassword("password123")
	validUser := &model.User{
		ID:       1,
		Username: "john",
		Email:    "john@example.com",
		Password: hashedPassword,
		Role:     model.RoleUser,
	}

	t.Run("success", func(t *testing.T) {
		mockRepo, svc := newAuthService()

		mockRepo.On("FindByEmail", "john@example.com").Return(validUser, nil)

		token, err := svc.Login(dto.LoginRequest{
			Email:    "john@example.com",
			Password: "password123",
		})
		assert.NoError(t, err)
		assert.NotEmpty(t, token)
	})

	t.Run("email not found", func(t *testing.T) {
		mockRepo, svc := newAuthService()

		mockRepo.On("FindByEmail", "unknown@example.com").Return(nil, fmt.Errorf("not found"))

		token, err := svc.Login(dto.LoginRequest{
			Email:    "unknown@example.com",
			Password: "password123",
		})
		assert.ErrorIs(t, err, service.ErrInvalidCredentials)
		assert.Empty(t, token)
	})

	t.Run("wrong password", func(t *testing.T) {
		mockRepo, svc := newAuthService()

		mockRepo.On("FindByEmail", "john@example.com").Return(validUser, nil)

		token, err := svc.Login(dto.LoginRequest{
			Email:    "john@example.com",
			Password: "wrongpassword",
		})
		assert.ErrorIs(t, err, service.ErrInvalidCredentials)
		assert.Empty(t, token)
	})
}
