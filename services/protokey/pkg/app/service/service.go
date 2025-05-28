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
	commands  chan Command
	responses chan Response
}

var (
	keyPattern    = regexp.MustCompile(`^[a-zA-Z0-9_\-.]{1,1000}$`)
	ErrBadRequest = errors.New("bad request")
)

func NewProtoKeyService(cmdCh chan Command, respCh chan Response) *ProtoKeyService {
	return &ProtoKeyService{commands: cmdCh, responses: respCh}
}

func (service *ProtoKeyService) Set(key string, value int) error {
	if !keyPattern.MatchString(key) {
		return ErrBadRequest
	}

	service.commands <- Command{Type: SetOperation, Key: key, Value: value}
	resp := <-service.responses

	return resp.Err
}

func (service *ProtoKeyService) Get(key string) (int, error) {
	if !keyPattern.MatchString(key) {
		return 0, ErrBadRequest
	}

	service.commands <- Command{Type: GetOperation, Key: key}
	resp := <-service.responses

	return resp.Value, resp.Err
}

func (service *ProtoKeyService) Keys(prefix string) ([]string, error) {
	if !keyPattern.MatchString(prefix) {
		return nil, ErrBadRequest
	}

	service.commands <- Command{Type: ListKeysOperation, Key: prefix}
	resp := <-service.responses

	return resp.Keys, resp.Err
}
