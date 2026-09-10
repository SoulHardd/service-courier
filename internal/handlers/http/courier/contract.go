//go:generate mockgen -source=contract.go -destination=../../../mocks/usecase/mock_courier_usecase.go -package=usecase
package courier

import (
	"context"

	"github.com/SoulHardd/service-courier/internal/domain"
)

type courierUseCase interface {
	GetCourier(ctx context.Context, id int64) (*domain.Courier, error)
	GetCouriers(ctx context.Context) ([]domain.Courier, error)
	CreateCourier(ctx context.Context, req *domain.Courier) (int64, error)
	UpdateCourier(ctx context.Context, req *domain.Courier) error
}
