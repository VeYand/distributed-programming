package event

import (
	"valuator/pkg/app/model"
)

type Event struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

func CreateSimilarityCalculatedEvent(textID model.TextID, isDuplicate bool) Event {
	return Event{
		Type: "valuator.similarity_calculated",
		Data: struct {
			TextID      model.TextID `json:"text_id"`
			IsDuplicate bool         `json:"is_duplicate"`
		}{
			TextID:      textID,
			IsDuplicate: isDuplicate,
		},
	}
}

type Dispatcher interface {
	Dispatch(event Event) error
}
