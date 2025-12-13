//go:build integration

package integration

import (
	"context"
	"testing"
	"time"

	"avito/internal/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type DeliveryTestSuite struct {
	suite.Suite
	deps *TestDependencies
	ctx  context.Context
}

func TestDeliveryTestSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("Skip integration tests")
	}
	suite.Run(t, new(DeliveryTestSuite))
}

func (s *DeliveryTestSuite) SetupSuite() {
	s.ctx = context.Background()
	deps, err := SetupTestDependencies(s.ctx, s.T())
	s.Require().NoError(err)
	s.deps = deps
}

func (s *DeliveryTestSuite) TearDownSuite() {
	if s.deps != nil {
		if s.deps.Pool != nil {
			s.deps.Pool.Close()
		}
		if s.deps.Container != nil {
			_ = s.deps.Container.Terminate(s.ctx)
		}
	}
}

func (s *DeliveryTestSuite) SetupTest() {
	s.Require().NoError(CleanTestData(s.ctx, s.deps.Pool))
}

func (s *DeliveryTestSuite) TestAssignDeliverySuccess() {
	c := &domain.Courier{
		Name:          "Vanya",
		Phone:         "+79456432345",
		Status:        domain.CourierStatusAvailable,
		TransportType: domain.TransportCar,
	}
	cID, err := s.deps.CourierUseCase.CreateCourier(s.ctx, c)
	s.Require().NoError(err)

	orderID := uuid.NewString()

	d, err := s.deps.DeliveryUseCase.AssignDelivery(s.ctx, orderID)
	s.Require().NoError(err)
	s.NotNil(d)
	s.Equal(cID, d.CourierId)
	s.Equal(orderID, d.OrderId)
	s.True(domain.TransportCar == d.TransportType, "Expected transport type %s, got %s", domain.TransportCar, d.TransportType)
	s.True(d.Deadline.After(time.Now()))

	got, err := s.deps.CourierUseCase.GetCourier(s.ctx, cID)
	s.Require().NoError(err)
	s.True(domain.CourierStatusBusy == got.Status, "Expected status %s, got %s", domain.CourierStatusBusy, got.Status)
}

func (s *DeliveryTestSuite) TestUnassignDeliverySuccess() {
	c := &domain.Courier{
		Name:          "Glent",
		Phone:         "+79183965839",
		Status:        domain.CourierStatusAvailable,
		TransportType: domain.TransportCar,
	}
	cID, err := s.deps.CourierUseCase.CreateCourier(s.ctx, c)
	s.Require().NoError(err)

	orderID := uuid.NewString()
	_, err = s.deps.DeliveryUseCase.AssignDelivery(s.ctx, orderID)
	s.Require().NoError(err)

	u, err := s.deps.DeliveryUseCase.UnassignDelivery(s.ctx, orderID)
	s.Require().NoError(err)
	s.Equal(cID, u.CourierId)
	s.Equal(orderID, u.OrderId)

	got, err := s.deps.CourierUseCase.GetCourier(s.ctx, cID)
	s.Require().NoError(err)
	s.True(domain.CourierStatusAvailable == got.Status, "Expected status %s, got %s", domain.CourierStatusAvailable, got.Status)
}

func (s *DeliveryTestSuite) TestAssignDeliveryDuplicateOrder() {
	c := &domain.Courier{
		Name:          "Ilya",
		Phone:         "+79273456754",
		Status:        domain.CourierStatusAvailable,
		TransportType: domain.TransportCar,
	}
	_, err := s.deps.CourierUseCase.CreateCourier(s.ctx, c)
	s.Require().NoError(err)

	orderID := uuid.NewString()
	_, err = s.deps.DeliveryUseCase.AssignDelivery(s.ctx, orderID)
	s.Require().NoError(err)

	_, err = s.deps.DeliveryUseCase.AssignDelivery(s.ctx, orderID)
	s.Error(err)
}

func (s *DeliveryTestSuite) TestAssignDeliveryNoAvailableCouriers() {
	orderID := uuid.NewString()
	d, err := s.deps.DeliveryUseCase.AssignDelivery(s.ctx, orderID)
	s.Error(err)
	s.Nil(d)
	s.Equal(domain.ErrNoAvailableCouriers, err)
}

func (s *DeliveryTestSuite) TestAssignDeliveryInvalidOrderID() {
	c := &domain.Courier{
		Name:          "Lev",
		Phone:         "+79123456782",
		Status:        domain.CourierStatusAvailable,
		TransportType: domain.TransportCar,
	}
	_, err := s.deps.CourierUseCase.CreateCourier(s.ctx, c)
	s.Require().NoError(err)

	invalid := "not-a-uuid"
	d, err := s.deps.DeliveryUseCase.AssignDelivery(s.ctx, invalid)
	s.Error(err)
	s.Nil(d)
	s.Equal(domain.ErrInvalidId, err)
}

