package courier

import (
	"avito/internal/domain"
	"avito/mocks/repository"
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetCourier(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockCourierRepository(ctrl)
	useCase := New(mockRepo)

	t.Run("Success", func(t *testing.T) {
		expectedCourier := &domain.Courier{
			ID:            1,
			Name:          "Max",
			Phone:         "+79485768392",
			Status:        "available",
			TransportType: "car",
		}

		mockRepo.EXPECT().
			GetOneById(gomock.Any(), int64(1)).
			Return(expectedCourier, nil)

		courier, err := useCase.GetCourier(context.Background(), 1)

		assert.NoError(t, err)
		assert.Equal(t, expectedCourier.ID, courier.ID)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockRepo.EXPECT().
			GetOneById(gomock.Any(), int64(999)).
			Return(nil, domain.ErrCourierNotFound)

		courier, err := useCase.GetCourier(context.Background(), 999)

		assert.ErrorIs(t, err, domain.ErrCourierNotFound)
		assert.Nil(t, courier)
	})

	t.Run("RepositoryError", func(t *testing.T) {
		expectedErr := errors.New("database error")

		mockRepo.EXPECT().
			GetOneById(gomock.Any(), int64(2)).
			Return(nil, expectedErr)

		courier, err := useCase.GetCourier(context.Background(), 2)

		assert.ErrorIs(t, err, expectedErr)
		assert.Nil(t, courier)
	})
}

func TestGetCouriers(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockCourierRepository(ctrl)
	useCase := New(mockRepo)

	t.Run("Success", func(t *testing.T) {
		expectedCouriers := []domain.Courier{
			{
				ID:            1,
				Name:          "Vlad",
				Phone:         "+79891645323",
				Status:        "available",
				TransportType: "car",
			},
			{
				ID:            2,
				Name:          "Vova",
				Phone:         "+77772556483",
				Status:        "busy",
				TransportType: "on_foot",
			},
		}

		mockRepo.EXPECT().
			GetAll(gomock.Any()).
			Return(expectedCouriers, nil)

		couriers, err := useCase.GetCouriers(context.Background())

		assert.NoError(t, err)
		assert.Len(t, couriers, len(expectedCouriers))
	})

	t.Run("EmptyList", func(t *testing.T) {
		mockRepo.EXPECT().
			GetAll(gomock.Any()).
			Return([]domain.Courier{}, nil)

		couriers, err := useCase.GetCouriers(context.Background())

		assert.NoError(t, err)
		assert.Len(t, couriers, 0)
	})

	t.Run("RepositoryError", func(t *testing.T) {
		expectedErr := errors.New("database error")

		mockRepo.EXPECT().
			GetAll(gomock.Any()).
			Return(nil, expectedErr)

		couriers, err := useCase.GetCouriers(context.Background())

		assert.ErrorIs(t, err, expectedErr)
		assert.Nil(t, couriers)
	})
}

func TestCreateCourier(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockCourierRepository(ctrl)
	useCase := New(mockRepo)

	t.Run("Success", func(t *testing.T) {
		courier := &domain.Courier{
			Name:          "Katerina",
			Phone:         "+71235783215",
			Status:        "available",
			TransportType: "car",
		}

		expectedID := int64(123)

		mockRepo.EXPECT().
			Create(gomock.Any(), courier).
			Return(expectedID, nil)

		id, err := useCase.CreateCourier(context.Background(), courier)

		assert.NoError(t, err)
		assert.Equal(t, expectedID, id)
	})

	t.Run("MissingRequiredFields", func(t *testing.T) {
		testCases := []struct {
			name    string
			courier *domain.Courier
			field   string
		}{
			{"MissingName", &domain.Courier{Phone: "+71284657384", Status: "available"}, "Name"},
			{"MissingPhone", &domain.Courier{Name: "Soul", Status: "available"}, "Phone"},
			{"MissingStatus", &domain.Courier{Name: "Walter", Phone: "+79893544625"}, "Status"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				id, err := useCase.CreateCourier(context.Background(), tc.courier)

				assert.ErrorIs(t, err, domain.ErrMissingRequiredFields)
				assert.Zero(t, id)
			})
		}
	})

	t.Run("InvalidPhone", func(t *testing.T) {
		invalidPhones := []string{
			"89123456789",
			"+7912345678",
			"+791234567890",
			"+7abc456789",
			"+7912 345 67 89",
		}

		for _, phone := range invalidPhones {
			t.Run(phone, func(t *testing.T) {
				courier := &domain.Courier{
					Name:   "Hank",
					Phone:  phone,
					Status: "available",
				}

				id, err := useCase.CreateCourier(context.Background(), courier)

				assert.ErrorIs(t, err, domain.ErrInvalidPhone)
				assert.Zero(t, id)
			})
		}
	})

	t.Run("PhoneExists", func(t *testing.T) {
		courier := &domain.Courier{
			Name:          "Jessie",
			Phone:         "+73647566325",
			Status:        "available",
			TransportType: "car",
		}

		mockRepo.EXPECT().
			Create(gomock.Any(), courier).
			Return(int64(0), domain.ErrPhoneExists)

		id, err := useCase.CreateCourier(context.Background(), courier)

		assert.ErrorIs(t, err, domain.ErrPhoneExists)
		assert.Zero(t, id)
	})

	t.Run("RepositoryError", func(t *testing.T) {
		courier := &domain.Courier{
			Name:          "Artem",
			Phone:         "+79894567832",
			Status:        "available",
			TransportType: "car",
		}

		expectedErr := errors.New("database error")

		mockRepo.EXPECT().
			Create(gomock.Any(), courier).
			Return(int64(0), expectedErr)

		id, err := useCase.CreateCourier(context.Background(), courier)

		assert.ErrorIs(t, err, expectedErr)
		assert.Zero(t, id)
	})
}

