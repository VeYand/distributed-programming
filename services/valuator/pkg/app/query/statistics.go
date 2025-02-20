package query

import (
	"github.com/gofrs/uuid"
	"unicode"
	"unicode/utf8"
	"valuator/pkg/app/errors"
	"valuator/pkg/app/model"
	"valuator/pkg/app/repository"
)

type StatisticsQueryService interface {
	GetSummary(textID uuid.UUID) (TextStatistics, error)
}

func NewStatisticsQueryService(textReadRepository repository.TextReadRepository) StatisticsQueryService {
	return &statisticsQueryService{
		textReadRepository: textReadRepository,
	}
}

type statisticsQueryService struct {
	textReadRepository repository.TextReadRepository
}

type TextStatistics struct {
	SymbolStatistics
	UniqueStatistics
}

type SymbolStatistics struct {
	AlphabetSymbolsCount int
	AllSymbolsCount      int
}

type UniqueStatistics struct {
	IsDuplicate bool
}

func (queryService *statisticsQueryService) GetSummary(textID uuid.UUID) (TextStatistics, error) {
	text, err := queryService.textReadRepository.Find(model.TextID(textID))
	if err != nil {
		return TextStatistics{}, err
	}
	if text.IsEmpty() {
		return TextStatistics{}, errors.ErrTextNotFound
	}
	allTexts, err := queryService.textReadRepository.ListAll()
	if err != nil {
		return TextStatistics{}, err
	}

	return TextStatistics{
		SymbolStatistics: queryService.SymbolStatistics(text.Value()),
		UniqueStatistics: queryService.UniqueStatistic(text.Value(), allTexts),
	}, nil
}

func (queryService *statisticsQueryService) SymbolStatistics(text model.Text) SymbolStatistics {
	value := text.Value

	allSymbolsCount := utf8.RuneCountInString(value)
	alphabetSymbolsCount := 0

	for _, r := range value {
		if unicode.IsLetter(r) {
			alphabetSymbolsCount++
		}
	}

	return SymbolStatistics{
		AlphabetSymbolsCount: alphabetSymbolsCount,
		AllSymbolsCount:      allSymbolsCount,
	}
}

func (queryService *statisticsQueryService) UniqueStatistic(targetText model.Text, allTexts []model.Text) UniqueStatistics {
	for _, otherText := range allTexts {
		if otherText.ID == targetText.ID {
			continue
		}
		if otherText.Value == targetText.Value {
			return UniqueStatistics{
				IsDuplicate: true,
			}
		}
	}
	return UniqueStatistics{
		IsDuplicate: false,
	}
}
