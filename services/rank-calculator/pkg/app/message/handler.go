package message

import (
	"context"
	"encoding/json"
	"rankcalculator/pkg/app/data"
	"rankcalculator/pkg/app/model"
	"rankcalculator/pkg/app/service"
)

type Message struct {
	Type string
	Data []byte
}

type TextAddedPayload struct {
	ID    string `json:"text_id"`
	Value string `json:"value"`
	Count int    `json:"count"`
}

type Handler interface {
	Handle(ctx context.Context, message Message) error
}

func NewHandler(rankCalculator service.RankCalculator) Handler {
	return &handler{
		rankCalculator: rankCalculator,
	}
}

type handler struct {
	rankCalculator service.RankCalculator
}

func (h *handler) Handle(_ context.Context, message Message) error {
	if message.Type == "TextAdded" {
		var payload TextAddedPayload
		err := json.Unmarshal(message.Data, &payload)
		if err != nil {
			return err
		}

		err = h.rankCalculator.Calculate(data.Text{
			ID:    model.TextID(payload.ID),
			Value: payload.Value,
			Count: payload.Count,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
