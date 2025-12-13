package domain

import (
	"time"
)

type Courier struct {
	ID            int64
	Name          string
	Phone         string
	Status        CourierStatus
	TransportType TransportType
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
