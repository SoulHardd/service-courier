package worker

import (
	"avito/internal/domain"
	"context"
	"time"
)

type OrderGateway interface {
	GetOrders(ctx context.Context, cursor time.Time) ([]domain.Order, error)
	GetOrderById(ctx context.Context, id string) (domain.Order, error)
}

type DeliveryUseCase interface {
	AssignDelivery(ctx context.Context, orderId string) (*domain.Delivery, error)
}
