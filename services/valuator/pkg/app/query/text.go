package query

import (
	"valuator/pkg/app/data"
	"valuator/pkg/app/errors"
	"valuator/pkg/app/model"
	"valuator/pkg/app/repository"
)

type TextQueryService interface {
	Get(id string) (data.TextData, error)
}

func NewTextQueryService(textReadRepository repository.TextReadRepository) TextQueryService {
	return &textQueryService{
		textReadRepository: textReadRepository,
	}
}

type textQueryService struct {
	textReadRepository repository.TextReadRepository
}

func (s *textQueryService) Get(id string) (data.TextData, error) {
	text, err := s.textReadRepository.Find(model.TextID(id))
	if err != nil {
		return data.TextData{}, err
	}

	if text.IsEmpty() {
		return data.TextData{}, errors.ErrTextNotFound
	}

	textValue := text.Value()

	return data.TextData{
		ID:    string(textValue.ID),
		Value: textValue.Value,
		Count: textValue.Count,
	}, nil
}
