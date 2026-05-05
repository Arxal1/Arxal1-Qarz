package model

import "time"

type Client struct {
	ID         int       `json:"id"`
	BusinessID int       `json:"business_id"`
	UserID     *int      `json:"user_id"`
	Name       string    `json:"name"`
	Phone      string    `json:"phone"`
	DebtLimit  int64     `json:"debt_limit"`
	CreatedAt  time.Time `json:"created_at"`
}
