package changed

import (
	"avito/internal/domain"
	"context"
	"fmt"
)

type Handler interface {
	Handle(ctx context.Context, orderId string) error
}

type Factory struct {
	created   Handler
	cancelled Handler
	completed Handler
}

func NewFactory(deliveryUC DeliveryUseCase, courierUC CourierUseCase) *Factory {
	createdHandler := NewCreatedOrder(deliveryUC)
	cancelledHandler := NewCancelledOrder(deliveryUC)
	completedHandler := NewCompletedOrder(courierUC)
	return &Factory{
		created:   createdHandler,
		cancelled: cancelledHandler,
		completed: completedHandler,
	}
}

func (f *Factory) GetHandler(status domain.OrderStatus) (Handler, error) {
	switch status {
	case domain.OrderStatusCreated:
		return f.created, nil
	case domain.OrderStatusCancelled:
		return f.cancelled, nil
	case domain.OrderStatusCompleted:
		return f.completed, nil
	default:
		return nil, fmt.Errorf("unknown order status: %s", status)
	}
}
