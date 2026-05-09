package handler

import (
	"encoding/json"
	"net/http"
	"qarzi/internal/model"
	"qarzi/internal/repository"
)

type ClientHandler struct {
	Repo *repository.ClientRepo
}

func NewClientHandler(repo *repository.ClientRepo) *ClientHandler {
	return &ClientHandler{Repo: repo}
}

func (h *ClientHandler) CreateClient(w http.ResponseWriter, r *http.Request) {
	var req struct {
		BusinessID int    `json:"business_id"`
		Name       string `json:"name"`
		Phone      string `json:"phone"`
		DebtLimit  int64  `json:"debt_limit"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Неверный формат данных"}`, http.StatusBadRequest)
		return
	}

	client := &model.Client{
		BusinessID: req.BusinessID,
		Name:       req.Name,
		Phone:      req.Phone,
		DebtLimit:  req.DebtLimit,
	}

	if err := h.Repo.CreateClient(client); err != nil {
		http.Error(w, `{"error": "Ошибка при создании клиента (возможно, телефон уже существует в этом бизнесе)"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Клиент успешно создан",
		"client":  client.ID,
	})
}
