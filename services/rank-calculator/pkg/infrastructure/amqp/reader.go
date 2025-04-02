package amqp

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

type RabbitMQClient struct {
	Conn      *amqp.Connection
	Channel   *amqp.Channel
	QueueName string
}

type RabbitMQConsumer interface {
	ConnectReadChannel() error
}

func NewRabbitMQConsumer(
	client *RabbitMQClient,
	handler Handler,
) RabbitMQConsumer {
	return &rabbitMQConsumer{
		client:  client,
		handler: handler,
	}
}

type rabbitMQConsumer struct {
	client  *RabbitMQClient
	handler Handler
}

func (r *RabbitMQClient) Close() {
	if err := r.Channel.Close(); err != nil {
		log.Printf("Error closing channel: %s", err)
	}
	if err := r.Conn.Close(); err != nil {
		log.Printf("Error closing connection: %s", err)
	}
}

func (r *rabbitMQConsumer) ConnectReadChannel() error {
	channel := r.client.Channel
	queueName := r.client.QueueName

	q, err := channel.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete
		// when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return err
	}

	msgs, err := channel.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return err
	}

	go func() {
		for d := range msgs {
			err = r.handler.Handle(context.Background(), d.Body)
			if err != nil {
				log.Printf("Error handling message: %s", err)
			} else {
				log.Printf("Successfully handled message: %s", d.Body)
			}
		}
	}()

	return err
}
