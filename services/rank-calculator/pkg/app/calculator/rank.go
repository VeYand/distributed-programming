package calculator

import "rankcalculator/pkg/app/model"

type RankCalculator struct{}

func NewRankCalculator() *RankCalculator {
	return &RankCalculator{}
}

func (rc *RankCalculator) Calculate(statistics model.Statistics) float64 {
	return 1 - float64(statistics.AlphabetSymbolsCount)/float64(statistics.AllSymbolsCount)
}
