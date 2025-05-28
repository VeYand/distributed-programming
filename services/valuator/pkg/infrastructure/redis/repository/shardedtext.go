package repo

import (
	"github.com/mono83/maybe"
	"valuator/pkg/app/model"
	"valuator/pkg/app/repository"
)

func NewTextShardedRepository(shardManager *ShardManager) repository.TextRepository {
	return &textShardedRepository{
		shardManager: shardManager,
	}
}

type textShardedRepository struct {
	shardManager *ShardManager
}

func (repository *textShardedRepository) Find(id model.TextID) (maybe.Maybe[model.Text], error) {
	repo, err := repository.shardManager.GetTextRepository(id)
	if err != nil {
		return maybe.Nothing[model.Text](), err
	}

	return repo.Find(id)
}

func (repository *textShardedRepository) Store(region string, text model.Text) error {
	err := repository.shardManager.Store(region, text.ID)
	if err != nil {
		return err
	}

	repo, err := repository.shardManager.GetTextRepository(text.ID)
	if err != nil {
		return err
	}

	return repo.Store(region, text)
}

func (repository *textShardedRepository) Remove(text model.Text) error {
	repo, err := repository.shardManager.GetTextRepository(text.ID)
	if err != nil {
		return err
	}

	return repo.Remove(text)
}
