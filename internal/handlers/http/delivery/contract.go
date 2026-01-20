//go:generate mockgen -source=contract.go -destination=../../../mocks/usecase/mock_delivery_usecase.go -package=usecase
package delivery

import (
	"avito/internal/domain"
	"context"
)

type deliveryUseCase interface {
	AssignDelivery(ctx context.Context, orderId string) (*domain.Delivery, error)
	UnassignDelivery(ctx context.Context, orderId string) (*domain.Delivery, error)
}
