package changed

import "context"

type CompletedOrder struct {
	courierUseCase CourierUseCase
}

func NewCompletedOrder(uc CourierUseCase) *CompletedOrder {
	return &CompletedOrder{courierUseCase: uc}
}

func (o *CompletedOrder) Handle(ctx context.Context, orderId string) error {
	err := o.courierUseCase.ReleaseCourierByOrderId(ctx, orderId)
	if err != nil {
		return err
	}
	return nil
}
