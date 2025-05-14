package repo

import (
	"context"
	stderrors "errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"rankcalculator/pkg/app/errors"
	"rankcalculator/pkg/app/model"
	"rankcalculator/pkg/app/repository"
	"rankcalculator/pkg/infrastructure/redis/keyvalue"
)

func NewShardManager(
	mainClient *redis.Client,
	shards map[string]*redis.Client,
) *ShardManager {
	return &ShardManager{
		textRegionStorage: keyvalue.NewStorage[textRegionSerializable](mainClient),
		shards:            shards,
	}
}

type ShardManager struct {
	textRegionStorage keyvalue.Storage[textRegionSerializable]
	shards            map[string]*redis.Client
}

type textRegionSerializable struct {
	TextID string `json:"id"`
	Region string `json:"region"`
}

func (m *ShardManager) Store(region string, textID model.TextID) error {
	_, err := m.getShard(region)
	if err != nil {
		return err
	}
	return m.textRegionStorage.Set(
		context.Background(),
		fmt.Sprintf("text_region:%s", textID),
		textRegionSerializable{
			TextID: string(textID),
			Region: region,
		},
		0,
	)
}

func (m *ShardManager) GetRepository(textID model.TextID) (repository.StatisticsRepository, error) {
	region, err := m.textRegionStorage.Get(context.Background(), fmt.Sprintf("text_region:%s", textID))
	if err != nil {
		if stderrors.Is(err, redis.Nil) {
			return nil, errors.ErrStatisticsNotFound
		}
		return nil, err
	}

	shard, err := m.getShard(region.Region)
	if err != nil {
		return nil, err
	}

	fmt.Println()
	log.Printf("LOOKUP: %s, %s", textID, region.Region)
	fmt.Println()

	return NewStatisticsRepository(shard), nil
}

func (m *ShardManager) getShard(region string) (*redis.Client, error) {
	shard, ok := m.shards[region]
	if !ok {
		return nil, fmt.Errorf("region %s not found in shard map", region)
	}

	return shard, nil
}
