package service

import (
	"log"
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
	centrifugoClient CentrifugoClient,
) RankCalculator {
	return &rankCalculator{
		statisticsRepository: statisticsRepository,
		eventDispatcher:      eventDispatcher,
		centrifugoClient:     centrifugoClient,
	}
}

type rankCalculator struct {
	statisticsRepository repository.StatisticsRepository
	eventDispatcher      event.Dispatcher
	centrifugoClient     CentrifugoClient
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

	statistics := model.Statistics{
		TextID:               text.ID,
		AllSymbolsCount:      allSymbolsCount,
		AlphabetSymbolsCount: alphabetSymbolsCount,
		IsDuplicate:          isDuplicate,
	}
	err := r.statisticsRepository.Store(statistics)
	if err != nil {
		return err
	}

	channel := "statistics#" + string(text.ID)
	err = r.centrifugoClient.Publish(
		channel,
		map[string]interface{}{
			"text_id": string(text.ID),
		},
	)
	if err != nil {
		log.Printf("Failed to publish centrifugo event: %v", err)
	}

	return r.eventDispatcher.Dispatch(event.CreateRankCalculatedEvent(text.ID, CalculateRank(statistics)))
}

func CalculateRank(statistics model.Statistics) float64 {
	if statistics.AllSymbolsCount == 0 {
		return 0
	}
	return 1 - float64(statistics.AlphabetSymbolsCount)/float64(statistics.AllSymbolsCount)
}
