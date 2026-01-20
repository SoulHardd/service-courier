package dto

import (
	"avito/internal/domain"
	"time"
)

type DeliveryCreateRequest struct {
	OrderId string `json:"order_id"`
}

type DeliveryCreateResponse struct {
	CourierId     int64                `json:"courier_id"`
	OrderId       string               `json:"order_id"`
	TransportType domain.TransportType `json:"transport_type"`
	Deadline      time.Time            `json:"delivery_deadline"`
}

type DeliveryDeleteRequest struct {
	OrderId string `json:"order_id"`
}

type DeliveryDeleteResponse struct {
	OrderId   string `json:"order_id"`
	Status    string `json:"status"`
	CourierId int64  `json:"courier_id"`
}
