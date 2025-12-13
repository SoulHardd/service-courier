package domain

import (
	"time"
)

type Delivery struct {
	CourierId     int64
	OrderId       string
	TransportType TransportType
	AssignedAt    time.Time
	Deadline      time.Time
}
