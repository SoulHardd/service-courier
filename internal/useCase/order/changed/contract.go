package changed

import (
	"avito/internal/domain"
	"context"
)

type OrderGateway interface {
	GetOrderById(ctx context.Context, id string) (domain.Order, error)
}

type DeliveryUseCase interface {
	AssignDelivery(ctx context.Context, orderId string) (*domain.Delivery, error)
	UnassignDelivery(ctx context.Context, orderId string) (*domain.Delivery, error)
}

type CourierUseCase interface {
	ReleaseCourierByOrderId(ctx context.Context, orderId string) error
}
