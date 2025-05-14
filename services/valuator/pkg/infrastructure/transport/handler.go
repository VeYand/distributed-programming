package transport

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	repo "valuator/pkg/infrastructure/redis/repository"

	"valuator/pkg/app/query"
	"valuator/pkg/app/service"
)

type Handler struct {
	textService      service.TextService
	textQueryService query.TextQueryService
}

func NewHandler(textService service.TextService, textQueryService query.TextQueryService) *Handler {
	return &Handler{
		textService:      textService,
		textQueryService: textQueryService,
	}
}

func (h *Handler) GetAddForm(w http.ResponseWriter, _ *http.Request) {
	tmpl, err := template.ParseFiles("./data/html/add-form.html")
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) CalculateStatistics(w http.ResponseWriter, r *http.Request) {
	id, err := h.textService.Add(r.FormValue("region"), r.FormValue("text"))
	if err != nil {
		log.Printf("Error adding text: %v", err)
		if errors.Is(err, repo.ErrInvalidRegion) {
			http.Error(w, "Invalid region", http.StatusBadRequest)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	dataStruct := struct {
		StatisticsURL string `json:"statistics_url"`
	}{
		StatisticsURL: fmt.Sprintf("/statistics/%s", id),
	}

	w.Header().Set("Content-Type", "application/json")

	if err = json.NewEncoder(w).Encode(dataStruct); err != nil {
		log.Printf("Error encoding JSON: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
