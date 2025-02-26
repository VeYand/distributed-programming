package statistics

import (
	"github.com/gofrs/uuid"
	"unicode"
	"unicode/utf8"
	"valuator/pkg/app/data"
	"valuator/pkg/app/query"
)

type StatisticsQueryService interface {
	GetSummary(textID uuid.UUID) (TextStatistics, error)
}

func NewStatisticsQueryService(textQueryService query.TextQueryService) StatisticsQueryService {
	return &statisticsQueryService{
		textQueryService: textQueryService,
	}
}

type statisticsQueryService struct {
	textQueryService query.TextQueryService
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
	text, err := queryService.textQueryService.Get(textID)
	if err != nil {
		return TextStatistics{}, err
	}
	allTexts, err := queryService.textQueryService.List()
	if err != nil {
		return TextStatistics{}, err
	}

	return TextStatistics{
		SymbolStatistics: queryService.SymbolStatistics(text),
		UniqueStatistics: queryService.UniqueStatistic(text, allTexts),
	}, nil
}

func (queryService *statisticsQueryService) SymbolStatistics(text data.TextData) SymbolStatistics {
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

func (queryService *statisticsQueryService) UniqueStatistic(targetText data.TextData, allTexts []data.TextData) UniqueStatistics {
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
