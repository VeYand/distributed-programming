package message

import (
	"context"
	"encoding/json"
	"rankcalculator/pkg/app/message"
	"time"
)

type Handler interface {
	Handle(ctx context.Context, body []byte) error
}

func NewHandler(messageHandler message.Handler) Handler {
	return &handler{messageHandler: messageHandler}
}

type handler struct {
	messageHandler message.Handler
}

type messageSerializable struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

func (h *handler) Handle(ctx context.Context, body []byte) error {
	time.Sleep(time.Second * 5)

	var evt messageSerializable
	err := json.Unmarshal(body, &evt)
	if err != nil {
		return err
	}

	return h.messageHandler.Handle(ctx, message.Message{
		Type: evt.Type,
		Data: evt.Data,
	})
}
