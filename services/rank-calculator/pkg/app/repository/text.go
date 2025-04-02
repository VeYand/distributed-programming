package repository

import (
	"github.com/mono83/maybe"
	"rankcalculator/pkg/app/model"
)

type StatisticsReadRepository interface {
	Find(ID model.TextID) (maybe.Maybe[model.Statistics], error)
}

type StatisticsRepository interface {
	StatisticsReadRepository
	Store(text model.Statistics) error
}
