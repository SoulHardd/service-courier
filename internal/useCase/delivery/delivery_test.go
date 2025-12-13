package delivery

import (
	"avito/internal/domain"
	"avito/mocks/repository"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestAssignDelivery(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDeliveryRepo := repository.NewMockDeliveryRepository(ctrl)
	mockCourierRepo := repository.NewMockCourierRepository(ctrl)
	useCase := New(mockDeliveryRepo, mockCourierRepo)

	t.Run("Success", func(t *testing.T) {
		orderID := "1g6dk4o6-ffff-ffff-ffff-17r5fff27495"
		expectedCourier := &domain.Courier{
			ID:            1,
			Name:          "David",
			Phone:         "+79183756435",
			Status:        "available",
			TransportType: "car",
		}

		expectedDeadline := time.Now().Add(30 * time.Minute)
		expectedDelivery := &domain.Delivery{
			CourierId:     1,
			OrderId:       orderID,
			TransportType: "car",
			Deadline:      expectedDeadline,
		}

		mockCourierRepo.EXPECT().
			GetOneForDelivery(gomock.Any()).
			Return(expectedCourier, nil)

		mockDeliveryRepo.EXPECT().
			Assign(gomock.Any(), orderID, expectedCourier, gomock.Any()).
			DoAndReturn(func(ctx context.Context, orderId string, courier *domain.Courier, deadline time.Time) (*domain.Delivery, error) {
				return &domain.Delivery{
					CourierId:     courier.ID,
					OrderId:       orderId,
					TransportType: courier.TransportType,
					Deadline:      expectedDeadline,
				}, nil
			})

		delivery, err := useCase.AssignDelivery(context.Background(), orderID)

		assert.NoError(t, err)
		assert.Equal(t, expectedDelivery.OrderId, delivery.OrderId)
		assert.Equal(t, expectedDelivery.CourierId, delivery.CourierId)
	})

	t.Run("MissingOrderId", func(t *testing.T) {
		delivery, err := useCase.AssignDelivery(context.Background(), "")

		assert.ErrorIs(t, err, domain.ErrMissingRequiredFields)
		assert.Nil(t, delivery)
	})

	t.Run("NoAvailableCouriers", func(t *testing.T) {
		orderID := "1g6dk4o6-ffff-ffff-ffff-17r5fff27495"

		mockCourierRepo.EXPECT().
			GetOneForDelivery(gomock.Any()).
			Return(nil, domain.ErrNoAvailableCouriers)

		delivery, err := useCase.AssignDelivery(context.Background(), orderID)

		assert.ErrorIs(t, err, domain.ErrNoAvailableCouriers)
		assert.Nil(t, delivery)
	})

	t.Run("UnsupportedTransport", func(t *testing.T) {
		orderID := "1g6dk4o6-ffff-ffff-ffff-17r5fff27495"
		courier := &domain.Courier{
			ID:            1,
			Name:          "Leonid",
			Phone:         "+79283745865",
			Status:        "available",
			TransportType: "unknown",
		}

		mockCourierRepo.EXPECT().
			GetOneForDelivery(gomock.Any()).
			Return(courier, nil)

		delivery, err := useCase.AssignDelivery(context.Background(), orderID)

		assert.ErrorIs(t, err, domain.ErrUnsupportedTransport)
		assert.Nil(t, delivery)
	})

	t.Run("DeliveryAlreadyExists", func(t *testing.T) {
		orderID := "1g6dk4o6-ffff-ffff-ffff-17r5fff27495"
		courier := &domain.Courier{
			ID:            1,
			Name:          "George",
			Phone:         "+77777777777",
			Status:        "available",
			TransportType: "scooter",
		}

		mockCourierRepo.EXPECT().
			GetOneForDelivery(gomock.Any()).
			Return(courier, nil)

		mockDeliveryRepo.EXPECT().
			Assign(gomock.Any(), orderID, courier, gomock.Any()).
			Return(nil, domain.ErrDeliveryAlreadyExists)

		delivery, err := useCase.AssignDelivery(context.Background(), orderID)

		assert.ErrorIs(t, err, domain.ErrDeliveryAlreadyExists)
		assert.Nil(t, delivery)
	})

	t.Run("CourierNotFound", func(t *testing.T) {
		orderID := "1g6dk4o6-ffff-ffff-ffff-17r5fff27495"
		courier := &domain.Courier{
			ID:            1,
			Name:          "Anton",
			Phone:         "+79124758639",
			Status:        "available",
			TransportType: "car",
		}

		mockCourierRepo.EXPECT().
			GetOneForDelivery(gomock.Any()).
			Return(courier, nil)

		mockDeliveryRepo.EXPECT().
			Assign(gomock.Any(), orderID, courier, gomock.Any()).
			Return(nil, domain.ErrCourierNotFound)

		delivery, err := useCase.AssignDelivery(context.Background(), orderID)

		assert.ErrorIs(t, err, domain.ErrCourierNotFound)
		assert.Nil(t, delivery)
	})

	t.Run("InvalidId", func(t *testing.T) {
		orderID := "invalid-order"
		courier := &domain.Courier{
			ID:            1,
			Name:          "Lena",
			Phone:         "+79123111111",
			Status:        "available",
			TransportType: "on_foot",
		}

		mockCourierRepo.EXPECT().
			GetOneForDelivery(gomock.Any()).
			Return(courier, nil)

		mockDeliveryRepo.EXPECT().
			Assign(gomock.Any(), orderID, courier, gomock.Any()).
			Return(nil, domain.ErrInvalidId)

		delivery, err := useCase.AssignDelivery(context.Background(), orderID)

		assert.ErrorIs(t, err, domain.ErrInvalidId)
		assert.Nil(t, delivery)
	})

	t.Run("RepositoryError", func(t *testing.T) {
		orderID := "1g6dk4o6-ffff-ffff-ffff-17r5fff27495"
		courier := &domain.Courier{
			ID:            1,
			Name:          "Misha",
			Phone:         "+79364758697",
			Status:        "available",
			TransportType: "car",
		}

		expectedErr := errors.New("database error")

		mockCourierRepo.EXPECT().
			GetOneForDelivery(gomock.Any()).
			Return(courier, nil)

		mockDeliveryRepo.EXPECT().
			Assign(gomock.Any(), orderID, courier, gomock.Any()).
			Return(nil, expectedErr)

		delivery, err := useCase.AssignDelivery(context.Background(), orderID)

		assert.ErrorIs(t, err, expectedErr)
		assert.Nil(t, delivery)
	})
}

func TestUnassignDelivery(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDeliveryRepo := repository.NewMockDeliveryRepository(ctrl)
	mockCourierRepo := repository.NewMockCourierRepository(ctrl)
	useCase := New(mockDeliveryRepo, mockCourierRepo)

	t.Run("Success", func(t *testing.T) {
		orderID := "1g6dk4o6-ffff-ffff-ffff-17r5fff27495"
		expectedDelivery := &domain.Delivery{
			CourierId: 1,
			OrderId:   orderID,
		}

		mockDeliveryRepo.EXPECT().
			Unassign(gomock.Any(), orderID).
			Return(expectedDelivery, nil)

		delivery, err := useCase.UnassignDelivery(context.Background(), orderID)

		assert.NoError(t, err)
		assert.Equal(t, expectedDelivery.OrderId, delivery.OrderId)
		assert.Equal(t, expectedDelivery.CourierId, delivery.CourierId)
	})

	t.Run("MissingOrderId", func(t *testing.T) {
		delivery, err := useCase.UnassignDelivery(context.Background(), "")

		assert.ErrorIs(t, err, domain.ErrMissingRequiredFields)
		assert.Nil(t, delivery)
	})

	t.Run("InvalidId", func(t *testing.T) {
		orderID := "invalid-order"

		mockDeliveryRepo.EXPECT().
			Unassign(gomock.Any(), orderID).
			Return(nil, domain.ErrInvalidId)

		delivery, err := useCase.UnassignDelivery(context.Background(), orderID)

		assert.ErrorIs(t, err, domain.ErrInvalidId)
		assert.Nil(t, delivery)
	})

	t.Run("DeliveryNotFound", func(t *testing.T) {
		orderID := "non-existent-order"

		mockDeliveryRepo.EXPECT().
			Unassign(gomock.Any(), orderID).
			Return(nil, domain.ErrDeliveryNotFound)

		delivery, err := useCase.UnassignDelivery(context.Background(), orderID)

		assert.ErrorIs(t, err, domain.ErrDeliveryNotFound)
		assert.Nil(t, delivery)
	})

	t.Run("RepositoryError", func(t *testing.T) {
		orderID := "1g6dk4o6-ffff-ffff-ffff-17r5fff27495"
		expectedErr := errors.New("database error")

		mockDeliveryRepo.EXPECT().
			Unassign(gomock.Any(), orderID).
			Return(nil, expectedErr)

		delivery, err := useCase.UnassignDelivery(context.Background(), orderID)

		assert.ErrorIs(t, err, expectedErr)
		assert.Nil(t, delivery)
	})
}
