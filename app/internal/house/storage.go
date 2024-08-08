package house

import (
	"context"
)

type Repository interface {
	Create(ctx context.Context, h CreateHouseDTO) (House, error)
}
