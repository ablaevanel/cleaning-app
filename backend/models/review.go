package models

import "time"

type Review struct {
    ID        int       `json:"id"`
    OrderID   int       `json:"order_id"`
    UserID    int       `json:"user_id"`
    Rating    int       `json:"rating"`
    Comment   string    `json:"comment"`
    CreatedAt time.Time `json:"created_at"`
}