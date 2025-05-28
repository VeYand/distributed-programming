package session

import (
	"errors"
	"user/pkg/app/model"
	"user/pkg/app/query"
	"user/pkg/app/utils"
)

var ErrInvalidCredentials = errors.New("invalid user credentials")

type UserSession interface {
	Identify(login, password string) (model.User, error)
}

func NewUserSession(userQueryService query.UserQueryService) UserSession {
	return &session{
		userQueryService: userQueryService,
	}
}

type session struct {
	userQueryService query.UserQueryService
}

func (s *session) Identify(login, password string) (model.User, error) {
	user, err := s.userQueryService.FindByLogin(login)
	if errors.Is(err, query.ErrUserNotFound) {
		return model.User{}, ErrInvalidCredentials
	}
	if err != nil {
		return model.User{}, err
	}

	if !utils.ComparePassword(user.Password, password) {
		return model.User{}, ErrInvalidCredentials
	}

	return user, nil
}
