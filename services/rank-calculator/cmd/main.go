package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func setupRoutes() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/", helloWorld).Methods("GET")

	return router
}

func main() {
	router := setupRoutes()

	log.Println("Server is listening on port 8082")
	if err := http.ListenAndServe(":8082", router); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}

func helloWorld(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	if _, err := w.Write([]byte("hello world")); err != nil {
		log.Printf("Error writing response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
