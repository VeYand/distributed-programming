package event

import (
	"context"
	amqpinf "eventslogger/pkg/infrastructure/amqp"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

type RabbitMQSubscriber interface {
	Subscribe(routingKeys []string) error
}

func NewRabbitMQSubscriber(
	client *amqpinf.RabbitMQClient,
	handler Handler,
	queueName string,
) RabbitMQSubscriber {
	return &rabbitMQSubscriber{
		client:    client,
		handler:   handler,
		queueName: queueName,
	}
}

type rabbitMQSubscriber struct {
	client    *amqpinf.RabbitMQClient
	handler   Handler
	queueName string
}

func (r *rabbitMQSubscriber) Subscribe(routingKeys []string) error {
	channel := r.client.Channel
	queueName := r.queueName

	err := channel.ExchangeDeclare(
		"events", // name
		"topic",  // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		return err
	}

	q, err := channel.QueueDeclare(
		queueName, // name - имя очереди.
		true,      // durable - если true, очередь сохраняется при перезапуске брокера; false — очередь не сохраняется после рестарта.
		false,     // delete when unused - если true, очередь будет удалена, когда перестанут быть к ней привязанные потребители; false — не удаляется автоматически.
		false,     // exclusive - если true, очередь используется только текущим соединением и удаляется при его закрытии; false — очередь может использоваться несколькими соединениями.
		false,     // no-wait - если true, сервер не будет ждать подтверждения о создании очереди; false — клиент ждёт подтверждения от сервера.
		nil,       // arguments - дополнительные аргументы очереди (например, настройки TTL, лимиты и пр.); nil — дополнительных настроек нет.
	)
	if err != nil {
		return err
	}

	for _, routingKey := range routingKeys {
		err = channel.QueueBind(
			q.Name,     // queue name
			routingKey, // routing key
			"events",   // exchange
			false,
			nil,
		)
		if err != nil {
			return err
		}
	}

	msgs, err := channel.Consume(
		q.Name, // queue - имя очереди, из которой будет происходить чтение сообщений.
		"",     // consumer - тег потребителя; пустая строка позволяет серверу сгенерировать его автоматически.
		false,  // auto-ack - если true, то сообщения автоматически подтверждаются сразу после получения, без явного ack от потребителя.
		false,  // exclusive - если true, то только данный потребитель может получать сообщения из очереди; false — очередь может обслуживать нескольких потребителей.
		false,  // no-local - если true, сообщения, опубликованные текущим соединением, не будут доставляться этому же соединению; false — таких ограничений нет.
		false,  // no-wait - если true, клиент не ждёт подтверждения от сервера о начале потребления; false — ожидание подтверждения.
		nil,    // args - дополнительные аргументы для потребителя; nil — дополнительных настроек нет.
	)
	if err != nil {
		return err
	}

	go r.processMessages(msgs)
	return nil
}

func (r *rabbitMQSubscriber) processMessages(msgs <-chan amqp.Delivery) {
	for d := range msgs {
		err := r.handler.Handle(context.Background(), d.Body)
		if err != nil {
			log.Printf("Error handling message: %s", err)
			if err = d.Nack(false, true); err != nil {
				log.Printf("Error sending Nack: %s", err)
			}
		} else {
			if err = d.Ack(false); err != nil {
				log.Printf("Error sending Ack: %s", err)
			}
		}
	}
}
