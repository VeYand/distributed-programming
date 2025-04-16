package event

import (
	"context"
	"encoding/json"
	"eventslogger/pkg/app/event"
	"time"
)

type Handler interface {
	Handle(ctx context.Context, body []byte) error
}

func NewHandler(eventHandler event.Handler) Handler {
	return &handler{eventHandler: eventHandler}
}

type handler struct {
	eventHandler event.Handler
}

type eventSerializable struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

func (h *handler) Handle(ctx context.Context, body []byte) error {
	time.Sleep(time.Second * 5)

	var evt eventSerializable
	err := json.Unmarshal(body, &evt)
	if err != nil {
		return err
	}

	return h.eventHandler.Handle(ctx, event.Event{
		Type: evt.Type,
		Data: evt.Data,
	})
}
