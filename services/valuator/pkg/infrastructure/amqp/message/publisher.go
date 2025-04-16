package message

import (
	"encoding/json"
	"valuator/pkg/app/message"

	amqp "github.com/rabbitmq/amqp091-go"
	amqpinf "valuator/pkg/infrastructure/amqp"
)

func NewMessagePublisher(
	publisher *amqpinf.RabbitMQClient,
	queueName string,
) message.Publisher {
	return &messagePublisher{publisher: publisher, queueName: queueName}
}

type messagePublisher struct {
	publisher *amqpinf.RabbitMQClient
	queueName string
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
		"",           // exchange (пустой exchange для direct отправки в очередь)
		ed.queueName, // routing key - имя очереди, куда должно быть отправлено сообщение.
		false,        // mandatory - если true, то сервер вернёт сообщение обратно, если оно не может быть маршрутизировано; false — не требуется возврат.
		false,        // immediate - если true, то сообщение будет доставлено немедленно (иначе ошибка, если нет готового потребителя); false — без проверки моментальной доставки.
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	return err
}
