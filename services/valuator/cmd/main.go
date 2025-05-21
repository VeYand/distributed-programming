package main

import (
	"log"
	"net/http"
	"os"
	amqp2 "valuator/pkg/infrastructure/amqp"
	"valuator/pkg/infrastructure/authentication"

	"valuator/pkg/infrastructure/amqp/event"
	"valuator/pkg/infrastructure/amqp/message"

	"github.com/gorilla/mux"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"

	"valuator/pkg/app/query"
	"valuator/pkg/app/service"
	"valuator/pkg/infrastructure/redis/repository"
	"valuator/pkg/infrastructure/transport"
)

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
		Conn:    conn,
		Channel: ch,
	}, nil
}

func createHandler(rabbitMQClient *amqp2.RabbitMQClient) *transport.Handler {
	mainRdb := redis.NewClient(&redis.Options{Addr: os.Getenv("REDIS_MAIN_URL")})
	shards := map[string]*redis.Client{
		"RU":   redis.NewClient(&redis.Options{Addr: os.Getenv("REDIS_RU_URL")}),
		"EU":   redis.NewClient(&redis.Options{Addr: os.Getenv("REDIS_EU_URL")}),
		"ASIA": redis.NewClient(&redis.Options{Addr: os.Getenv("REDIS_ASIA_URL")}),
	}

	publisher := message.NewMessagePublisher(rabbitMQClient, "valuator_queue")
	dispatcher := event.NewEventDispatcher(rabbitMQClient, "events", "valuator")
	shardManager := repo.NewShardManager(mainRdb, shards)
	textRepo := repo.NewTextShardedRepository(shardManager)
	textService := service.NewTextService(textRepo, publisher, dispatcher)
	textQueryService := query.NewTextQueryService(textRepo)
	authChecker := authentication.NewClient(os.Getenv("USER_INTERNAL_URL"))

	return transport.NewHandler(textService, textQueryService, authChecker)
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

	handler := createHandler(rabbitMQClient)
	router := setupRoutes(handler)

	log.Println("Server is listening on port 8082")
	if err := http.ListenAndServe(":8082", router); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}
