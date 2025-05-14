package repository

import (
	"github.com/mono83/maybe"

	"valuator/pkg/app/model"
)

type TextReadRepository interface {
	Find(id model.TextID) (maybe.Maybe[model.Text], error)
}

type TextRepository interface {
	TextReadRepository
	Store(region string, text model.Text) error
	Remove(text model.Text) error
}
