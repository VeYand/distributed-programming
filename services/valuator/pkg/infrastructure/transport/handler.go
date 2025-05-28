package transport

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"valuator/pkg/app/auth"
	errors2 "valuator/pkg/app/errors"
	"valuator/pkg/app/model"
	"valuator/pkg/app/query"
	"valuator/pkg/app/service"
	"valuator/pkg/infrastructure/authentication"
	repo "valuator/pkg/infrastructure/redis/repository"
)

type Handler struct {
	textService       service.TextService
	textQueryService  query.TextQueryService
	permissionChecker auth.PermissionChecker
	authChecker       *authentication.Client
}

func NewHandler(
	textService service.TextService,
	textQueryService query.TextQueryService,
	authChecker *authentication.Client,
	permissionChecker auth.PermissionChecker,
) *Handler {
	return &Handler{
		textService:       textService,
		textQueryService:  textQueryService,
		authChecker:       authChecker,
		permissionChecker: permissionChecker,
	}
}

func (h *Handler) GetAddForm(w http.ResponseWriter, r *http.Request) {
	_, ok, err := h.authenticate(w, r)
	if err != nil {
		log.Printf("Error authenticating: %v", err)
		return
	}
	if !ok {
		return
	}

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
	userID, ok, err := h.authenticate(w, r)
	if err != nil {
		log.Printf("Error authenticating: %v", err)
		return
	}
	if !ok {
		return
	}

	id, err := h.textService.Add(userID, r.FormValue("region"), r.FormValue("text"))
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

func (h *Handler) CheckPermissions(w http.ResponseWriter, r *http.Request) {
	textID := r.URL.Query().Get("text_id")
	userID := r.URL.Query().Get("user_id")

	if textID == "" || userID == "" {
		log.Printf("Error checking permissions: text_id, user_id: %v", userID)
		http.Error(w, "Missing 'text_id' or 'user_id' query parameter", http.StatusBadRequest)
		return
	}

	canReadText, err := h.permissionChecker.CanReadText(userID, model.TextID(textID))
	if errors.Is(err, errors2.ErrTextNotFound) {
		log.Printf("Error text not found: %v", err)
		http.Error(w, "Text not found", http.StatusNotFound)
		return
	}
	if err != nil {
		log.Printf("Error checking permissions: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		CanReadText bool `json:"can_read_text"`
	}{
		CanReadText: canReadText,
	})
}

func (h *Handler) authenticate(w http.ResponseWriter, r *http.Request) (string, bool, error) {
	userID, ok, err := h.authChecker.IsAuthenticatedFromRequest(r)
	if err != nil {
		http.Error(w, "Auth service error: "+err.Error(), http.StatusInternalServerError)
		return "", false, err
	}
	if !ok {
		http.Redirect(w, r, "/user/signin", http.StatusSeeOther)
		return "", false, nil
	}
	return userID, ok, nil
}
