package delivery

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/SoulHardd/service-courier/internal/domain"
	"github.com/SoulHardd/service-courier/internal/handlers/http/delivery/dto"
	"github.com/SoulHardd/service-courier/mocks/usecase"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := usecase.NewMockdeliveryUseCase(ctrl)
	controller := New(mockUseCase)

	t.Run("Success", func(t *testing.T) {
		createRequest := dto.DeliveryCreateRequest{
			OrderId: "12345678-1234-1234-1234-123456789012",
		}

		expectedDelivery := &domain.Delivery{
			CourierId:     1,
			OrderId:       "12345678-1234-1234-1234-123456789012",
			TransportType: domain.TransportCar,
			Deadline:      time.Now().Add(30 * time.Minute),
		}

		mockUseCase.EXPECT().
			AssignDelivery(gomock.Any(), "12345678-1234-1234-1234-123456789012").
			Return(expectedDelivery, nil)

		body, _ := json.Marshal(createRequest)
		req := httptest.NewRequest("POST", "/delivery/assign", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		controller.Create(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response dto.DeliveryCreateResponse
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response))

		assert.Equal(t, expectedDelivery.OrderId, response.OrderId)
		assert.Equal(t, expectedDelivery.CourierId, response.CourierId)
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/delivery/assign", bytes.NewReader([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		controller.Create(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("NoAvailableCouriers", func(t *testing.T) {
		createRequest := dto.DeliveryCreateRequest{
			OrderId: "12345678-1234-1234-1234-123456789012",
		}

		mockUseCase.EXPECT().
			AssignDelivery(gomock.Any(), "12345678-1234-1234-1234-123456789012").
			Return(nil, domain.ErrNoAvailableCouriers)

		body, _ := json.Marshal(createRequest)
		req := httptest.NewRequest("POST", "/delivery/assign", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		controller.Create(rec, req)

		assert.Equal(t, http.StatusConflict, rec.Code)
	})

	t.Run("DeliveryAlreadyExists", func(t *testing.T) {
		createRequest := dto.DeliveryCreateRequest{
			OrderId: "12345678-1234-1234-1234-123456789012",
		}

		mockUseCase.EXPECT().
			AssignDelivery(gomock.Any(), "12345678-1234-1234-1234-123456789012").
			Return(nil, domain.ErrDeliveryAlreadyExists)

		body, _ := json.Marshal(createRequest)
		req := httptest.NewRequest("POST", "/delivery/assign", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		controller.Create(rec, req)

		assert.Equal(t, http.StatusConflict, rec.Code)
	})

	t.Run("MissingRequiredFields", func(t *testing.T) {
		createRequest := dto.DeliveryCreateRequest{
			OrderId: "",
		}

		mockUseCase.EXPECT().
			AssignDelivery(gomock.Any(), "").
			Return(nil, domain.ErrMissingRequiredFields)

		body, _ := json.Marshal(createRequest)
		req := httptest.NewRequest("POST", "/delivery/assign", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		controller.Create(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("InvalidId", func(t *testing.T) {
		createRequest := dto.DeliveryCreateRequest{
			OrderId: "invalid-order",
		}

		mockUseCase.EXPECT().
			AssignDelivery(gomock.Any(), "invalid-order").
			Return(nil, domain.ErrInvalidId)

		body, _ := json.Marshal(createRequest)
		req := httptest.NewRequest("POST", "/delivery/assign", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		controller.Create(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("InternalServerError", func(t *testing.T) {
		createRequest := dto.DeliveryCreateRequest{
			OrderId: "12345678-1234-1234-1234-123456789012",
		}

		mockUseCase.EXPECT().
			AssignDelivery(gomock.Any(), "12345678-1234-1234-1234-123456789012").
			Return(nil, errors.New("database error"))

		body, _ := json.Marshal(createRequest)
		req := httptest.NewRequest("POST", "/delivery/assign", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		controller.Create(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestDelete(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := usecase.NewMockdeliveryUseCase(ctrl)
	controller := New(mockUseCase)

	t.Run("Success", func(t *testing.T) {
		deleteRequest := dto.DeliveryDeleteRequest{
			OrderId: "12345678-1234-1234-1234-123456789012",
		}

		expectedDelivery := &domain.Delivery{
			CourierId: 1,
			OrderId:   "12345678-1234-1234-1234-123456789012",
		}

		mockUseCase.EXPECT().
			UnassignDelivery(gomock.Any(), "12345678-1234-1234-1234-123456789012").
			Return(expectedDelivery, nil)

		body, _ := json.Marshal(deleteRequest)
		req := httptest.NewRequest("POST", "/delivery/unassign", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		controller.Delete(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
		}

		assert.Equal(t, http.StatusOK, rec.Code)

		var response dto.DeliveryDeleteResponse
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response))

		assert.Equal(t, expectedDelivery.OrderId, response.OrderId)
		assert.Equal(t, expectedDelivery.CourierId, response.CourierId)
		assert.Equal(t, "unassigned", response.Status)
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/delivery/unassign", bytes.NewReader([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		controller.Delete(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("InvalidId", func(t *testing.T) {
		deleteRequest := dto.DeliveryDeleteRequest{
			OrderId: "invalid-order",
		}

		mockUseCase.EXPECT().
			UnassignDelivery(gomock.Any(), "invalid-order").
			Return(nil, domain.ErrInvalidId)

		body, _ := json.Marshal(deleteRequest)
		req := httptest.NewRequest("POST", "/delivery/unassign", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		controller.Delete(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("MissingRequiredFields", func(t *testing.T) {
		deleteRequest := dto.DeliveryDeleteRequest{
			OrderId: "",
		}

		mockUseCase.EXPECT().
			UnassignDelivery(gomock.Any(), "").
			Return(nil, domain.ErrMissingRequiredFields)

		body, _ := json.Marshal(deleteRequest)
		req := httptest.NewRequest("POST", "/delivery/unassign", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		controller.Delete(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("DeliveryNotFound", func(t *testing.T) {
		deleteRequest := dto.DeliveryDeleteRequest{
			OrderId: "non-existent-order",
		}

		mockUseCase.EXPECT().
			UnassignDelivery(gomock.Any(), "non-existent-order").
			Return(nil, domain.ErrDeliveryNotFound)

		body, _ := json.Marshal(deleteRequest)
		req := httptest.NewRequest("POST", "/delivery/unassign", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		controller.Delete(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("InternalServerError", func(t *testing.T) {
		deleteRequest := dto.DeliveryDeleteRequest{
			OrderId: "12345678-1234-1234-1234-123456789012",
		}

		mockUseCase.EXPECT().
			UnassignDelivery(gomock.Any(), "12345678-1234-1234-1234-123456789012").
			Return(nil, errors.New("database error"))

		body, _ := json.Marshal(deleteRequest)
		req := httptest.NewRequest("POST", "/delivery/unassign", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		controller.Delete(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}
