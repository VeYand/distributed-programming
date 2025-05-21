package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/mono83/maybe"
	"github.com/redis/go-redis/v9"

	"user/pkg/app/model"
	"user/pkg/infrastructure/redis/keyvalue"
)

func NewUserRepository(client *redis.Client) model.UserRepository {
	return &userRepository{
		userStorage:   keyvalue.NewStorage[userSerializable](client),
		userIDStorage: keyvalue.NewStorage[userIDSerializable](client),
	}
}

type userRepository struct {
	userStorage   keyvalue.Storage[userSerializable]
	userIDStorage keyvalue.Storage[userIDSerializable]
}

type userSerializable struct {
	UserID   string `json:"user_id"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

type userIDSerializable struct {
	UserID string `json:"user_id"`
}

func (repository *userRepository) FindByLogin(login string) (maybe.Maybe[model.User], error) {
	data, err := repository.userIDStorage.Get(context.Background(), fmt.Sprintf("user:login:%s", login))
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return maybe.Nothing[model.User](), nil
		}
		return maybe.Nothing[model.User](), err
	}

	return repository.Find(model.UserID(data.UserID))
}

func (repository *userRepository) Find(id model.UserID) (maybe.Maybe[model.User], error) {
	v, err := repository.userStorage.Get(context.Background(), fmt.Sprintf("user:id:%s", id))
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return maybe.Nothing[model.User](), nil
		}
		return maybe.Nothing[model.User](), err
	}

	userModel, err := repository.convertToModel(v)
	if err != nil {
		return maybe.Nothing[model.User](), err
	}
	return maybe.Just(userModel), nil
}

func (repository *userRepository) Store(user model.User) error {
	ctx := context.Background()
	keyID := fmt.Sprintf("user:id:%s", user.UserID)
	err := repository.userStorage.Set(context.Background(), keyID, userSerializable{
		UserID:   string(user.UserID),
		Login:    user.Login,
		Password: user.Password,
	}, 0)
	if err != nil {
		return err
	}
	keyLogin := fmt.Sprintf("user:login:%s", user.Login)
	return repository.userIDStorage.Set(ctx, keyLogin,
		userIDSerializable{
			UserID: string(user.UserID),
		}, 0)
}

func (repository *userRepository) convertToModel(user userSerializable) (model.User, error) {
	return model.User{
		UserID:   model.UserID(user.UserID),
		Login:    user.Login,
		Password: user.Password,
	}, nil
}
