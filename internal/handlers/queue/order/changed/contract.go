package changed

import (
	"context"

	model "github.com/SoulHardd/service-courier/internal/domain"
)

type OrderUseCase interface {
	Process(ctx context.Context, o model.Order) error
}
