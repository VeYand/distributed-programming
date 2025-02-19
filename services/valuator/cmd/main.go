package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"

	"valuator/pkg/app/query"
	"valuator/pkg/app/service"
	"valuator/pkg/infrastructure/redis/repository"
	"valuator/pkg/infrastructure/transport"
)

func initRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})
}

func initHandler(rdb *redis.Client) (*transport.Handler, error) {
	textRepo := repo.NewTextRepository(rdb)
	textService := service.NewTextService(textRepo)
	statisticsQueryService := query.NewStatisticsQueryService(textRepo)
	textQueryService := query.NewTextQueryService(textRepo)

	return transport.NewHandler(textService, statisticsQueryService, textQueryService), nil
}

func setupRoutes(handler *transport.Handler) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/add-form", handler.GetAddForm).Methods("GET")

	router.HandleFunc("/add", handler.AddText).Methods("POST")
	router.HandleFunc("/summary", handler.GetStatistics).Methods("GET")
	router.HandleFunc("/delete", handler.DeleteText).Methods("POST")
	router.HandleFunc("/list", handler.ListTexts).Methods("GET")
	router.HandleFunc("/", handler.ListTexts).Methods("GET")

	return router
}

func main() {
	rdb := initRedis()
	handler, err := initHandler(rdb)
	if err != nil {
		log.Fatalf("Could not initialize services: %v", err)
	}

	router := setupRoutes(handler)

	log.Println("Server is listening on port 8082...")
	if err := http.ListenAndServe(":8082", router); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}
