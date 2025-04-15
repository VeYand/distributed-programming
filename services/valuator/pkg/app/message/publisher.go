package message

import "valuator/pkg/app/model"

type Message struct {
	Type string
	Data interface{}
}

type addedTextSerializable struct {
	ID    string `json:"text_id"`
	Value string `json:"value"`
	Count int    `json:"count"`
}

func NewTextAddedMessage(text model.Text) Message {
	return Message{
		Type: "TextAdded",
		Data: addedTextSerializable{
			ID:    string(text.ID),
			Value: text.Value,
			Count: text.Count,
		},
	}
}

type Publisher interface {
	Publish(event Message) error
}
