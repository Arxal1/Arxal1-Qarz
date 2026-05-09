package handler

import (
	"encoding/json"
	"net/http"

	"qarzi/internal/model"
	"qarzi/internal/repository"
)

type EventHandler struct {
	Repo *repository.EventRepo
}

func NewEventHandler(repo *repository.EventRepo) *EventHandler {
	return &EventHandler{Repo: repo}
}

func (h *EventHandler) CreateShipment(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ClientID    int    `json:"client_id"`
		Amount      int64  `json:"amount"`
		Description string `json:"description"`
		CreatedBy   int    `json:"created_by"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Неверный формат данных"}`, http.StatusBadRequest)
		return
	}

	event := &model.Event{
		ClientID:    req.ClientID,
		Amount:      req.Amount,
		Description: req.Description,
		CreatedBy:   &req.CreatedBy,
	}

	if err := h.Repo.CreateShipment(event); err != nil {
		http.Error(w, `{"error": "Ошибка при оформлении долга"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":   "success",
		"message":  "Долг успешно оформлен",
		"event_id": event.ID,
	})
}
