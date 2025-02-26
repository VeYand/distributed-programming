package main

import (
	"log"
	"net/http"
	"valuator/pkg/app/statistics"

	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"

	"valuator/pkg/app/query"
	"valuator/pkg/app/service"
	"valuator/pkg/infrastructure/redis/repository"
	"valuator/pkg/infrastructure/transport"
)

func createRedisClient() *redis.Client { // todo: rename
	return redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "12345Q",
	})
}

func createHandler(rdb *redis.Client) *transport.Handler {
	textRepo := repo.NewTextRepository(rdb)
	textService := service.NewTextService(textRepo)
	textQueryService := query.NewTextQueryService(textRepo)
	statisticsQueryService := statistics.NewStatisticsQueryService(textQueryService)

	return transport.NewHandler(textService, statisticsQueryService, textQueryService)
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

func main() { // TODO хранить в редиске только одну копию текста, избавиться от получения значений в цикле
	rdb := createRedisClient()
	handler := createHandler(rdb)
	router := setupRoutes(handler)

	log.Println("Server is listening on port 8082")
	if err := http.ListenAndServe(":8082", router); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}
