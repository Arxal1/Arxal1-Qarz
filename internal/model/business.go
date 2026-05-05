package model

import "time"

type Business struct {
	ID        int       `json:"id"`
	OwnerID   int       `json:"owner_id"`
	Name      string    `json:"name"`
	Phone     *string   `json:"phone"`
	CreatedAt time.Time `json:"created_at"`
}
