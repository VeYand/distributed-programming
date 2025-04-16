package message

import (
	"context"
	"log"
	amqpinf "rankcalculator/pkg/infrastructure/amqp"
)

type RabbitMQConsumer interface {
	ConnectReadChannel() error
}

func NewRabbitMQConsumer(
	client *amqpinf.RabbitMQClient,
	handler Handler,
	queueName string,
) RabbitMQConsumer {
	return &rabbitMQConsumer{
		client:    client,
		handler:   handler,
		queueName: queueName,
	}
}

type rabbitMQConsumer struct {
	client    *amqpinf.RabbitMQClient
	handler   Handler
	queueName string
}

func (r *rabbitMQConsumer) ConnectReadChannel() error {
	channel := r.client.Channel
	queueName := r.queueName

	q, err := channel.QueueDeclare(
		queueName, // name - имя очереди.
		false,     // durable - если true, очередь сохраняется при перезапуске брокера; false — очередь не сохраняется после рестарта.
		false,     // delete when unused - если true, очередь будет удалена, когда перестанут быть к ней привязанные потребители; false — не удаляется автоматически.
		false,     // exclusive - если true, очередь используется только текущим соединением и удаляется при его закрытии; false — очередь может использоваться несколькими соединениями.
		false,     // no-wait - если true, сервер не будет ждать подтверждения о создании очереди; false — клиент ждёт подтверждения от сервера.
		nil,       // arguments - дополнительные аргументы очереди (например, настройки TTL, лимиты и пр.); nil — дополнительных настроек нет.
	)
	if err != nil {
		return err
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

	go func() {
		for d := range msgs {
			err = r.handler.Handle(context.Background(), d.Body)
			if err != nil {
				log.Printf("Error handling message: %s", err)
				errAck := d.Nack(false, true)
				if errAck != nil {
					log.Printf("Error sending Nack: %s", errAck)
				}
			} else {
				errAck := d.Ack(false)
				if errAck != nil {
					log.Printf("Error sending Ack: %s", errAck)
				} else {
					log.Printf("Successfully handled message: %s", d.Body)
				}
			}
		}
	}()

	return nil
}
