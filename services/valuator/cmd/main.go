package main

import (
	"log"
	"net/http"
	amqp2 "valuator/pkg/infrastructure/amqp"

	"github.com/gorilla/mux"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"

	"valuator/pkg/app/query"
	"valuator/pkg/app/service"
	"valuator/pkg/infrastructure/redis/repository"
	"valuator/pkg/infrastructure/transport"
)

func createRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "12345Q",
	})
}

func newRabbitMQClient() (*amqp2.RabbitMQClient, error) {
	amqpURL := "amqp://guest:guest@rabbitmq:5672/"
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
		QueueName: "valuator_queue",
	}, nil
}

func createHandler(rdb *redis.Client, rabbitMQClient *amqp2.RabbitMQClient) *transport.Handler {
	publisher := amqp2.NewMessagePublisher(rabbitMQClient)
	textRepo := repo.NewTextRepository(rdb)
	textService := service.NewTextService(textRepo, publisher)
	textQueryService := query.NewTextQueryService(textRepo)

	return transport.NewHandler(textService, textQueryService)
}

func setupRoutes(handler *transport.Handler) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/text/calculate", handler.CalculateStatistics).Methods("POST")

	router.HandleFunc("/text/add-form", handler.GetAddForm).Methods("GET")
	router.HandleFunc("/", handler.GetAddForm).Methods("GET")

	return router
}

func main() {
	rabbitMQClient, err := newRabbitMQClient()
	if err != nil {
		log.Fatalf("Ошибка инициализации RabbitMQ: %s", err)
	}
	defer rabbitMQClient.Close()

	rdb := createRedisClient()
	handler := createHandler(rdb, rabbitMQClient)
	router := setupRoutes(handler)

	log.Println("Server is listening on port 8082")
	if err := http.ListenAndServe(":8082", router); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}

// TODO: Для чего нужен exchange
// TODO: Чем модель передачи сообщений отличается от событийной модели
// TODO: не использовать reader writer, использовать publisher consumer
// TODO: разобраться с флагом durable
// TODO: разобраться с флагом autoAck
