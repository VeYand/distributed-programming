package event

import (
	"rankcalculator/pkg/app/model"
)

type Event struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

func CreateRankCalculatedEvent(statisticsID model.TextID, rank float64) Event {
	return Event{
		Type: "rankcalculator.rank_calculated",
		Data: struct {
			StatisticID model.TextID `json:"statistic_id"`
			Rank        float64      `json:"rank"`
		}{
			StatisticID: statisticsID,
			Rank:        rank,
		},
	}
}

type Dispatcher interface {
	Dispatch(event Event) error
}
