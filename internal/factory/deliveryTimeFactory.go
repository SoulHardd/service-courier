package factory

import (
	"time"

	"github.com/SoulHardd/service-courier/internal/domain"
)

type DeliveryDeadline interface {
	Deadline() time.Time
}

type OnFoot struct{}

func (o OnFoot) Deadline() time.Time {
	return time.Now().Add(time.Minute * 30).UTC()
}

type Scooter struct{}

func (s Scooter) Deadline() time.Time {
	return time.Now().Add(time.Minute * 15).UTC()
}

type Car struct{}

func (c Car) Deadline() time.Time {
	return time.Now().Add(time.Minute * 5).UTC()
}

type DeliveryTimeFactory struct{}

func (d DeliveryTimeFactory) Create(transportType domain.TransportType) (DeliveryDeadline, error) {
	switch transportType {
	case domain.TransportOnFoot:
		return OnFoot{}, nil
	case domain.TransportScooter:
		return Scooter{}, nil
	case domain.TransportCar:
		return Car{}, nil
	default:
		return nil, domain.ErrUnsupportedTransport
	}
}
