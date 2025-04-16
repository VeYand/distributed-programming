package event

import (
	"encoding/json"
	"valuator/pkg/app/event"

	amqp "github.com/rabbitmq/amqp091-go"
	amqpinf "valuator/pkg/infrastructure/amqp"
)

func NewEventDispatcher(dispatcher *amqpinf.RabbitMQClient, exchangeName string, routingKey string) event.Dispatcher {
	return &messageDispatcher{
		dispatcher:   dispatcher,
		exchangeName: exchangeName,
		routingKey:   routingKey,
	}
}

type messageDispatcher struct {
	dispatcher   *amqpinf.RabbitMQClient
	exchangeName string
	routingKey   string
}

type eventSerializable struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

func (ed *messageDispatcher) Dispatch(event event.Event) error {
	body, err := json.Marshal(eventSerializable{
		Type: event.Type,
		Data: event.Data,
	})
	if err != nil {
		return err
	}

	err = ed.dispatcher.Channel.ExchangeDeclare(
		ed.exchangeName, // имя exchange
		"topic",         // тип exchange - fanout - широковещательная рассылка
		true,            // durable - exchange сохраняется даже после перезапуска брокера
		false,           // auto-deleted - exchange не удаляется автоматически, когда перестают быть привязанные очереди
		false,           // internal - false, чтобы exchange мог использоваться для публикации сообщений извне
		false,           // no-wait - false, клиент ждет подтверждения от сервера
		nil,             // args - дополнительные аргументы, здесь не используются
	)
	if err != nil {
		return err
	}

	return ed.dispatcher.Channel.Publish(
		ed.exchangeName, // exchange - имя ранее объявленного exchange
		ed.routingKey,   // routing key - имя очереди, куда должно быть отправлено сообщение - не требуется для fanout
		false,           // mandatory - если true, то сервер вернёт сообщение обратно, если оно не может быть маршрутизировано; false — не требуется возврат.
		false,           // immediate - если true, то сообщение будет доставлено немедленно (иначе ошибка, если нет готового потребителя); false — без проверки моментальной доставки.
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
}
