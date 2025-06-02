package service

import (
	"errors"
	"regexp"
)

type Service interface {
	Set(key string, value int) error
	Get(key string) (int, error)
	Keys(prefix string) ([]string, error)
}

type ProtoKeyService struct {
	commands chan Command
}

var (
	keyPattern    = regexp.MustCompile(`^[a-zA-Z0-9_\-.]{1,1000}$`)
	ErrBadRequest = errors.New("bad request")
)

func NewProtoKeyService(cmdCh chan Command) *ProtoKeyService {
	return &ProtoKeyService{commands: cmdCh}
}

func (service *ProtoKeyService) Set(key string, value int) error {
	if !keyPattern.MatchString(key) {
		return ErrBadRequest
	}

	reply := make(chan Response, 1)
	command := Command{
		Type:  SetOperation,
		Key:   key,
		Value: value,
		Reply: reply,
	}

	service.commands <- command
	response := <-reply
	return response.Err
}

func (service *ProtoKeyService) Get(key string) (int, error) {
	if !keyPattern.MatchString(key) {
		return 0, ErrBadRequest
	}

	reply := make(chan Response, 1)
	command := Command{
		Type:  GetOperation,
		Key:   key,
		Reply: reply,
	}

	service.commands <- command
	response := <-reply
	return response.Value, response.Err
}

func (service *ProtoKeyService) Keys(prefix string) ([]string, error) {
	if !keyPattern.MatchString(prefix) {
		return nil, ErrBadRequest
	}

	reply := make(chan Response, 1)
	command := Command{
		Type:  ListKeysOperation,
		Key:   prefix,
		Reply: reply,
	}

	service.commands <- command
	response := <-reply
	return response.Keys, response.Err
}
