package model

import (
	"github.com/mono83/maybe"
)

type UserID string

type User struct {
	UserID   UserID
	Login    string
	Password string
}

type UserReadRepository interface {
	Find(id UserID) (maybe.Maybe[User], error)
	FindByLogin(login string) (maybe.Maybe[User], error)
}

type UserRepository interface {
	UserReadRepository
	Store(user User) error
}
