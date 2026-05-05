package model

import "time"

type User struct {
	ID           int       `json:"id"`
	TelegramID   int64     `json:"telegram_id"`
	Phone        *string   `json:"phone"`
	FullName     *string   `json:"full_name"`
	Role         string    `json:"role"`
	IsIdentified bool      `json:"is_identified"`
	PINFL        *string   `json:"pinfl"`
	CreatedAt    time.Time `json:"created_at"`
}
