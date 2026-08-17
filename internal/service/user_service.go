package service

import (
	"errors"
	"go-income-expense-tracker-app/internal/dto"
	"go-income-expense-tracker-app/internal/model"
	"go-income-expense-tracker-app/internal/repository"
	"go-income-expense-tracker-app/internal/utils"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type UserService interface {
	GetByID(id int) (*model.User, error)
	GetAll() ([]model.User, error)
	Update(id int, req dto.UpdateUserRequest) (*model.User, error)
	Delete(id int) error
}

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return &userService{userRepository: userRepository}
}

func (s *userService) GetByID(id int) (*model.User, error) {
	user, err := s.userRepository.FindByID(id)
	if err != nil {
		return nil, ErrUserNotFound
	}

	return user, nil
}

func (s *userService) GetAll() ([]model.User, error) {
	return s.userRepository.FindAll()
}

func (s *userService) Update(id int, req dto.UpdateUserRequest) (*model.User, error) {
	user, err := s.userRepository.FindByID(id)
	if err != nil {
		return nil, ErrUserNotFound
	}

	if req.Username != "" {
		user.Username = req.Username
	}

	if req.Email != "" {
		user.Email = req.Email
	}

	if req.Password != "" {
		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
			return nil, err
		}

		user.Password = hashedPassword
	}

	if err := s.userRepository.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) Delete(id int) error {
	if _, err := s.userRepository.FindByID(id); err != nil {
		return ErrUserNotFound
	}

	return s.userRepository.Delete(id)
}