func TestUpdateCourier(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockCourierRepository(ctrl)
	useCase := New(mockRepo)

	t.Run("Success", func(t *testing.T) {
		courier := &domain.Courier{
			ID:            1,
			Name:          "Natasha",
			Phone:         "+79142645728",
			Status:        "busy",
			TransportType: "scooter",
		}

		mockRepo.EXPECT().
			Update(gomock.Any(), courier).
			Return(nil)

		err := useCase.UpdateCourier(context.Background(), courier)

		assert.NoError(t, err)
	})

	t.Run("MissingRequiredFields", func(t *testing.T) {
		testCases := []struct {
			name    string
			courier *domain.Courier
			field   string
		}{
			{"MissingName", &domain.Courier{ID: 1, Phone: "+79898988989", Status: "available"}, "Name"},
			{"MissingPhone", &domain.Courier{ID: 1, Name: "Vlad", Status: "available"}, "Phone"},
			{"MissingStatus", &domain.Courier{ID: 1, Name: "Ravshan", Phone: "+79485768345"}, "Status"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				err := useCase.UpdateCourier(context.Background(), tc.courier)

				assert.ErrorIs(t, err, domain.ErrMissingRequiredFields)
			})
		}
	})

	t.Run("InvalidPhone", func(t *testing.T) {
		courier := &domain.Courier{
			ID:     1,
			Name:   "Peter",
			Phone:  "phone",
			Status: "available",
		}

		err := useCase.UpdateCourier(context.Background(), courier)

		assert.ErrorIs(t, err, domain.ErrInvalidPhone)
	})

	t.Run("CourierNotFound", func(t *testing.T) {
		courier := &domain.Courier{
			ID:            999,
			Name:          "Arsen",
			Phone:         "+78293456764",
			Status:        "available",
			TransportType: "car",
		}

		mockRepo.EXPECT().
			Update(gomock.Any(), courier).
			Return(domain.ErrCourierNotFound)

		err := useCase.UpdateCourier(context.Background(), courier)

		assert.ErrorIs(t, err, domain.ErrCourierNotFound)
	})

	t.Run("PhoneExists", func(t *testing.T) {
		courier := &domain.Courier{
			ID:            1,
			Name:          "David",
			Phone:         "+79892647581",
			Status:        "available",
			TransportType: "car",
		}

		mockRepo.EXPECT().
			Update(gomock.Any(), courier).
			Return(domain.ErrPhoneExists)

		err := useCase.UpdateCourier(context.Background(), courier)

		assert.ErrorIs(t, err, domain.ErrPhoneExists)
	})

	t.Run("RepositoryError", func(t *testing.T) {
		courier := &domain.Courier{
			ID:            1,
			Name:          "Nikita",
			Phone:         "+79183957639",
			Status:        "available",
			TransportType: "on_foot",
		}

		expectedErr := errors.New("database error")

		mockRepo.EXPECT().
			Update(gomock.Any(), courier).
			Return(expectedErr)

		err := useCase.UpdateCourier(context.Background(), courier)

		assert.ErrorIs(t, err, expectedErr)
	})
}

func TestValidatePhoneNumber(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		phone    string
		expected bool
	}{
		{"+79123456789", true},
		{"+791234567890", false},
		{"+7912345678", false},
		{"89123456789", false},
		{"+7abc456789", false},
		{"+7912 3456789", false},
		{"", false},
	}

	for _, tc := range testCases {
		t.Run(tc.phone, func(t *testing.T) {
			result := validatePhoneNumber(tc.phone)
			assert.Equal(t, tc.expected, result)
		})
	}
}
