package query

import (
	"errors"
	"user/pkg/app/model"
)

var ErrUserNotFound = errors.New("user not found")

type UserQueryService interface {
	FindByLogin(login string) (model.User, error)
	FindByID(id string) (model.User, error)
}

func NewUserQueryService(repository model.UserReadRepository) UserQueryService {
	return &service{repository: repository}
}

type service struct {
	repository model.UserReadRepository
}

func (s *service) FindByLogin(login string) (model.User, error) {
	user, err := s.repository.FindByLogin(login)
	if err != nil {
		return model.User{}, err
	}
	justUser, isJust := user.Get()
	if !isJust {
		return model.User{}, ErrUserNotFound
	}
	return justUser, nil
}

func (s *service) FindByID(id string) (model.User, error) {
	user, err := s.repository.Find(model.UserID(id))
	if err != nil {
		return model.User{}, err
	}
	justUser, isJust := user.Get()
	if !isJust {
		return model.User{}, ErrUserNotFound
	}
	return justUser, nil
}
