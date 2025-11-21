package courier

import (
	"avito/internal/model"
	"context"
)

type CourierRepository interface {
	GetOneById(ctx context.Context, id int64) (*model.Courier, error)
	GetAll(ctx context.Context) ([]model.Courier, error)
	Create(ctx context.Context, courier *model.Courier) (int64, error)
	Update(ctx context.Context, courier *model.Courier) error
}
