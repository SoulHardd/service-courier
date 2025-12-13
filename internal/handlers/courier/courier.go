package courier

import (
	"avito/internal/handlers"
	"avito/internal/handlers/courier/dto"
	"avito/internal/handlers/httperror"
	"encoding/json"
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
		handlers.WriteErrorResponse(w, httperror.ErrInvalidCourierId)
		return
	}

	Courier, err := c.useCase.GetCourier(r.Context(), id)

	if err != nil {
		handlers.WriteErrorResponse(w, err)
		return
	}

	response := dto.ToCourierResponse(*Courier)

	handlers.WriteResponse(w, http.StatusOK, response)
}

func (c *CourierController) GetAll(w http.ResponseWriter, r *http.Request) {
	Couriers, err := c.useCase.GetCouriers(r.Context())
	if err != nil {
		handlers.WriteErrorResponse(w, err)
		return
	}

	response := dto.ToCourierResponses(Couriers)

	handlers.WriteResponse(w, http.StatusOK, response)
}

func (c *CourierController) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CourierCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handlers.WriteErrorResponse(w, httperror.ErrInvalidJSON)
		return
	}

	domainCourier := dto.ToModelCreate(&req)
	id, err := c.useCase.CreateCourier(r.Context(), &domainCourier)

	if err != nil {
		handlers.WriteErrorResponse(w, err)
		return
	}

	response := map[string]interface{}{
		"id":      id,
		"message": "Courier created successfully",
	}

	handlers.WriteResponse(w, http.StatusOK, response)
}

func (c *CourierController) Update(w http.ResponseWriter, r *http.Request) {
	var req dto.CourierUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handlers.WriteErrorResponse(w, httperror.ErrInvalidJSON)
		return
	}

	domainCourier := dto.ToModelUpdate(&req)
	err := c.useCase.UpdateCourier(r.Context(), &domainCourier)

	if err != nil {
		handlers.WriteErrorResponse(w, err)
		return
	}

	response := map[string]string{
		"message": "Courier updated successfully",
	}

	handlers.WriteResponse(w, http.StatusOK, response)
}
