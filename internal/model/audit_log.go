package model

import "time"

type AuditLog struct {
	ID         int       `json:"id"`
	UserID     int       `json:"user_id"`
	Action     string    `json:"action"`
	EntityType string    `json:"entity_type"`
	EntityID   int       `json:"entity_id"`
	OldValues  *string   `json:"old_values"`
	NewValues  *string   `json:"new_values"`
	CreatedAt  time.Time `json:"created_at"`
}
