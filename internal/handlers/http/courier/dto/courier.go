package dto

import "github.com/SoulHardd/service-courier/internal/domain"

type CourierResponse struct {
	ID            int64                `json:"id"`
	Name          string               `json:"name"`
	Phone         string               `json:"phone"`
	Status        domain.CourierStatus `json:"status"`
	TransportType domain.TransportType `json:"transport_type"`
}

type CourierResponses struct {
	ID            int64                `json:"id"`
	Name          string               `json:"name"`
	Status        domain.CourierStatus `json:"status"`
	TransportType domain.TransportType `json:"transport_type"`
}

type CourierCreateRequest struct {
	Name          string               `json:"name"`
	Phone         string               `json:"phone"`
	Status        domain.CourierStatus `json:"status"`
	TransportType domain.TransportType `json:"transport_type"`
}

type CourierUpdateRequest struct {
	ID            int64                `json:"id"`
	Name          string               `json:"name"`
	Phone         string               `json:"phone"`
	Status        domain.CourierStatus `json:"status"`
	TransportType domain.TransportType `json:"transport_type"`
}
