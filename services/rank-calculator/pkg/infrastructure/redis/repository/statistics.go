package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/mono83/maybe"
	"github.com/redis/go-redis/v9"

	"rankcalculator/pkg/app/model"
	"rankcalculator/pkg/app/repository"
	"rankcalculator/pkg/infrastructure/redis/keyvalue"
)

func NewStatisticsRepository(client *redis.Client) repository.StatisticsRepository {
	return &statisticsRepository{
		storage: keyvalue.NewStorage[statisticsSerializable](client),
	}
}

type statisticsRepository struct {
	storage keyvalue.Storage[statisticsSerializable]
}

type statisticsSerializable struct {
	TextID               string `json:"text_id"`
	AlphabetSymbolsCount int    `json:"alphabet_symbols_count"`
	AllSymbolsCount      int    `json:"all_symbols_count"`
	IsDuplicate          bool   `json:"is_duplicate"`
}

func (repository *statisticsRepository) Store(statistics model.Statistics) error {
	return repository.storage.Set(context.Background(), fmt.Sprintf("statistics:%s", statistics.TextID), statisticsSerializable{
		TextID:               string(statistics.TextID),
		AlphabetSymbolsCount: statistics.AlphabetSymbolsCount,
		AllSymbolsCount:      statistics.AllSymbolsCount,
		IsDuplicate:          statistics.IsDuplicate,
	}, 0)
}

func (repository *statisticsRepository) Find(ID model.TextID) (maybe.Maybe[model.Statistics], error) {
	v, err := repository.storage.Get(context.Background(), fmt.Sprintf("statistics:%s", ID))
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return maybe.Nothing[model.Statistics](), nil
		}
		return maybe.Nothing[model.Statistics](), err
	}

	statisticsModel, err := repository.convertToModel(v)
	if err != nil {
		return maybe.Nothing[model.Statistics](), err
	}
	return maybe.Just(statisticsModel), nil
}

func (repository *statisticsRepository) convertToModel(statistics statisticsSerializable) (model.Statistics, error) {
	return model.Statistics{
		TextID:               model.TextID(statistics.TextID),
		AlphabetSymbolsCount: statistics.AlphabetSymbolsCount,
		AllSymbolsCount:      statistics.AllSymbolsCount,
		IsDuplicate:          statistics.IsDuplicate,
	}, nil
}
