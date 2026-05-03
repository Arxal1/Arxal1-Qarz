package model

import "time"

type User struct {
	ID         int   `json:"id"`
	TelegramID int64 `json:"telegram_id"`

	// Уровень 1
	Phone        *string `json:"phone"`
	TelegramName string  `json:"telegram_name"`

	// Уровень 2
	PINFL        *string `json:"pinfl"`
	IsIdentified bool    `json:"is_identified"`
	RealName     *string `json:"real_name"`

	CreatedAt time.Time `json:"created_at"`
}
