package service

import (
	"errors"

	"github.com/wikasdude/blog-backend/model"
	repository "github.com/wikasdude/blog-backend/repositories"
)

type UserService interface {
	RegisterUser(user *model.User) error
	GetUserByID(id string) (*model.User, error)
	UpdateUser(user *model.User) error
	DeleteUser(id int) error
	GetUserByEmail(email string) (*model.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) RegisterUser(user *model.User) error {
	// Add validation or password hashing here
	return s.repo.CreateUser(user)
}

func (s *userService) GetUserByID(id string) (*model.User, error) {
	return s.repo.GetUserByID(id)
}
func (s *userService) UpdateUser(user *model.User) error {
	exists, err := s.repo.IsEmailTaken(user.Email, user.ID)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("email is already taken")
	}
	return s.repo.UpdateUser(user)
}
func (s *userService) DeleteUser(id int) error {
	return s.repo.DeleteUser(id)
}
func (s *userService) GetUserByEmail(email string) (*model.User, error) {
	return s.repo.GetUserByEmail(email)
}
