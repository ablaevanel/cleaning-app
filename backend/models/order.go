package models

import "time"

type Order struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	ServiceID int       `json:"service_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}
