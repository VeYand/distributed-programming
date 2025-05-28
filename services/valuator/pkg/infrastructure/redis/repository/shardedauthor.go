package repo

import (
	"valuator/pkg/app/model"
	"valuator/pkg/app/repository"
)

func NewAuthorShardedRepository(shardManager *ShardManager) repository.AuthorRepository {
	return &authorShardedRepository{
		shardManager: shardManager,
	}
}

type authorShardedRepository struct {
	shardManager *ShardManager
}

func (repository *authorShardedRepository) FindAuthors(textID model.TextID) ([]model.AuthorID, error) {
	repo, err := repository.shardManager.GetAuthorRepository(textID)
	if err != nil {
		return nil, err
	}

	return repo.FindAuthors(textID)
}

func (repository *authorShardedRepository) Store(author model.Author) error {
	repo, err := repository.shardManager.GetAuthorRepository(author.TextID)
	if err != nil {
		return err
	}

	return repo.Store(author)
}
