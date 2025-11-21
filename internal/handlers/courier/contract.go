package courier

import (
	"avito/internal/model"
	"context"
)

type courierUseCase interface {
	GetCourier(ctx context.Context, id int64) (*model.Courier, error)
	GetCouriers(ctx context.Context) ([]model.Courier, error)
	CreateCourier(ctx context.Context, req *model.Courier) (int64, error)
	UpdateCourier(ctx context.Context, req *model.Courier) error
}
