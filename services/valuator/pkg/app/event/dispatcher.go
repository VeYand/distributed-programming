package event

import "valuator/pkg/app/model"

type Event struct {
	Type string
	Data interface{}
}

type addedTextSerializable struct {
	ID    string `json:"text_id"`
	Value string `json:"value"`
	Count int    `json:"count"`
}

func NewTextAddedEvent(text model.Text) Event {
	return Event{
		Type: "TextAdded",
		Data: addedTextSerializable{
			ID:    string(text.ID),
			Value: text.Value,
			Count: text.Count,
		},
	}
}

type Dispatcher interface {
	Dispatch(event Event) error
}
