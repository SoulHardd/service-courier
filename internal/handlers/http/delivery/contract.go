//go:generate mockgen -source=contract.go -destination=../../../mocks/usecase/mock_delivery_usecase.go -package=usecase
package delivery

import (
	"context"

	"github.com/SoulHardd/service-courier/internal/domain"
)

type deliveryUseCase interface {
	AssignDelivery(ctx context.Context, orderId string) (*domain.Delivery, error)
	UnassignDelivery(ctx context.Context, orderId string) (*domain.Delivery, error)
}
