package handler

import (
	"encoding/json"
	"net/http"
	"qarzi/internal/model"
	"qarzi/internal/repository"
)

type BusinessHandler struct {
	Repo *repository.BusinessRepo
}

func NewBusinessHandler(repo *repository.BusinessRepo) *BusinessHandler {
	return &BusinessHandler{Repo: repo}
}

func (h *BusinessHandler) RegisterBusiness(w http.ResponseWriter, r *http.Request) {
	var req struct {
		OwnerID int64  `json:"owner_id"`
		Name    string `json:"name"`
		Phone   string `json:"phone"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Неверный формат данных"}`, http.StatusBadRequest)
		return
	}

	business := &model.Business{
		OwnerID: int(req.OwnerID),
		Name:    req.Name,
		Phone:   &req.Phone,
	}
	if err := h.Repo.CreateBusiness(business); err != nil {
		http.Error(w, `{"error": "Ошибка при создании бизнеса"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":   "success",
		"message":  "Бизнес успешно зарегистрирован",
		"business": business.ID,
	})
}
