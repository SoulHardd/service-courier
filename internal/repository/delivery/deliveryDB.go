package delivery

import "time"

type DeliveryDB struct {
	ID         int64
	CourierId  int64
	OrderId    string
	AssignedAt time.Time
	Deadline   time.Time
}
