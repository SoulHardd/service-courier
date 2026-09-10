package delivery

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/SoulHardd/service-courier/internal/domain"
	"github.com/SoulHardd/service-courier/internal/factory"
	"github.com/SoulHardd/service-courier/internal/useCase/courier"
)

type DeliveryUseCase struct {
	deliveryRepo DeliveryRepository
	courierRepo  courier.CourierRepository
}

func New(d DeliveryRepository, c courier.CourierRepository) *DeliveryUseCase {
	return &DeliveryUseCase{
		deliveryRepo: d,
		courierRepo:  c,
	}
}

func (u *DeliveryUseCase) AssignDelivery(ctx context.Context, orderId string) (*domain.Delivery, error) {
	if orderId == "" {
		return nil, domain.ErrMissingRequiredFields
	}

	c, err := u.courierRepo.GetOneForDelivery(ctx)
	if err != nil {
		if errors.Is(err, domain.ErrNoAvailableCouriers) {
			return nil, domain.ErrNoAvailableCouriers
		}
		return nil, err
	}

	timeFactory := factory.DeliveryTimeFactory{}
	deliveryTime, err := timeFactory.Create(c.TransportType)
	if err != nil {
		if errors.Is(err, domain.ErrUnsupportedTransport) {
			return nil, domain.ErrUnsupportedTransport
		}
		return nil, err
	}
	deadline := deliveryTime.Deadline()

	d, err := u.deliveryRepo.Assign(ctx, orderId, c, deadline)

	if err != nil {
		switch {
		case errors.Is(err, domain.ErrCourierNotFound):
			return nil, domain.ErrCourierNotFound
		case errors.Is(err, domain.ErrInvalidId):
			return nil, domain.ErrInvalidId
		case errors.Is(err, domain.ErrDeliveryAlreadyExists):
			return nil, domain.ErrDeliveryAlreadyExists
		}
		return nil, err
	}

	return d, nil
}

func (u *DeliveryUseCase) UnassignDelivery(ctx context.Context, orderId string) (*domain.Delivery, error) {
	if orderId == "" {
		return nil, domain.ErrMissingRequiredFields
	}

	d, err := u.deliveryRepo.Unassign(ctx, orderId)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidId):
			return nil, domain.ErrInvalidId
		case errors.Is(err, domain.ErrDeliveryNotFound):
			return nil, domain.ErrDeliveryNotFound
		}
		return nil, err
	}
	return d, nil
}

func (u *DeliveryUseCase) StartMonitorDeliveries(ctx context.Context, interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if err := u.deliveryRepo.ReleaseExpireDeliveries(ctx); err != nil {
					fmt.Printf("Error releasing expired deliveries: %v", err)
				}
			}
		}
	}()
}
