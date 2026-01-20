package changed

import (
	model "avito/internal/domain"
	"context"
)

type OrderUseCase interface {
	Process(ctx context.Context, o model.Order) error
}
