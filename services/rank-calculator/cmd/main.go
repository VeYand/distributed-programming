package main

import (
	"github.com/gorilla/mux"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
	"rankcalculator/pkg/app/event"
	"rankcalculator/pkg/app/query"
	"rankcalculator/pkg/app/service"
	amqp2 "rankcalculator/pkg/infrastructure/amqp"
	"rankcalculator/pkg/infrastructure/redis/repository"
	"rankcalculator/pkg/infrastructure/transport"
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

func setupRoutes(handler *transport.Handler) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/statistics/{id}", handler.GetStatisticsPage).Methods("GET")
	router.HandleFunc("/statistics/{id}", handler.GetStatisticsAPI).Methods("POST")

	return router
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

	statisticsQueryService := query.NewStatisticsQueryService(statisticsRepository)
	httpHandler := transport.NewHandler(statisticsQueryService)
	router := setupRoutes(httpHandler)

	err = rabbitMQConsumer.ConnectReadChannel()
	if err != nil {
		log.Fatalf("RabbitMQ initialization error: %s", err)
		return
	}

	log.Println("Server is listening on port 8082")
	if err := http.ListenAndServe(":8082", router); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}
