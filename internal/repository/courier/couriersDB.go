package courier

import (
	"avito/internal/domain"
	"time"
)

type CourierDB struct {
	ID            int64
	Name          string
	Phone         string
	Status        domain.CourierStatus
	TransportType domain.TransportType
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
