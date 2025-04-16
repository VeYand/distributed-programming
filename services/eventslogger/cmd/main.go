package main

import (
	appevent "eventslogger/pkg/app/event"
	amqpinf "eventslogger/pkg/infrastructure/amqp"
	"eventslogger/pkg/infrastructure/amqp/event"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

func newRabbitMQClient() (*amqpinf.RabbitMQClient, error) {
	amqpURL := "amqp://guest:guest@rabbitmq:5672/"
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &amqpinf.RabbitMQClient{
		Conn:    conn,
		Channel: ch,
	}, nil
}

func main() {
	rabbitMQClient, err := newRabbitMQClient()
	if err != nil {
		log.Fatalf("RabbitMQ initialization error: %s", err)
	}
	defer rabbitMQClient.Close()

	eventHandler := appevent.NewHandler()
	rabbitHandler := event.NewHandler(eventHandler)
	rabbitMQSubscriber := event.NewRabbitMQSubscriber(rabbitMQClient, rabbitHandler, "events.eventslogger")

	err = rabbitMQSubscriber.Subscribe([]string{"rankcalculator", "valuator"})
	if err != nil {
		log.Fatalf("RabbitMQ initialization error: %s", err)
		return
	}

	log.Println("Server initialized")
	select {}
}
