package user

import (
	"context"

	"github.com/Polyrom/houses_api/pkg/logging"
)

type Service struct {
	repo   Repository
	logger logging.Logger
}

func (s *Service) Register(ctx context.Context, u User) (UserID, error) {
	err := u.HashPassword(u.Password)
	if err != nil {
		return "", err
	}
	return s.repo.Create(ctx, u)
}

func (s *Service) GetByID(ctx context.Context, uid UserID) (User, error) {
	return s.repo.GetByID(ctx, uid)
}

func (s *Service) AddToken(ctx context.Context, uid UserID, token Token) error {
	return s.repo.AddToken(ctx, uid, token)
}

func NewService(r Repository, l logging.Logger) *Service {
	return &Service{repo: r, logger: l}
}
