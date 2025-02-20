package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/gofrs/uuid"
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
}

func (repository *textRepository) Store(text model.Text) error {
	return repository.storage.Set(context.Background(), fmt.Sprintf("text:%s", uuid.UUID(text.ID).String()), textSerializable{
		ID:    uuid.UUID(text.ID).String(),
		Value: text.Value,
	}, 0)
}

func (repository *textRepository) Remove(text model.Text) error {
	return repository.storage.Delete(context.Background(), fmt.Sprintf("text:%s", uuid.UUID(text.ID).String()))
}

func (repository *textRepository) Find(id model.TextID) (maybe.Maybe[model.Text], error) {
	v, err := repository.storage.Get(context.Background(), fmt.Sprintf("text:%s", uuid.UUID(id).String()))
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

func (repository *textRepository) ListAll() ([]model.Text, error) {
	vs, err := repository.storage.ListAll(context.Background(), "text:*")
	if err != nil {
		return nil, err
	}
	texts := make([]model.Text, 0, len(vs))
	for _, v := range vs {
		textModel, err := repository.convertToModel(v)
		if err != nil {
			return nil, err
		}
		texts = append(texts, textModel)
	}
	return texts, nil
}

func (repository *textRepository) convertToModel(text textSerializable) (model.Text, error) {
	id, err := uuid.FromString(text.ID)
	if err != nil {
		return model.Text{}, err
	}
	return model.Text{
		ID:    model.TextID(id),
		Value: text.Value,
	}, nil
}