func (s *DeliveryTestSuite) TestAssignDeliveryMultipleCouriers() {
	c1 := &domain.Courier{
		Name:          "Alex",
		Phone:         "+79123456783",
		Status:        domain.CourierStatusAvailable,
		TransportType: domain.TransportCar,
	}
	c2 := &domain.Courier{
		Name:          "Kirill",
		Phone:         "+79123456784",
		Status:        domain.CourierStatusAvailable,
		TransportType: domain.TransportOnFoot,
	}
	_, err := s.deps.CourierUseCase.CreateCourier(s.ctx, c1)
	s.Require().NoError(err)
	_, err = s.deps.CourierUseCase.CreateCourier(s.ctx, c2)
	s.Require().NoError(err)

	ids := []string{uuid.NewString(), uuid.NewString()}
	var deliveries []*domain.Delivery
	for _, id := range ids {
		d, err := s.deps.DeliveryUseCase.AssignDelivery(s.ctx, id)
		s.Require().NoError(err)
		deliveries = append(deliveries, d)
	}

	s.Len(deliveries, 2)
	s.NotEqual(deliveries[0].CourierId, deliveries[1].CourierId)
	for _, d := range deliveries {
		got, err := s.deps.CourierUseCase.GetCourier(s.ctx, d.CourierId)
		s.Require().NoError(err)
		s.True(domain.CourierStatusBusy == got.Status, "Expected status %s, got %s", domain.CourierStatusBusy, got.Status)
	}
}

func (s *DeliveryTestSuite) TestReleaseExpiredDeliveries() {
	c := &domain.Courier{
		Name:          "Artur",
		Phone:         "+79123412345",
		Status:        domain.CourierStatusAvailable,
		TransportType: domain.TransportCar,
	}
	cID, err := s.deps.CourierUseCase.CreateCourier(s.ctx, c)
	s.Require().NoError(err)

	order := uuid.NewString()
	expired := time.Now().UTC().Add(-2 * time.Hour)

	_, err = s.deps.DeliveryRepo.Assign(s.ctx, order, &domain.Courier{ID: cID, TransportType: domain.TransportCar}, expired)
	s.Require().NoError(err)

	got, err := s.deps.CourierUseCase.GetCourier(s.ctx, cID)
	s.Require().NoError(err)
	s.True(domain.CourierStatusBusy == got.Status, "Expected status %s, got %s", domain.CourierStatusBusy, got.Status)

	s.Require().NoError(s.deps.DeliveryRepo.ReleaseExpireDeliveries(s.ctx))

	got2, err := s.deps.CourierUseCase.GetCourier(s.ctx, cID)
	s.Require().NoError(err)
	s.True(domain.CourierStatusAvailable == got2.Status, "Expected status %s, got %s", domain.CourierStatusAvailable, got2.Status)
}

func (s *DeliveryTestSuite) TestFullDeliveryProcess() {
	c := &domain.Courier{
		Name:          "Natasha",
		Phone:         "+79183647586",
		Status:        domain.CourierStatusAvailable,
		TransportType: domain.TransportCar,
	}
	cID, err := s.deps.CourierUseCase.CreateCourier(s.ctx, c)
	s.Require().NoError(err)

	order := uuid.NewString()
	d, err := s.deps.DeliveryUseCase.AssignDelivery(s.ctx, order)
	s.Require().NoError(err)
	s.Equal(cID, d.CourierId)

	got, err := s.deps.CourierUseCase.GetCourier(s.ctx, cID)
	s.Require().NoError(err)
	s.True(domain.CourierStatusBusy == got.Status, "Expected status %s, got %s", domain.CourierStatusBusy, got.Status)

	_, err = s.deps.DeliveryUseCase.UnassignDelivery(s.ctx, order)
	s.Require().NoError(err)

	got2, err := s.deps.CourierUseCase.GetCourier(s.ctx, cID)
	s.Require().NoError(err)
	s.True(domain.CourierStatusAvailable == got2.Status, "Expected status %s, got %s", domain.CourierStatusAvailable, got2.Status)

	newOrder := uuid.NewString()
	d2, err := s.deps.DeliveryUseCase.AssignDelivery(s.ctx, newOrder)
	s.Require().NoError(err)
	s.Equal(cID, d2.CourierId)
}
