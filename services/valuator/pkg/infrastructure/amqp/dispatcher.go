package amqp

import (
	"encoding/json"
	"log"
	"valuator/pkg/app/event"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQClient struct {
	Conn      *amqp.Connection
	Channel   *amqp.Channel
	QueueName string
}

func (r *RabbitMQClient) Close() {
	if err := r.Channel.Close(); err != nil {
		log.Printf("Error closing Channel: %s", err)
	}
	if err := r.Conn.Close(); err != nil {
		log.Printf("Error closing connection: %s", err)
	}
}

func NewEventDispatcher(publisher *RabbitMQClient) event.Dispatcher {
	return &eventDispatcher{publisher: publisher}
}

type eventDispatcher struct {
	publisher *RabbitMQClient
}

type eventSerializable struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

func (ed *eventDispatcher) Dispatch(event event.Event) error {
	body, err := json.Marshal(eventSerializable{
		Type: event.Type,
		Data: event.Data,
	})
	if err != nil {
		return err
	}
	err = ed.publisher.Channel.Publish(
		"",                     // exchange (пустой exchange для direct отправки в очередь)
		ed.publisher.QueueName, // routing key = имя очереди
		false,                  // mandatory
		false,                  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	return err
}
