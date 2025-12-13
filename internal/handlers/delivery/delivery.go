package delivery

import (
	"avito/internal/handlers"
	"avito/internal/handlers/delivery/dto"
	"avito/internal/handlers/httperror"
	"encoding/json"
	"net/http"
)

type DeliveryController struct {
	useCase deliveryUseCase
}

func New(useCase deliveryUseCase) *DeliveryController {
	return &DeliveryController{useCase: useCase}
}

func (c *DeliveryController) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.DeliveryCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handlers.WriteErrorResponse(w, httperror.ErrInvalidJSON)
		return
	}

	domainDelivery := dto.ToModelCreate(&req)
	delivery, err := c.useCase.AssignDelivery(r.Context(), domainDelivery.OrderId)

	if err != nil {
		handlers.WriteErrorResponse(w, err)
		return
	}

	response := dto.ToCreateResponse(*delivery)

	handlers.WriteResponse(w, http.StatusOK, response)
}

func (c *DeliveryController) Delete(w http.ResponseWriter, r *http.Request) {
	var req dto.DeliveryDeleteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handlers.WriteErrorResponse(w, httperror.ErrInvalidJSON)
		return
	}

	domainDelivery := dto.ToModelDelete(&req)
	delivery, err := c.useCase.UnassignDelivery(r.Context(), domainDelivery.OrderId)

	if err != nil {
		handlers.WriteErrorResponse(w, err)
		return
	}

	response := dto.ToDeleteResponse(*delivery)

	handlers.WriteResponse(w, http.StatusOK, response)
}
