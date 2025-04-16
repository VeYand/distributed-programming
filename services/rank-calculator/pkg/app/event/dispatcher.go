package event

import (
	"rankcalculator/pkg/app/model"
)

type Event struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

func CreateRankCalculatedEvent(statisticsID model.TextID) Event {
	return Event{
		Type: "rankcalculator.rank_calculated",
		Data: string(statisticsID),
	}
}

type Dispatcher interface {
	Dispatch(event Event) error
}
