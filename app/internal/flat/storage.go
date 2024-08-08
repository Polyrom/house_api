package flat

import "context"

type Repository interface {
	GetByHouseIDModerator(ctx context.Context, fl FlatID) ([]FlatDTO, error)
	GetByHouseIDClient(ctx context.Context, fl FlatID) ([]FlatDTO, error)
	GetByID(ctx context.Context, fl GetFlatByIDDTO) (FlatDTO, error)
	Create(ctx context.Context, fl CreateFlatDTO) (FlatDTO, error)
	Update(ctx context.Context, fl UpdateFlatStatusDTO) (FlatDTO, error)
	UpdateWithNewMod(ctx context.Context, uid string, fl UpdateFlatStatusDTO) (FlatDTO, error)
}
