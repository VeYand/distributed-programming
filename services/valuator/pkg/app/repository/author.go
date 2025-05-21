package repository

import (
	"valuator/pkg/app/model"
)

type AuthorReadRepository interface {
	FindAuthors(textID model.TextID) ([]model.AuthorID, error)
}

type AuthorRepository interface {
	AuthorReadRepository
	Store(author model.Author) error
}
