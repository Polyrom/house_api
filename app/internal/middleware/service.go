package middleware

import (
	"context"

	"github.com/Polyrom/houses_api/pkg/logging"
)

type Service struct {
	repo   Repository
	logger logging.Logger
}

func (s *Service) GetRoleByToken(ctx context.Context, token Token) (UserIDRoleDTO, error) {
	return s.repo.GetRoleByToken(ctx, token)
}

func NewService(r Repository, l logging.Logger) Service {
	return Service{repo: r, logger: l}
}
