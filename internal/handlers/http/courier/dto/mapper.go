package dto

import (
	"github.com/SoulHardd/service-courier/internal/domain"
)

func ToCourierResponse(c domain.Courier) CourierResponse {
	return CourierResponse{
		ID:            c.ID,
		Name:          c.Name,
		Phone:         c.Phone,
		Status:        c.Status,
		TransportType: c.TransportType,
	}
}

func ToCourierResponses(list []domain.Courier) []CourierResponses {
	res := make([]CourierResponses, len(list))
	for i, c := range list {
		res[i] = CourierResponses{
			ID:            c.ID,
			Name:          c.Name,
			Status:        c.Status,
			TransportType: c.TransportType,
		}
	}
	return res
}

func ToModelCreate(req *CourierCreateRequest) domain.Courier {
	return domain.Courier{
		Name:          req.Name,
		Phone:         req.Phone,
		Status:        req.Status,
		TransportType: req.TransportType,
	}
}

func ToModelUpdate(req *CourierUpdateRequest) domain.Courier {
	return domain.Courier{
		ID:            req.ID,
		Name:          req.Name,
		Phone:         req.Phone,
		Status:        req.Status,
		TransportType: req.TransportType,
	}
}
