package changed

import "context"

type CancelledOrder struct {
	deliveryUseCase DeliveryUseCase
}

func NewCancelledOrder(uc DeliveryUseCase) *CancelledOrder {
	return &CancelledOrder{deliveryUseCase: uc}
}

func (o *CancelledOrder) Handle(ctx context.Context, orderId string) error {
	_, err := o.deliveryUseCase.UnassignDelivery(ctx, orderId)
	return err
}
