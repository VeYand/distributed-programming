package statistics

import (
	"unicode"
	"unicode/utf8"
	"valuator/pkg/app/data"
	"valuator/pkg/app/query"
)

type TextStatistics interface {
	GetSummary(textID string) (TextStatisticsData, error)
}

func NewStatisticsQueryService(textQueryService query.TextQueryService) TextStatistics {
	return &statisticsQueryService{
		textQueryService: textQueryService,
	}
}

type statisticsQueryService struct {
	textQueryService query.TextQueryService
}

type TextStatisticsData struct {
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

func (queryService *statisticsQueryService) GetSummary(textID string) (TextStatisticsData, error) {
	text, err := queryService.textQueryService.Get(textID)
	if err != nil {
		return TextStatisticsData{}, err
	}

	return TextStatisticsData{
		SymbolStatistics: queryService.SymbolStatistics(text),
		UniqueStatistics: queryService.UniqueStatistic(text),
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

func (queryService *statisticsQueryService) UniqueStatistic(targetText data.TextData) UniqueStatistics {
	return UniqueStatistics{
		IsDuplicate: targetText.Count > 1,
	}
}
