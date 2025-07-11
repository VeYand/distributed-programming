package transport

import (
	"encoding/json"
	stderrors "errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"time"

	"rankcalculator/pkg/app/errors"
	"rankcalculator/pkg/app/query"
)

type Handler struct {
	statisticsQueryService query.StatisticsQueryService
}

func NewHandler(statisticsQueryService query.StatisticsQueryService) *Handler {
	return &Handler{
		statisticsQueryService: statisticsQueryService,
	}
}

func (h *Handler) GetStatisticsPage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	textID := vars["id"]
	channel := "statistics#" + string(textID)
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = r.RemoteAddr
	}
	dataStruct := struct {
		Title           string
		CentrifugoToken string
		CentrifugoURL   string
		Channel         string
		ProcessingID    string
	}{
		Title:           "Результаты",
		CentrifugoToken: generateCentrifugoToken(ip, channel),
		CentrifugoURL:   "ws://127.0.0.1:8000/connection/websocket",
		Channel:         channel,
		ProcessingID:    textID,
	}

	tmpl, err := template.ParseFiles("./data/html/summary.html")
	if err != nil {
		log.Printf("Error parsing summary template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, dataStruct)
	if err != nil {
		log.Printf("Error executing summary template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetStatisticsAPI(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	summary, err := h.statisticsQueryService.Get(id)
	if stderrors.Is(err, errors.ErrStatisticsNotFound) {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if err != nil {
		log.Printf("Error getting summary: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	dataStruct := struct {
		TextID      string  `json:"text_id"`
		Rank        float64 `json:"rank"`
		IsDuplicate bool    `json:"is_duplicate"`
	}{
		TextID:      id,
		Rank:        summary.Rank,
		IsDuplicate: summary.IsDuplicate,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(dataStruct); err != nil {
		log.Printf("Error encoding JSON: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func generateCentrifugoToken(identifier string, channel string) string {
	claims := jwt.MapClaims{
		"sub":      identifier,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
		"channels": []string{channel},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte("my_secret"))
	if err != nil {
		log.Printf("Ошибка генерации токена: %v", err)
		return ""
	}

	return signedToken
}
