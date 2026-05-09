package model

import "time"

type Installment struct {
	ID           int       `json:"id"`
	EventID      int       `json:"event_id"`
	DueDate      time.Time `json:"due_date"`
	Amount       int64     `json:"amount"`
	Status       string    `json:"status"` // pending, paid, overdue
	ReminderSent bool      `json:"reminder_sent"`
	CreatedAt    time.Time `json:"created_at"`
}

type InstallmentWithClient struct {
	ID          int       `json:"id"`
	Amount      int64     `json:"amount"`
	DueDate     time.Time `json:"due_date"`
	ClientPhone string    `json:"client_phone"`
	ClientName  string    `json:"client_name"`

	TelegramID *int64 `json:"telegram_id"`
}
