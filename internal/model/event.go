package model

import "time"

type Event struct {
	ID          int       `json:"id"`
	ClientID    int       `json:"client_id"`
	EventType   string    `json:"event_type"`
	Amount      int64     `json:"amount"`
	Status      string    `json:"status"`
	Description string    `json:"Description"`
	CreatedBy   *int      `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"update_at"`
}
