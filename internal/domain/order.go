package domain

import "time"

type Order struct {
	ID        string
	Status    OrderStatus
	CreatedAt time.Time
}
