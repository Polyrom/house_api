package house

import (
	"context"

	"github.com/Polyrom/houses_api/pkg/logging"
)

type Service struct {
	repo   Repository
	logger logging.Logger
}

func (s *Service) Create(ctx context.Context, h CreateHouseDTO) (House, error) {
	return s.repo.Create(ctx, h)
}

func NewService(r Repository, l logging.Logger) *Service {
	return &Service{repo: r, logger: l}
}
