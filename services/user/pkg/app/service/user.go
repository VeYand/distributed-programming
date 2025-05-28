package service

import (
	"errors"
	"github.com/gofrs/uuid"
	"user/pkg/app/model"
	"user/pkg/app/utils"
)

var ErrUserAlreadyExists = errors.New("user already exists")

type UserService interface {
	Create(login, password string) error
}

func NewUserService(repository model.UserRepository) UserService {
	return &service{repository: repository}
}

type service struct {
	repository model.UserRepository
}

func (s *service) Create(login, password string) error {
	newUuid, err := uuid.NewV7()
	if err != nil {
		return err
	}

	existingUser, err := s.repository.FindByLogin(login)
	if err != nil {
		return err
	}

	if existingUser.IsPresent() {
		return ErrUserAlreadyExists
	}

	return s.repository.Store(model.User{
		UserID:   model.UserID(newUuid.String()),
		Login:    login,
		Password: utils.HashPassword(password),
	})
}
