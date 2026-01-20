package changed

import "context"

type CreatedOrder struct {
	deliveryUseCase DeliveryUseCase
}

func NewCreatedOrder(uc DeliveryUseCase) *CreatedOrder {
	return &CreatedOrder{deliveryUseCase: uc}
}

func (o *CreatedOrder) Handle(ctx context.Context, orderId string) error {
	_, err := o.deliveryUseCase.AssignDelivery(ctx, orderId)
	return err
}
