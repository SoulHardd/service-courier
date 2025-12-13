//go:generate mockgen -source=contract.go -destination=../../../mocks/usecase/mock_courier_usecase.go -package=usecase
package courier

import (
	"avito/internal/domain"
	"context"
)

type courierUseCase interface {
	GetCourier(ctx context.Context, id int64) (*domain.Courier, error)
	GetCouriers(ctx context.Context) ([]domain.Courier, error)
	CreateCourier(ctx context.Context, req *domain.Courier) (int64, error)
	UpdateCourier(ctx context.Context, req *domain.Courier) error
}
