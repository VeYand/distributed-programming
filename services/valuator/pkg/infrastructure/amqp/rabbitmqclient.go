package amqp

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

type RabbitMQClient struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
}

func (r *RabbitMQClient) Close() {
	if err := r.Channel.Close(); err != nil {
		log.Printf("Error closing Channel: %s", err)
	}
	if err := r.Conn.Close(); err != nil {
		log.Printf("Error closing connection: %s", err)
	}
}
