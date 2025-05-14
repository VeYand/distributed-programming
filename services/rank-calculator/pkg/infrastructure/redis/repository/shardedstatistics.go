package repo

import (
	"github.com/mono83/maybe"
	"rankcalculator/pkg/app/model"
	"rankcalculator/pkg/app/repository"
)

func NewStatisticsShardedRepository(shardManager *ShardManager) repository.StatisticsRepository {
	return &statisticsShardedRepository{
		shardManager: shardManager,
	}
}

type statisticsShardedRepository struct {
	shardManager *ShardManager
}

func (repository *statisticsShardedRepository) Find(ID model.TextID) (maybe.Maybe[model.Statistics], error) {
	repo, err := repository.shardManager.GetRepository(ID)
	if err != nil {
		return maybe.Nothing[model.Statistics](), err
	}

	return repo.Find(ID)
}

func (repository *statisticsShardedRepository) Store(text model.Statistics) error {
	repo, err := repository.shardManager.GetRepository(text.TextID)
	if err != nil {
		return err
	}

	return repo.Store(text)
}
