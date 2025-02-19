package repository

import (
	"github.com/mono83/maybe"

	"valuator/pkg/app/model"
)

type TextReadRepository interface {
	Find(id model.TextID) (maybe.Maybe[model.Text], error)
	ListAll() ([]model.Text, error)
}

type TextRepository interface {
	TextReadRepository
	Store(text model.Text) error
	Remove(text model.Text) error
}
