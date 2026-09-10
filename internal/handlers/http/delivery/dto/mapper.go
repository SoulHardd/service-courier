package dto

import (
	"github.com/SoulHardd/service-courier/internal/domain"
)

func ToCreateResponse(d domain.Delivery) DeliveryCreateResponse {
	return DeliveryCreateResponse{
		CourierId:     d.CourierId,
		OrderId:       d.OrderId,
		TransportType: d.TransportType,
		Deadline:      d.Deadline,
	}
}

func ToDeleteResponse(d domain.Delivery) DeliveryDeleteResponse {
	return DeliveryDeleteResponse{
		OrderId:   d.OrderId,
		Status:    "unassigned",
		CourierId: d.CourierId,
	}
}

func ToModelCreate(req *DeliveryCreateRequest) domain.Delivery {
	return domain.Delivery{
		OrderId: req.OrderId,
	}
}

func ToModelDelete(req *DeliveryDeleteRequest) domain.Delivery {
	return domain.Delivery{
		OrderId: req.OrderId,
	}
}
