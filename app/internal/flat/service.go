package flat

import (
	"context"

	"github.com/Polyrom/houses_api/internal/middleware"
	"github.com/Polyrom/houses_api/pkg/logging"
)

type Service struct {
	repo   Repository
	logger logging.Logger
}

func (s *Service) GetByHouseID(ctx context.Context, f FlatID) ([]FlatDTO, error) {
	userRole := ctx.Value(middleware.UserRole).(middleware.Role)
	if userRole == middleware.Moderator {
		return s.repo.GetByHouseIDModerator(ctx, f)
	}
	return s.repo.GetByHouseIDClient(ctx, f)
}

func (s *Service) Create(ctx context.Context, f CreateFlatDTO) (FlatDTO, error) {
	return s.repo.Create(ctx, f)
}

func (s *Service) Update(ctx context.Context, f UpdateFlatStatusDTO) (FlatDTO, error) {
	return s.repo.Update(ctx, f)
}

func NewService(r Repository, l logging.Logger) *Service {
	return &Service{repo: r, logger: l}
}
