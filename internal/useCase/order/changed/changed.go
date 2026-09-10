package changed

import (
	"context"
	"fmt"
	"log"

	model "github.com/SoulHardd/service-courier/internal/domain"
)

type OrderUseCase struct {
	orderGateway    OrderGateway
	deliveryUseCase DeliveryUseCase
	courierUseCase  CourierUseCase
}

func NewOrderUseCase(g OrderGateway, d DeliveryUseCase, c CourierUseCase) *OrderUseCase {
	return &OrderUseCase{
		orderGateway:    g,
		deliveryUseCase: d,
		courierUseCase:  c,
	}
}

func (uc *OrderUseCase) Process(ctx context.Context, o model.Order) error {

	checkOrder, err := uc.orderGateway.GetOrderById(ctx, o.ID)
	if err != nil {
		log.Printf("failed to check order: %v", err)
	}

	factory := NewFactory(uc.deliveryUseCase, uc.courierUseCase)

	if checkOrder.Status != o.Status {
		log.Printf("order statuses mismatch: got %s, actual %s. Using actual status",
			o.Status, checkOrder.Status)

		handler, err := factory.GetHandler(checkOrder.Status)
		if err != nil {
			return fmt.Errorf("failed to get handler for status %s: %w", checkOrder.Status, err)
		}
		return handler.Handle(ctx, o.ID)
	}

	handler, err := factory.GetHandler(o.Status)
	if err != nil {
		return fmt.Errorf("failed to get handler: %w", err)
	}
	return handler.Handle(ctx, o.ID)
}
