//go:generate mockgen -source=contract.go -destination=../../../mocks/repository/mock_courier_repo.go -package=repository
package courier

import (
	"avito/internal/domain"
	"context"
)

type CourierRepository interface {
	GetOneById(ctx context.Context, id int64) (*domain.Courier, error)
	GetAll(ctx context.Context) ([]domain.Courier, error)
	Create(ctx context.Context, courier *domain.Courier) (int64, error)
	Update(ctx context.Context, courier *domain.Courier) error
	GetOneForDelivery(ctx context.Context) (*domain.Courier, error)
	ReleaseOneByOrderId(ctx context.Context, orderId string) error
}
