package query

import (
	"rankcalculator/pkg/app/calculator"
	"rankcalculator/pkg/app/data"
	"rankcalculator/pkg/app/errors"
	"rankcalculator/pkg/app/model"
	"rankcalculator/pkg/app/repository"
)

type StatisticsQueryService interface {
	Get(ID string) (data.Statistics, error)
}

func NewStatisticsQueryService(statisticsReadRepository repository.StatisticsReadRepository) StatisticsQueryService {
	return &statisticsQueryService{
		statisticsReadRepository: statisticsReadRepository,
	}
}

type statisticsQueryService struct {
	statisticsReadRepository repository.StatisticsReadRepository
}

func (s *statisticsQueryService) Get(ID string) (data.Statistics, error) {
	stat, err := s.statisticsReadRepository.Find(model.TextID(ID))
	if err != nil {
		return data.Statistics{}, err
	}

	if stat.IsEmpty() {
		return data.Statistics{}, errors.ErrStatisticsNotFound
	}

	statisticsValue := stat.Value()

	return data.Statistics{
		TextID:               string(statisticsValue.TextID),
		AlphabetSymbolsCount: statisticsValue.AlphabetSymbolsCount,
		AllSymbolsCount:      statisticsValue.AllSymbolsCount,
		IsDuplicate:          statisticsValue.IsDuplicate,
		Rank:                 calculator.NewRankCalculator().Calculate(statisticsValue),
	}, nil
}
