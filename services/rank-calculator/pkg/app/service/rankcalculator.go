package service

import (
	"rankcalculator/pkg/app/data"
	"rankcalculator/pkg/app/event"
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
	eventDispatcher event.Dispatcher,
) RankCalculator {
	return &rankCalculator{
		statisticsRepository: statisticsRepository,
		eventDispatcher:      eventDispatcher,
	}
}

type rankCalculator struct {
	statisticsRepository repository.StatisticsRepository
	eventDispatcher      event.Dispatcher
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

	err := r.statisticsRepository.Store(model.Statistics{
		TextID:               text.ID,
		AllSymbolsCount:      allSymbolsCount,
		AlphabetSymbolsCount: alphabetSymbolsCount,
		IsDuplicate:          isDuplicate,
	})
	if err != nil {
		return err
	}

	return r.eventDispatcher.Dispatch(event.CreateRankCalculatedEvent(text.ID))
}
