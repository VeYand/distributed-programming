package service

import (
	"github.com/gofrs/uuid"
	"valuator/pkg/app/model"
	"valuator/pkg/app/repository"
)

type TextService interface {
	Add(value string) (uuid.UUID, error)
	Remove(id uuid.UUID) error
}

func NewTextService(repository repository.TextRepository) TextService {
	return &textService{repository: repository}
}

type textService struct {
	repository repository.TextRepository
}

func (s *textService) Add(value string) (uuid.UUID, error) {
	text := s.createText(value)
	err := s.repository.Store(text)
	if err != nil {
		return uuid.UUID{}, err
	}
	return uuid.UUID(text.ID), nil
}

func (s *textService) Remove(id uuid.UUID) error {
	text, err := s.repository.Find(model.TextID(id))
	if err != nil {
		return err
	}
	if text.IsPresent() {
		return s.repository.Remove(text.Value())
	}
	return nil
}

func (s *textService) createText(value string) model.Text {
	newUuid, err := uuid.NewV7()
	if err != nil {
		panic(err)
	}

	return model.Text{
		ID:    model.TextID(newUuid),
		Value: value,
	}
}
