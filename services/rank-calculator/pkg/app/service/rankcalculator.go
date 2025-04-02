package service

import (
	"rankcalculator/pkg/app/data"
	"rankcalculator/pkg/app/model"
	"rankcalculator/pkg/app/repository"
	"unicode"
	"unicode/utf8"
)

type RankCalculator interface {
	Calculate(text data.Text) error
}

func NewRankCalculator(
	statisticsRepository repository.StatisticsRepository,
) RankCalculator {
	return &rankCalculator{
		statisticsRepository: statisticsRepository,
	}
}

type rankCalculator struct {
	statisticsRepository repository.StatisticsRepository
}

func (r rankCalculator) Calculate(text data.Text) error {
	value := text.Value

	allSymbolsCount := utf8.RuneCountInString(value)
	alphabetSymbolsCount := 0

	for _, r := range value {
		if unicode.IsLetter(r) {
			alphabetSymbolsCount++
		}
	}

	isDuplicate := text.Count > 1

	return r.statisticsRepository.Store(model.Statistics{
		TextID:               text.ID,
		AllSymbolsCount:      allSymbolsCount,
		AlphabetSymbolsCount: alphabetSymbolsCount,
		IsDuplicate:          isDuplicate,
	})
}
