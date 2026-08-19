package service

import (
	"errors"
	"go-income-expense-tracker-app/internal/dto"
	"go-income-expense-tracker-app/internal/middleware"
	"go-income-expense-tracker-app/internal/model"
	"go-income-expense-tracker-app/internal/repository"
	"go-income-expense-tracker-app/internal/utils"

	"github.com/go-playground/validator/v10"
)

var (
	ErrEmailExists        = errors.New("email already registered")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrPasswordTooShort   = errors.New("password must be at least 6 characters")
	ErrInvalidEmailFormat = errors.New("invalid email format")
)

type AuthService interface {
	Register(req dto.RegisterRequest) (*model.User, error)
	Login(req dto.LoginRequest) (string, error)
}

type authService struct {
	authRepository repository.AuthRepository
	jwtConfig      middleware.JWTConfig
}

func NewAuthService(authRepository repository.AuthRepository, jwtConfig middleware.JWTConfig) AuthService {
	return &authService{authRepository: authRepository, jwtConfig: jwtConfig}
}

func (s *authService) Register(req dto.RegisterRequest) (*model.User, error) {
	v := validator.New()
	if err := v.Var(req.Email, "email"); err != nil {
		return nil, ErrInvalidEmailFormat
	}

	if len(req.Password) < 6 {
		return nil, ErrPasswordTooShort
	}

	if _, err := s.authRepository.FindByEmail(req.Email); err == nil {
		return nil, ErrEmailExists
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
		Role:     model.RoleUser,
	}

	if err := s.authRepository.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authService) Login(req dto.LoginRequest) (string, error) {
	user, err := s.authRepository.FindByEmail(req.Email)
	if err != nil {
		return "", ErrInvalidCredentials
	}

	if !utils.CheckPassword(user.Password, req.Password) {
		return "", ErrInvalidCredentials
	}

	token, err := s.jwtConfig.GenerateToken(int(user.ID), user.Role)
	if err != nil {
		return "", err
	}

	return token, nil
}
