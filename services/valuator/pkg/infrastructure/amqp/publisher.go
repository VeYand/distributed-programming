package amqp

import (
	"encoding/json"
	"log"
	"valuator/pkg/app/message"

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

func NewMessagePublisher(publisher *RabbitMQClient) message.Publisher {
	return &messagePublisher{publisher: publisher}
}

type messagePublisher struct {
	publisher *RabbitMQClient
}

type messageSerializable struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

func (ed *messagePublisher) Publish(message message.Message) error {
	body, err := json.Marshal(messageSerializable{
		Type: message.Type,
		Data: message.Data,
	})
	if err != nil {
		return err
	}
	err = ed.publisher.Channel.Publish(
		"",                     // exchange (пустой exchange для direct отправки в очередь)
		ed.publisher.QueueName, // routing key - имя очереди, куда должно быть отправлено сообщение.
		false,                  // mandatory - если true, то сервер вернёт сообщение обратно, если оно не может быть маршрутизировано; false — не требуется возврат.
		false,                  // immediate - если true, то сообщение будет доставлено немедленно (иначе ошибка, если нет готового потребителя); false — без проверки моментальной доставки.
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	return err
}
