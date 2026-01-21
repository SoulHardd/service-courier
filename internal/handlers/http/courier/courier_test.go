package courier

import (
	"avito/internal/domain"
	"avito/internal/handlers/http/courier/dto"
	"avito/mocks/usecase"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	chi "github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := usecase.NewMockcourierUseCase(ctrl)
	controller := New(mockUseCase)

	t.Run("Success", func(t *testing.T) {

		expectedCourier := &domain.Courier{
			ID:            1,
			Name:          "Vladislav",
			Phone:         "+77778899098",
			Status:        "available",
			TransportType: "car",
		}

		mockUseCase.EXPECT().
			GetCourier(gomock.Any(), int64(1)).
			Return(expectedCourier, nil)

		req := httptest.NewRequest("GET", "/courier/1", nil)
		rec := httptest.NewRecorder()

		r := chi.NewRouter()
		r.Get("/courier/{id}", controller.Get)
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response dto.CourierResponse
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response))
		assert.Equal(t, expectedCourier.ID, response.ID)

	})

	t.Run("InvalidID", func(t *testing.T) {

		req := httptest.NewRequest("GET", "/courier/invalid", nil)
		rec := httptest.NewRecorder()

		r := chi.NewRouter()
		r.Get("/courier/{id}", controller.Get)
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("NotFound", func(t *testing.T) {

		mockUseCase.EXPECT().
			GetCourier(gomock.Any(), int64(999)).
			Return(nil, domain.ErrCourierNotFound)

		req := httptest.NewRequest("GET", "/courier/999", nil)
		rec := httptest.NewRecorder()

		r := chi.NewRouter()
		r.Get("/courier/{id}", controller.Get)
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("InternalServerError", func(t *testing.T) {

		mockUseCase.EXPECT().
			GetCourier(gomock.Any(), int64(2)).
			Return(nil, errors.New("database error"))

		req := httptest.NewRequest("GET", "/courier/2", nil)
		rec := httptest.NewRecorder()

		r := chi.NewRouter()
		r.Get("/courier/{id}", controller.Get)
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestGetAll(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := usecase.NewMockcourierUseCase(ctrl)
	controller := New(mockUseCase)

	t.Run("Success", func(t *testing.T) {

		expectedCouriers := []domain.Courier{
			{
				ID:            1,
				Name:          "Max",
				Phone:         "+79123456789",
				Status:        "available",
				TransportType: "car",
			},
			{
				ID:            2,
				Name:          "Anton",
				Phone:         "+77384922034",
				Status:        "busy",
				TransportType: "on_foot",
			},
		}

		mockUseCase.EXPECT().
			GetCouriers(gomock.Any()).
			Return(expectedCouriers, nil)

		req := httptest.NewRequest("GET", "/couriers", nil)
		rec := httptest.NewRecorder()

		controller.GetAll(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response []dto.CourierResponses
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response))
		assert.Equal(t, len(expectedCouriers), len(response))
	})

	t.Run("EmptyList", func(t *testing.T) {

		mockUseCase.EXPECT().
			GetCouriers(gomock.Any()).
			Return([]domain.Courier{}, nil)

		req := httptest.NewRequest("GET", "/couriers", nil)
		rec := httptest.NewRecorder()

		controller.GetAll(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response []dto.CourierResponses
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response))
		assert.Len(t, response, 0)
	})

	t.Run("DatabaseError", func(t *testing.T) {

		mockUseCase.EXPECT().
			GetCouriers(gomock.Any()).
			Return(nil, errors.New("database error"))

		req := httptest.NewRequest("GET", "/couriers", nil)
		rec := httptest.NewRecorder()

		controller.GetAll(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestCreate(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := usecase.NewMockcourierUseCase(ctrl)
	controller := New(mockUseCase)

	t.Run("Success", func(t *testing.T) {

		createRequest := dto.CourierCreateRequest{
			Name:          "Sasha",
			Phone:         "+79892345234",
			Status:        "available",
			TransportType: "car",
		}

		expectedID := int64(123)

		mockUseCase.EXPECT().
			CreateCourier(gomock.Any(), gomock.Any()).
			Return(expectedID, nil)

		body, _ := json.Marshal(createRequest)
		req := httptest.NewRequest("POST", "/courier", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		controller.Create(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response map[string]interface{}
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response))
		assert.Equal(t, float64(expectedID), response["id"])
	})

	t.Run("InvalidJSON", func(t *testing.T) {

		req := httptest.NewRequest("POST", "/courier", bytes.NewReader([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		controller.Create(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("MissingRequiredFields", func(t *testing.T) {

		createRequest := dto.CourierCreateRequest{
			Name:  "",
			Phone: "+79283745869",
		}

		mockUseCase.EXPECT().
			CreateCourier(gomock.Any(), gomock.Any()).
			Return(int64(0), domain.ErrMissingRequiredFields)

		body, _ := json.Marshal(createRequest)
		req := httptest.NewRequest("POST", "/courier", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		controller.Create(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("InvalidPhone", func(t *testing.T) {

		createRequest := dto.CourierCreateRequest{
			Name:  "Vladimir",
			Phone: "not a number",
		}

		mockUseCase.EXPECT().
			CreateCourier(gomock.Any(), gomock.Any()).
			Return(int64(0), domain.ErrInvalidPhone)

		body, _ := json.Marshal(createRequest)
		req := httptest.NewRequest("POST", "/courier", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		controller.Create(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("PhoneExists", func(t *testing.T) {

		createRequest := dto.CourierCreateRequest{
			Name:  "Kirill",
			Phone: "+79892734600",
		}

		mockUseCase.EXPECT().
			CreateCourier(gomock.Any(), gomock.Any()).
			Return(int64(0), domain.ErrPhoneExists)

		body, _ := json.Marshal(createRequest)
		req := httptest.NewRequest("POST", "/courier", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		controller.Create(rec, req)

		assert.Equal(t, http.StatusConflict, rec.Code)
	})

	t.Run("InternalServerError", func(t *testing.T) {

		createRequest := dto.CourierCreateRequest{
			Name:  "Vlad",
			Phone: "+70003855693",
		}

		mockUseCase.EXPECT().
			CreateCourier(gomock.Any(), gomock.Any()).
			Return(int64(0), errors.New("unexpected error"))

		body, _ := json.Marshal(createRequest)
		req := httptest.NewRequest("POST", "/courier", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		controller.Create(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestUpdate(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := usecase.NewMockcourierUseCase(ctrl)
	controller := New(mockUseCase)

	t.Run("Success", func(t *testing.T) {

		updateRequest := dto.CourierUpdateRequest{
			ID:            1,
			Name:          "Masha",
			Phone:         "+71112223345",
			Status:        "busy",
			TransportType: "scooter",
		}

		mockUseCase.EXPECT().
			UpdateCourier(gomock.Any(), gomock.Any()).
			Return(nil)

		body, _ := json.Marshal(updateRequest)
		req := httptest.NewRequest("PUT", "/courier", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		controller.Update(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("InvalidJSON", func(t *testing.T) {

		req := httptest.NewRequest("PUT", "/courier", bytes.NewReader([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		controller.Update(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("CourierNotFound", func(t *testing.T) {

		updateRequest := dto.CourierUpdateRequest{
			ID:    999,
			Name:  "Kim",
			Phone: "+79273890022",
		}

		mockUseCase.EXPECT().
			UpdateCourier(gomock.Any(), gomock.Any()).
			Return(domain.ErrCourierNotFound)

		body, _ := json.Marshal(updateRequest)
		req := httptest.NewRequest("PUT", "/courier", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		controller.Update(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("PhoneExists", func(t *testing.T) {

		updateRequest := dto.CourierUpdateRequest{
			ID:    1,
			Name:  "Alex",
			Phone: "+70009998876",
		}

		mockUseCase.EXPECT().
			UpdateCourier(gomock.Any(), gomock.Any()).
			Return(domain.ErrPhoneExists)

		body, _ := json.Marshal(updateRequest)
		req := httptest.NewRequest("PUT", "/courier", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		controller.Update(rec, req)

		assert.Equal(t, http.StatusConflict, rec.Code)
	})

	t.Run("MissingRequiredFields", func(t *testing.T) {

		updateRequest := dto.CourierUpdateRequest{
			ID:    1,
			Name:  "",
			Phone: "+79997777787",
		}

		mockUseCase.EXPECT().
			UpdateCourier(gomock.Any(), gomock.Any()).
			Return(domain.ErrMissingRequiredFields)

		body, _ := json.Marshal(updateRequest)
		req := httptest.NewRequest("PUT", "/courier", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		controller.Update(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("InvalidPhone", func(t *testing.T) {

		updateRequest := dto.CourierUpdateRequest{
			ID:    1,
			Name:  "Sam",
			Phone: "not a number",
		}

		mockUseCase.EXPECT().
			UpdateCourier(gomock.Any(), gomock.Any()).
			Return(domain.ErrInvalidPhone)

		body, _ := json.Marshal(updateRequest)
		req := httptest.NewRequest("PUT", "/courier", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		controller.Update(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("InternalServerError", func(t *testing.T) {

		updateRequest := dto.CourierUpdateRequest{
			ID:    1,
			Name:  "Max",
			Phone: "+77777777777",
		}

		mockUseCase.EXPECT().
			UpdateCourier(gomock.Any(), gomock.Any()).
			Return(errors.New("unexpected error"))

		body, _ := json.Marshal(updateRequest)
		req := httptest.NewRequest("PUT", "/courier", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		controller.Update(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}
