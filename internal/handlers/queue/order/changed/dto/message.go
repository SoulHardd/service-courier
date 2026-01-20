package dto

import "time"

type Message struct {
	OrderId   string    `json:"order_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}
