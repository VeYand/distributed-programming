package service

import (
	"crypto/sha256"
	"encoding/hex"
	"valuator/pkg/app/model"
	"valuator/pkg/app/repository"
)

type TextService interface {
	Add(value string) (string, error)
}

func NewTextService(repository repository.TextRepository) TextService {
	return &textService{repository: repository}
}

type textService struct {
	repository repository.TextRepository
}

func (s *textService) Add(value string) (string, error) {
	text := s.createText(value)
	existingTextModel, err := s.repository.Find(text.ID)
	if err != nil {
		return "", err
	}
	existingText, isPresent := existingTextModel.Get()
	if isPresent {
		existingText.Count++
	} else {
		existingText = text
	}

	err = s.repository.Store(existingText)
	if err != nil {
		return "", err
	}
	return string(text.ID), nil
}

func (s *textService) createText(value string) model.Text {
	newID := hashText(value)
	return model.Text{
		ID:    model.TextID(newID),
		Value: value,
		Count: 1,
	}
}

func hashText(text string) string {
	hash := sha256.New()
	hash.Write([]byte(text))
	return hex.EncodeToString(hash.Sum(nil))
}
