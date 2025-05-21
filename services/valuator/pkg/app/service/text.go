package service

import (
	"crypto/sha256"
	"encoding/hex"
	stderrors "errors"
	"valuator/pkg/app/errors"
	"valuator/pkg/app/event"
	"valuator/pkg/app/message"
	"valuator/pkg/app/model"
	"valuator/pkg/app/repository"
)

type TextService interface {
	Add(authorID model.AuthorID, region, value string) (string, error)
}

func NewTextService(
	textRepository repository.TextRepository,
	authorRepository repository.AuthorRepository,
	messagePublisher message.Publisher,
	eventDispatcher event.Dispatcher,
) TextService {
	return &textService{
		textRepository:   textRepository,
		authorRepository: authorRepository,
		messagePublisher: messagePublisher,
		eventDispatcher:  eventDispatcher,
	}
}

type textService struct {
	textRepository   repository.TextRepository
	authorRepository repository.AuthorRepository
	messagePublisher message.Publisher
	eventDispatcher  event.Dispatcher
}

func (s *textService) Add(authorID model.AuthorID, region, value string) (string, error) {
	text := s.createText(value)
	existingTextModel, err := s.textRepository.Find(text.ID)
	if err != nil && !stderrors.Is(err, errors.ErrTextNotFound) {
		return "", err
	}
	existingText, isPresent := existingTextModel.Get()
	if isPresent {
		existingText.Count++
	} else {
		existingText = text
	}

	err = s.textRepository.Store(region, existingText)
	if err != nil {
		return "", err
	}

	err = s.authorRepository.Store(model.Author{
		AuthorID: authorID,
		TextID:   existingText.ID,
	})
	if err != nil {
		return "", err
	}

	err = s.eventDispatcher.Dispatch(event.CreateSimilarityCalculatedEvent(existingText.ID, existingText.Count > 0))
	if err != nil {
		return "", err
	}

	return string(text.ID), s.messagePublisher.Publish(message.NewTextAddedMessage(existingText))
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
