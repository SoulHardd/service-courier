package courier

import (
	"avito/internal/domain"
	"avito/internal/handlers"
	dto2 "avito/internal/handlers/courier/dto"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type CourierController struct {
	useCase courierUseCase
}

func New(useCase courierUseCase) *CourierController {
	return &CourierController{useCase: useCase}
}

func (c *CourierController) Get(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, `{"error": "Invalid Courier ID"}`, http.StatusBadRequest)
		return
	}

	Courier, err := c.useCase.GetCourier(r.Context(), id)

	if err != nil {
		switch {
		case errors.Is(err, domain.ErrCourierNotFound):
			http.Error(w, `{"error": "Courier not found"}`, http.StatusNotFound)
		default:
			http.Error(w, `{"error": "Database error"}`, http.StatusInternalServerError)
		}
		return
	}

	response := dto2.ToCourierResponse(*Courier)

	handlers.WriteResponse(w, http.StatusOK, response)
}

func (c *CourierController) GetAll(w http.ResponseWriter, r *http.Request) {
	Couriers, err := c.useCase.GetCouriers(r.Context())
	if err != nil {
		http.Error(w, `{"error": "Database error"}`, http.StatusInternalServerError)
		return
	}

	response := dto2.ToCourierResponses(Couriers)

	handlers.WriteResponse(w, http.StatusOK, response)
}

func (c *CourierController) Create(w http.ResponseWriter, r *http.Request) {
	var req dto2.CourierCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Invalid JSON"}`, http.StatusBadRequest)
		return
	}

	domainCourier := dto2.ToModelCreate(&req)
	id, err := c.useCase.CreateCourier(r.Context(), &domainCourier)

	if err != nil {
		switch {
		case errors.Is(err, domain.ErrMissingRequiredFields):
			http.Error(w, `{"error": "Missing required fields"}`, http.StatusBadRequest)
		case errors.Is(err, domain.ErrInvalidPhone):
			http.Error(w, `{"error": "Invalid phone"}`, http.StatusBadRequest)
		case errors.Is(err, domain.ErrPhoneExists):
			http.Error(w, `{"error": "Courier with this phone already exists"}`, http.StatusConflict)
		default:
			http.Error(w, `{"error": "Database error"}`, http.StatusInternalServerError)
		}
		return
	}

	response := map[string]interface{}{
		"id":      id,
		"message": "Courier created successfully",
	}

	handlers.WriteResponse(w, http.StatusOK, response)
}

func (c *CourierController) Update(w http.ResponseWriter, r *http.Request) {
	var req dto2.CourierUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Invalid JSON"}`, http.StatusBadRequest)
		return
	}

	domainCourier := dto2.ToModelUpdate(&req)
	err := c.useCase.UpdateCourier(r.Context(), &domainCourier)

	if err != nil {
		switch {
		case errors.Is(err, domain.ErrMissingRequiredFields):
			http.Error(w, `{"error": "Missing required fields"}`, http.StatusBadRequest)
		case errors.Is(err, domain.ErrInvalidPhone):
			http.Error(w, `{"error": "Invalid phone"}`, http.StatusBadRequest)
		case errors.Is(err, domain.ErrPhoneExists):
			http.Error(w, `{"error": "Courier with this phone already exists"}`, http.StatusConflict)
		case errors.Is(err, domain.ErrCourierNotFound):
			http.Error(w, `{"error": "Courier not found"}`, http.StatusNotFound)
		default:
			http.Error(w, `{"error": "Database error"}`, http.StatusInternalServerError)
		}
		return
	}

	response := map[string]string{
		"message": "Courier updated successfully",
	}

	handlers.WriteResponse(w, http.StatusOK, response)
}
