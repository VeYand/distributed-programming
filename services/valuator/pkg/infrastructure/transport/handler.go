package transport

import (
	"encoding/json"
	stderrors "errors"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"

	"valuator/pkg/app/errors"
	"valuator/pkg/app/query"
	"valuator/pkg/app/service"
	"valuator/pkg/app/statistics"
)

type Handler struct {
	textService       service.TextService
	statisticsService statistics.TextStatistics
	textQueryService  query.TextQueryService
}

func NewHandler(textService service.TextService, statisticsQueryService statistics.TextStatistics, textQueryService query.TextQueryService) *Handler {
	return &Handler{
		textService:       textService,
		statisticsService: statisticsQueryService,
		textQueryService:  textQueryService,
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
	id, err := h.textService.Add(r.FormValue("text"))
	if err != nil {
		log.Printf("Error adding text: %v", err)
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

func (h *Handler) GetStatistics(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	summary, err := h.statisticsService.GetSummary(id)
	if stderrors.Is(err, errors.ErrTextNotFound) {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if err != nil {
		log.Printf("Error getting summary: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	rank := 1 - float64(summary.SymbolStatistics.AlphabetSymbolsCount)/float64(summary.SymbolStatistics.AllSymbolsCount) //todo: вынести логику в app
	similarity := 0
	if summary.UniqueStatistics.IsDuplicate {
		similarity = 1
	}

	dataStruct := struct {
		Title      string
		TextID     string
		Rank       float64
		Similarity int
	}{
		Title:      "Результаты",
		TextID:     id,
		Rank:       rank,
		Similarity: similarity,
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
