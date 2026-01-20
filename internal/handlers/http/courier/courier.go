package courier

import (
	http2 "avito/internal/handlers/http"
	dto2 "avito/internal/handlers/http/courier/dto"
	"avito/internal/handlers/http/httperror"
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
		http2.WriteErrorResponse(w, httperror.ErrInvalidCourierId)
		return
	}

	Courier, err := c.useCase.GetCourier(r.Context(), id)

	if err != nil {
		http2.WriteErrorResponse(w, err)
		return
	}

	response := dto2.ToCourierResponse(*Courier)

	http2.WriteResponse(w, http.StatusOK, response)
}

func (c *CourierController) GetAll(w http.ResponseWriter, r *http.Request) {
	Couriers, err := c.useCase.GetCouriers(r.Context())
	if err != nil {
		http2.WriteErrorResponse(w, err)
		return
	}

	response := dto2.ToCourierResponses(Couriers)

	http2.WriteResponse(w, http.StatusOK, response)
}

func (c *CourierController) Create(w http.ResponseWriter, r *http.Request) {
	var req dto2.CourierCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http2.WriteErrorResponse(w, httperror.ErrInvalidJSON)
		return
	}

	domainCourier := dto2.ToModelCreate(&req)
	id, err := c.useCase.CreateCourier(r.Context(), &domainCourier)

	if err != nil {
		http2.WriteErrorResponse(w, err)
		return
	}

	response := map[string]interface{}{
		"id":      id,
		"message": "Courier created successfully",
	}

	http2.WriteResponse(w, http.StatusOK, response)
}

func (c *CourierController) Update(w http.ResponseWriter, r *http.Request) {
	var req dto2.CourierUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http2.WriteErrorResponse(w, httperror.ErrInvalidJSON)
		return
	}

	domainCourier := dto2.ToModelUpdate(&req)
	err := c.useCase.UpdateCourier(r.Context(), &domainCourier)

	if err != nil {
		http2.WriteErrorResponse(w, err)
		return
	}

	response := map[string]string{
		"message": "Courier updated successfully",
	}

	http2.WriteResponse(w, http.StatusOK, response)
}
