package main

import (
	"github.com/gorilla/mux"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
	"os"
	appmessage "rankcalculator/pkg/app/message"
	"rankcalculator/pkg/app/query"
	"rankcalculator/pkg/app/service"
	amqpinf "rankcalculator/pkg/infrastructure/amqp"
	"rankcalculator/pkg/infrastructure/amqp/event"
	"rankcalculator/pkg/infrastructure/amqp/message"
	"rankcalculator/pkg/infrastructure/centrifugo"
	"rankcalculator/pkg/infrastructure/redis/repository"
	"rankcalculator/pkg/infrastructure/transport"
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

	mainRdb := redis.NewClient(&redis.Options{Addr: os.Getenv("REDIS_MAIN_URL")})
	shards := map[string]*redis.Client{
		"RU":   redis.NewClient(&redis.Options{Addr: os.Getenv("REDIS_RU_URL")}),
		"EU":   redis.NewClient(&redis.Options{Addr: os.Getenv("REDIS_EU_URL")}),
		"ASIA": redis.NewClient(&redis.Options{Addr: os.Getenv("REDIS_ASIA_URL")}),
	}

	eventDispatcher := event.NewEventDispatcher(rabbitMQClient, "events", "rankcalculator")
	shardManager := repo.NewShardManager(mainRdb, shards)
	statisticsRepository := repo.NewStatisticsShardedRepository(shardManager)
	centrifugoClient := centrifugo.NewClient(
		"http://centrifugo:8000/api/publish",
		"_salt",
	)
	rankCalculator := service.NewRankCalculator(statisticsRepository, eventDispatcher, centrifugoClient)
	messageHandler := appmessage.NewHandler(rankCalculator)
	rabbitHandler := message.NewHandler(messageHandler)
	rabbitMQConsumer := message.NewRabbitMQConsumer(rabbitMQClient, rabbitHandler, "valuator_queue")

	statisticsQueryService := query.NewStatisticsQueryService(statisticsRepository)
	httpHandler := transport.NewHandler(statisticsQueryService)
	router := setupRoutes(httpHandler)

	err = rabbitMQConsumer.Subscribe()
	if err != nil {
		log.Fatalf("RabbitMQ initialization error: %s", err)
		return
	}

	log.Println("Server is listening on port 8082")
	if err := http.ListenAndServe(":8082", router); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}
