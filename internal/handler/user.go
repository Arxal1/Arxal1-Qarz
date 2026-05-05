package handler

import (
	"encoding/json"
	"net/http"
	"qarzi/internal/model"
	"qarzi/internal/repository"
)

type UserHandler struct {
	Repo *repository.UserRepo
}

func NewUserHandler(repo *repository.UserRepo) *UserHandler {
	return &UserHandler{Repo: repo}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		TelegramID int64  `json:"telegram_id"`
		FullName   string `json:"full_name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Неверный формат данных"}`, http.StatusBadRequest)
		return
	}

	if req.TelegramID <= 0 {
		http.Error(w, `{"error": "Invalid telegram_id"}`, http.StatusBadRequest)
		return
	}

	if req.FullName == "" {
		http.Error(w, `{"error": "full_name is required"}`, http.StatusBadRequest)
		return
	}

	user := &model.User{
		TelegramID: req.TelegramID,
		FullName:   &req.FullName,
	}

	if err := h.Repo.CreateUser(user); err != nil {
		http.Error(w, `{"error": "Ошибка при создании пользователя"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}
