//go:generate mockgen -source=contract.go -destination=../../../mocks/repository/mock_delivery_repo.go -package=repository
package delivery

import (
	"avito/internal/domain"
	"context"
	"time"
)

type DeliveryRepository interface {
	Assign(ctx context.Context, orderId string, c *domain.Courier, deadline time.Time) (*domain.Delivery, error)
	Unassign(ctx context.Context, orderId string) (*domain.Delivery, error)
	ReleaseExpireDeliveries(ctx context.Context) error
}
