package main

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"log"
	"rankcalculator/pkg/app/event"
	"rankcalculator/pkg/app/service"
	amqp2 "rankcalculator/pkg/infrastructure/amqp"
	"rankcalculator/pkg/infrastructure/redis/repository"
)

func createRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "12345Q",
	})
}

func newRabbitMQClient() (*amqp2.RabbitMQClient, error) {
	amqpURL := "amqp://guest:guest@rabbitmq:5672/"
	queueName := "valuator_queue"
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &amqp2.RabbitMQClient{
		Conn:      conn,
		Channel:   ch,
		QueueName: queueName,
	}, nil
}

func main() {
	rabbitMQClient, err := newRabbitMQClient()
	if err != nil {
		log.Fatalf("RabbitMQ initialization error: %s", err)
	}
	defer rabbitMQClient.Close()

	statisticsRepository := repo.NewStatisticsRepository(createRedisClient())
	rankCalculator := service.NewRankCalculator(statisticsRepository)
	eventHandler := event.NewHandler(rankCalculator)
	rabbitHandler := amqp2.NewHandler(eventHandler)
	rabbitMQConsumer := amqp2.NewRabbitMQConsumer(rabbitMQClient, rabbitHandler)

	err = rabbitMQConsumer.ConnectReadChannel()
	if err != nil {
		log.Fatalf("RabbitMQ initialization error: %s", err)
		return
	}

	fmt.Println("rankcalculator service is running")

	select {}
}
