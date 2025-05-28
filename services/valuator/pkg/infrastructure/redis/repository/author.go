package repo

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"

	"valuator/pkg/app/model"
	"valuator/pkg/app/repository"
	"valuator/pkg/infrastructure/redis/keyvalue"
)

const (
	authorsKeyPattern = "text:%s:authors"
)

type authorRepository struct {
	storage keyvalue.Storage[[]string]
}

func NewAuthorRepository(client *redis.Client) repository.AuthorRepository {
	return &authorRepository{
		storage: keyvalue.NewStorage[[]string](client),
	}
}

func (r *authorRepository) FindAuthors(textID model.TextID) ([]model.AuthorID, error) {
	key := fmt.Sprintf(authorsKeyPattern, textID)
	list, err := r.storage.Get(context.Background(), key)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return []model.AuthorID{}, nil
		}
		return nil, err
	}

	result := make([]model.AuthorID, 0, len(list))
	for _, id := range list {
		result = append(result, id)
	}
	return result, nil
}

func (r *authorRepository) Store(author model.Author) error {
	ctx := context.Background()
	key := fmt.Sprintf(authorsKeyPattern, author.TextID)

	list, err := r.storage.Get(ctx, key)
	if err != nil && !errors.Is(err, redis.Nil) {
		return err
	}
	if errors.Is(err, redis.Nil) {
		list = []string{}
	}

	for _, existing := range list {
		if existing == author.AuthorID {
			log.Println("no set")
			return nil
		}
	}

	list = append(list, author.AuthorID)

	log.Println("set authors+", key, "    :    ", list)

	return r.storage.Set(ctx, key, list, 0)
}
