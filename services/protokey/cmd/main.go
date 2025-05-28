package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"protokey/pkg/app/service"
	"protokey/pkg/infrastructure/storage"
	"protokey/pkg/infrastructure/transport"
)

func setupRoutes(handler *transport.Handler) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/get", handler.GetValue).Methods("GET")
	router.HandleFunc("/set", handler.SetValue).Methods("POST")
	router.HandleFunc("/keys", handler.ListKeys).Methods("GET")

	return router
}

func main() {
	dataFile, ok := os.LookupEnv("DATA_FILE")
	if !ok {
		dataFile = "data/protokey.data"
	}

	store := storage.NewStore(storage.Config{DataFile: dataFile})
	svc := service.NewProtoKeyService(store.CommandChan, store.ResponseChan)
	httpHandler := transport.NewHandler(svc)
	router := setupRoutes(httpHandler)

	log.Println("Server is listening on port 8082")
	if err := http.ListenAndServe(":8082", router); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}
