package query

import (
	"github.com/gofrs/uuid"
	"valuator/pkg/app/data"
	"valuator/pkg/app/errors"
	"valuator/pkg/app/model"
	"valuator/pkg/app/repository"
)

type TextQueryService interface {
	List() ([]data.TextData, error)
	Get(id uuid.UUID) (data.TextData, error)
}

func NewTextQueryService(textReadRepository repository.TextReadRepository) TextQueryService {
	return &textQueryService{
		textReadRepository: textReadRepository,
	}
}

type textQueryService struct {
	textReadRepository repository.TextReadRepository
}

func (s *textQueryService) List() ([]data.TextData, error) {
	texts, err := s.textReadRepository.ListAll()
	if err != nil {
		return nil, err
	}

	results := make([]data.TextData, 0, len(texts))
	for _, text := range texts {
		results = append(results, data.TextData{
			ID:    uuid.UUID(text.ID),
			Value: text.Value,
		})
	}
	return results, nil
}

func (s *textQueryService) Get(id uuid.UUID) (data.TextData, error) {
	text, err := s.textReadRepository.Find(model.TextID(id))
	if err != nil {
		return data.TextData{}, err
	}

	if text.IsEmpty() {
		return data.TextData{}, errors.ErrTextNotFound
	}

	textValue := text.Value()

	return data.TextData{
		ID:    uuid.UUID(textValue.ID),
		Value: textValue.Value,
	}, nil
}
