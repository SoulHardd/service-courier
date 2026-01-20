package worker

import (
	"context"
	"log"
	"time"
)

type OrderWorker struct {
	orderGateway    OrderGateway
	deliveryUseCase DeliveryUseCase
	cursor          time.Time
	interval        time.Duration
}

func New(o OrderGateway, d DeliveryUseCase, interval time.Duration) *OrderWorker {
	return &OrderWorker{
		orderGateway:    o,
		deliveryUseCase: d,
		cursor:          time.Now().Add(-interval),
		interval:        interval,
	}
}

func (o *OrderWorker) Run(ctx context.Context) {
	go func() {
		log.Println("starting order worker...")
		ticker := time.NewTicker(o.interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				orders, err := o.orderGateway.GetOrders(ctx, o.cursor)
				if err != nil {
					log.Printf("Error getting orders: %v", err)
					continue
				}
				cursor := o.cursor
				for _, order := range orders {
					if _, err := o.deliveryUseCase.AssignDelivery(ctx, order.ID); err != nil {
						log.Printf("Error assign delivery: %v", err)
						continue
					}

					if order.CreatedAt.After(cursor) {
						cursor = order.CreatedAt
					}
				}

				o.cursor = cursor
			}
		}
	}()
}
