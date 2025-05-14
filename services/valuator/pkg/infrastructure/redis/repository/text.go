package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/mono83/maybe"
	"github.com/redis/go-redis/v9"

	"valuator/pkg/app/model"
	"valuator/pkg/app/repository"
	"valuator/pkg/infrastructure/redis/keyvalue"
)

func NewTextRepository(client *redis.Client) repository.TextRepository {
	return &textRepository{
		storage: keyvalue.NewStorage[textSerializable](client),
	}
}

type textRepository struct {
	storage keyvalue.Storage[textSerializable]
}

type textSerializable struct {
	ID    string `json:"id"`
	Value string `json:"value"`
	Count int    `json:"count"`
}

func (repository *textRepository) Store(_ string, text model.Text) error {
	return repository.storage.Set(context.Background(), fmt.Sprintf("text:%s", text.ID), textSerializable{
		ID:    string(text.ID),
		Value: text.Value,
		Count: text.Count,
	}, 0)
}

func (repository *textRepository) Remove(text model.Text) error {
	return repository.storage.Delete(context.Background(), fmt.Sprintf("text:%s", text.ID))
}

func (repository *textRepository) Find(id model.TextID) (maybe.Maybe[model.Text], error) {
	v, err := repository.storage.Get(context.Background(), fmt.Sprintf("text:%s", id))
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return maybe.Nothing[model.Text](), nil
		}
		return maybe.Nothing[model.Text](), err
	}

	textModel, err := repository.convertToModel(v)
	if err != nil {
		return maybe.Nothing[model.Text](), err
	}
	return maybe.Just(textModel), nil
}

func (repository *textRepository) convertToModel(text textSerializable) (model.Text, error) {
	return model.Text{
		ID:    model.TextID(text.ID),
		Value: text.Value,
		Count: text.Count,
	}, nil
}
