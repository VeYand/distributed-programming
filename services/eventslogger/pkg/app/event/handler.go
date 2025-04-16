package event

import (
	"context"
	"encoding/json"
	"log"
)

type Event struct {
	Type string
	Data []byte
}

type RankCalculatedEventPayload struct {
	StatisticID string  `json:"statistic_id"`
	Rank        float64 `json:"rank"`
}

type SimilarityCalculatedEventPayload struct {
	TextID      string `json:"text_id"`
	IsDuplicate bool   `json:"is_duplicate"`
}

type Handler interface {
	Handle(ctx context.Context, message Event) error
}

func NewHandler() Handler {
	return &handler{}
}

type handler struct{}

func (h *handler) Handle(_ context.Context, message Event) error {
	if message.Type == "rankcalculator.rank_calculated" {
		var payload RankCalculatedEventPayload
		err := json.Unmarshal(message.Data, &payload)
		if err != nil {
			return err
		}

		log.Println("Rank calculated, statistic ID: ", payload.StatisticID, ", Rank: ", payload.Rank)
	} else if message.Type == "valuator.similarity_calculated" {
		var payload SimilarityCalculatedEventPayload
		err := json.Unmarshal(message.Data, &payload)
		if err != nil {
			return err
		}

		log.Println("Similarity calculated, textID: ", payload.TextID, ", isDuplicate: ", payload.IsDuplicate)
	} else {
		log.Println("Unknown message type: ", message.Type)
	}

	return nil
}
