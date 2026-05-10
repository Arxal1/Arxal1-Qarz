package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"qarzi/internal/model"
	"qarzi/internal/repository"

	"github.com/go-chi/chi/v5"
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

func (h *EventHandler) GetClientStatement(w http.ResponseWriter, r *http.Request) {
	clientIDStr := chi.URLParam(r, "clientID")

	clientID, err := strconv.Atoi(clientIDStr)
	if err != nil {
		http.Error(w, `{"error": "Неверный ID клиента"}`, http.StatusBadRequest)
		return
	}

	events, err := h.Repo.GetClientEvents(clientID)
	if err != nil {
		http.Error(w, `{"error": "Ошибка при получении выписки"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"events": events,
	})
}

func (h *EventHandler) InitiatePayment(w http.ResponseWriter, r *http.Request) {
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

	if err := h.Repo.InitiatePayment(event); err != nil {
		http.Error(w, `{"error": "Ошибка при создании запроса на оплату"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":   "pending",
		"message":  "Запрос на оплату отправлен менеджеру",
		"event_id": event.ID,
	})
}

func (h *EventHandler) ConfirmPayment(w http.ResponseWriter, r *http.Request) {

	eventIDStr := chi.URLParam(r, "eventID")
	eventID, err := strconv.Atoi(eventIDStr)
	if err != nil {
		http.Error(w, `{"error": "Неверный ID события"}`, http.StatusBadRequest)
		return
	}

	if err := h.Repo.ConfirmPayment(eventID); err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%v"}`, err), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Оплата подтверждена, баланс обновлен",
	})
}

func (h *EventHandler) DisputeEvent(w http.ResponseWriter, r *http.Request) {

	eventIDStr := chi.URLParam(r, "eventID")
	eventID, err := strconv.Atoi(eventIDStr)
	if err != nil {
		http.Error(w, `{"error": "Неверный ID события"}`, http.StatusBadRequest)
		return
	}

	var req struct {
		ClientID int `json:"client_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Неверный формат данных"}`, http.StatusBadRequest)
		return
	}

	if err := h.Repo.DisputeEvent(eventID, req.ClientID); err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%v"}`, err), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Операция оспорена. Владелец бизнеса получит уведомление.",
	})
}

func (h *EventHandler) AdjustEvent(w http.ResponseWriter, r *http.Request) {
	eventIDStr := chi.URLParam(r, "eventID")
	eventID, err := strconv.Atoi(eventIDStr)
	if err != nil {
		http.Error(w, `{"error": "Неверный ID события"}`, http.StatusBadRequest)
		return
	}

	var req struct {
		NewAmount int64  `json:"new_amount"`
		UserID    int    `json:"user_id"`
		Comment   string `json:"comment"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Неверный формат данных"}`, http.StatusBadRequest)
		return
	}

	if err := h.Repo.AdjustEvent(eventID, req.NewAmount, req.UserID, req.Comment); err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%v"}`, err), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Сумма скорректирована, изменения занесены в аудит",
	})
}
