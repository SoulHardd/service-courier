package dto

import "avito/internal/model"

func ToCourierResponse(c model.Courier) CourierResponse {
	return CourierResponse{
		ID:     c.ID,
		Name:   c.Name,
		Phone:  c.Phone,
		Status: c.Status,
	}
}

func ToCourierResponses(list []model.Courier) []CourierResponse {
	res := make([]CourierResponse, len(list))
	for i, c := range list {
		res[i] = ToCourierResponse(c)
	}
	return res
}

func ToModelCreate(req *CourierCreateRequest) model.Courier {
	return model.Courier{
		Name:   req.Name,
		Phone:  req.Phone,
		Status: req.Status,
	}
}

func ToModelUpdate(req *CourierUpdateRequest) model.Courier {
	return model.Courier{
		ID:     req.ID,
		Name:   req.Name,
		Phone:  req.Phone,
		Status: req.Status,
	}
}
