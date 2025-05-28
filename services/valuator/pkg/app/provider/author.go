package provider

import (
	"valuator/pkg/app/model"
	"valuator/pkg/app/repository"
)

type AuthorProvider interface {
	ListTextAuthors(textID model.TextID) ([]model.AuthorID, error)
}

func NewAuthorProvider(readRepository repository.AuthorReadRepository) AuthorProvider {
	return &authorProvider{authorReadRepository: readRepository}
}

type authorProvider struct {
	authorReadRepository repository.AuthorReadRepository
}

func (a *authorProvider) ListTextAuthors(textID model.TextID) ([]model.AuthorID, error) {
	return a.authorReadRepository.FindAuthors(textID)
}
