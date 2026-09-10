//go:build integration

package integration

import (
	"context"
	"testing"

	"github.com/SoulHardd/service-courier/internal/domain"

	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
)

type CourierTestSuite struct {
	suite.Suite
	deps      *TestDependencies
	container testcontainers.Container
	ctx       context.Context
}

func TestCourierTestSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("Skip integration tests")
	}
	suite.Run(t, new(CourierTestSuite))
}

func (s *CourierTestSuite) SetupSuite() {
	s.ctx = context.Background()
	deps, err := SetupTestDependencies(s.ctx, s.T())
	s.Require().NoError(err)
	s.deps = deps
	s.container = deps.Container
}

func (s *CourierTestSuite) TearDownSuite() {
	if s.deps != nil {
		if s.deps.Pool != nil {
			s.deps.Pool.Close()
		}
		if s.container != nil {
			_ = s.container.Terminate(s.ctx)
		}
	}
}

func (s *CourierTestSuite) SetupTest() {
	s.Require().NoError(CleanTestData(s.ctx, s.deps.Pool))
}

func (s *CourierTestSuite) TestCreateAndGetCourier() {
	c := &domain.Courier{
		Name:          "Alex",
		Phone:         "+79893848586",
		Status:        domain.CourierStatusAvailable,
		TransportType: domain.TransportCar,
	}

	id, err := s.deps.CourierUseCase.CreateCourier(s.ctx, c)
	s.Require().NoError(err)
	s.Require().Greater(id, int64(0))

	got, err := s.deps.CourierUseCase.GetCourier(s.ctx, id)
	s.Require().NoError(err)
	s.Equal(c.Name, got.Name)
	s.Equal(c.Phone, got.Phone)
	s.True(c.Status == got.Status, "Expected status %s, got %s", c.Status, got.Status)
	s.True(c.TransportType == got.TransportType, "Expected transport type %s, got %s", c.TransportType, got.TransportType)
}

func (s *CourierTestSuite) TestGetAllCouriers() {
	c1 := &domain.Courier{
		Name:          "Kirill",
		Phone:         "+79112345781",
		Status:        domain.CourierStatusAvailable,
		TransportType: domain.TransportCar,
	}
	c2 := &domain.Courier{
		Name:          "Artur",
		Phone:         "+71234567890",
		Status:        domain.CourierStatusBusy,
		TransportType: domain.TransportOnFoot,
	}

	_, err := s.deps.CourierUseCase.CreateCourier(s.ctx, c1)
	s.Require().NoError(err)
	_, err = s.deps.CourierUseCase.CreateCourier(s.ctx, c2)
	s.Require().NoError(err)

	all, err := s.deps.CourierUseCase.GetCouriers(s.ctx)
	s.Require().NoError(err)
	s.Len(all, 2)
}

func (s *CourierTestSuite) TestUpdateCourier() {
	c := &domain.Courier{
		Name:          "Vlad",
		Phone:         "+77777777777",
		Status:        domain.CourierStatusAvailable,
		TransportType: domain.TransportCar,
	}

	id, err := s.deps.CourierUseCase.CreateCourier(s.ctx, c)
	s.Require().NoError(err)

	updated := &domain.Courier{
		ID:            id,
		Name:          "Anton",
		Phone:         "+77777777779",
		Status:        domain.CourierStatusBusy,
		TransportType: domain.TransportScooter,
	}

	s.Require().NoError(s.deps.CourierUseCase.UpdateCourier(s.ctx, updated))

	got, err := s.deps.CourierUseCase.GetCourier(s.ctx, id)
	s.Require().NoError(err)
	s.Equal(updated.Name, got.Name)
	s.Equal(updated.Phone, got.Phone)
	s.True(updated.Status == got.Status, "Expected status %s, got %s", updated.Status, got.Status)
	s.True(updated.TransportType == got.TransportType, "Expected transport type %s, got %s", updated.TransportType, got.TransportType)
}

func (s *CourierTestSuite) TestCreateCourierDuplicatePhone() {
	c1 := &domain.Courier{
		Name:          "Sasha",
		Phone:         "+79999999999",
		Status:        domain.CourierStatusAvailable,
		TransportType: domain.TransportCar,
	}
	_, err := s.deps.CourierUseCase.CreateCourier(s.ctx, c1)
	s.Require().NoError(err)

	c2 := &domain.Courier{
		Name:          "Masha",
		Phone:         "+79999999999",
		Status:        domain.CourierStatusAvailable,
		TransportType: domain.TransportCar,
	}
	_, err = s.deps.CourierUseCase.CreateCourier(s.ctx, c2)
	s.Error(err)
	s.Equal(domain.ErrPhoneExists, err)
}

func (s *CourierTestSuite) TestGetCourierNotFound() {
	_, err := s.deps.CourierUseCase.GetCourier(s.ctx, 9999)
	s.Error(err)
	s.Equal(domain.ErrCourierNotFound, err)
}
