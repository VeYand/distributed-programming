package transport

import (
	"encoding/json"
	"errors"
	"net/http"
	"protokey/pkg/app/service"
)

type Handler struct {
	service service.Service
}

func NewHandler(svc service.Service) *Handler {
	return &Handler{
		service: svc,
	}
}

func (handler *Handler) GetValue(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")

	value, err := handler.service.Get(key)
	if errors.Is(err, service.ErrBadRequest) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(
		struct {
			Value int `json:"value"`
		}{
			Value: value,
		},
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (handler *Handler) SetValue(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Key   string `json:"key"`
		Value int    `json:"value"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON payload: "+err.Error(), http.StatusBadRequest)
		return
	}

	err := handler.service.Set(req.Key, req.Value)
	if errors.Is(err, service.ErrBadRequest) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (handler *Handler) ListKeys(w http.ResponseWriter, r *http.Request) {
	prefix := r.URL.Query().Get("prefix")

	keys, err := handler.service.Keys(prefix)
	if errors.Is(err, service.ErrBadRequest) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(
		struct {
			Keys []string `json:"keys"`
		}{
			Keys: keys,
		},
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
