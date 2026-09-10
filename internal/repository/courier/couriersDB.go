package courier

import (
	"time"

	"github.com/SoulHardd/service-courier/internal/domain"
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
