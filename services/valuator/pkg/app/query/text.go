package query

import (
	"github.com/gofrs/uuid"
	"valuator/pkg/app/errors"
	"valuator/pkg/app/model"
	"valuator/pkg/app/repository"
)

type TextQueryService interface {
	List() ([]TextData, error)
	Get(id uuid.UUID) (TextData, error)
}

func NewTextQueryService(textReadRepository repository.TextReadRepository) TextQueryService {
	return &textQueryService{
		textReadRepository: textReadRepository,
	}
}

type textQueryService struct {
	textReadRepository repository.TextReadRepository
}

type TextData struct {
	ID    uuid.UUID
	Value string
}

func (s *textQueryService) List() ([]TextData, error) {
	texts, err := s.textReadRepository.ListAll()
	if err != nil {
		return nil, err
	}

	results := make([]TextData, 0, len(texts))
	for _, text := range texts {
		results = append(results, TextData{
			ID:    uuid.UUID(text.ID),
			Value: text.Value,
		})
	}
	return results, nil
}

func (s *textQueryService) Get(id uuid.UUID) (TextData, error) {
	text, err := s.textReadRepository.Find(model.TextID(id))
	if err != nil {
		return TextData{}, err
	}

	if text.IsEmpty() {
		return TextData{}, errors.ErrTextNotFound
	}

	textValue := text.Value()

	return TextData{
		ID:    uuid.UUID(textValue.ID),
		Value: textValue.Value,
	}, nil
}
