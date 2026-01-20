package delivery

import (
	http2 "avito/internal/handlers/http"
	dto2 "avito/internal/handlers/http/delivery/dto"
	"avito/internal/handlers/http/httperror"
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
	var req dto2.DeliveryCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http2.WriteErrorResponse(w, httperror.ErrInvalidJSON)
		return
	}

	domainDelivery := dto2.ToModelCreate(&req)
	delivery, err := c.useCase.AssignDelivery(r.Context(), domainDelivery.OrderId)

	if err != nil {
		http2.WriteErrorResponse(w, err)
		return
	}

	response := dto2.ToCreateResponse(*delivery)

	http2.WriteResponse(w, http.StatusOK, response)
}

func (c *DeliveryController) Delete(w http.ResponseWriter, r *http.Request) {
	var req dto2.DeliveryDeleteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http2.WriteErrorResponse(w, httperror.ErrInvalidJSON)
		return
	}

	domainDelivery := dto2.ToModelDelete(&req)
	delivery, err := c.useCase.UnassignDelivery(r.Context(), domainDelivery.OrderId)

	if err != nil {
		http2.WriteErrorResponse(w, err)
		return
	}

	response := dto2.ToDeleteResponse(*delivery)

	http2.WriteResponse(w, http.StatusOK, response)
}
