package transport

import (
	stderrors "errors"
	"html/template"
	"log"
	"net/http"

	"github.com/gofrs/uuid"

	"valuator/pkg/app/errors"
	"valuator/pkg/app/query"
	"valuator/pkg/app/service"
)

type Handler struct {
	textService       service.TextService
	statisticsService query.StatisticsQueryService
	textQueryService  query.TextQueryService
}

func NewHandler(textService service.TextService, statisticsQueryService query.StatisticsQueryService, textQueryService query.TextQueryService) *Handler {
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

func (h *Handler) AddText(w http.ResponseWriter, r *http.Request) {
	_, err := h.textService.Add(r.FormValue("text"))
	if err != nil {
		log.Printf("Error adding text: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetStatistics(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.FromString(r.FormValue("id"))
	if err != nil {
		log.Printf("Invalid UUID: %v", err)
		http.Error(w, "Bad identifier", http.StatusBadRequest)
		return
	}

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

	rank := float64(summary.SymbolStatistics.AlphabetSymbolsCount) / float64(summary.SymbolStatistics.AllSymbolsCount)
	similarity := 0
	if summary.UniqueStatistics.IsDuplicate {
		similarity = 1
	}

	data := struct {
		Title      string
		TextID     uuid.UUID
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
	err = tmpl.Execute(w, data)
	if err != nil {
		log.Printf("Error executing summary template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) DeleteText(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.FromString(r.FormValue("id"))
	if err != nil {
		log.Printf("Invalid UUID for deletion: %v", err)
		http.Error(w, "Bad identifier", http.StatusBadRequest)
		return
	}

	err = h.textService.Remove(id)
	if err != nil {
		log.Printf("Error deleting text: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) ListTexts(w http.ResponseWriter, _ *http.Request) {
	texts, err := h.textQueryService.List()
	if err != nil {
		log.Printf("Error listing texts: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("./data/html/list.html")
	if err != nil {
		log.Printf("Error parsing list template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, struct {
		Texts []query.TextData
	}{
		Texts: texts,
	})
	if err != nil {
		log.Printf("Error executing list template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
