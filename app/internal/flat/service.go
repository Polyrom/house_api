package flat

import (
	"context"
	"errors"

	"github.com/Polyrom/houses_api/internal/middleware"
	"github.com/Polyrom/houses_api/internal/modstatus"
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
	userID := ctx.Value(middleware.UserID).(string)
	fldto := GetFlatByIDDTO{ID: f.ID, HouseID: f.HouseID}
	storedFlat, err := s.repo.GetByID(ctx, fldto)
	if err != nil {
		userNotFoundErr := errors.New("user not found")
		return FlatDTO{}, userNotFoundErr
	}
	if storedFlat.Status == modstatus.Created.String() {
		if f.Status != modstatus.OnModeration.String() {
			canJumpStatusErr := errors.New("cannot approve/decline without moderation")
			return FlatDTO{}, canJumpStatusErr
		}
		return s.repo.UpdateWithNewMod(ctx, userID, f)
	}
	if storedFlat.Status == modstatus.OnModeration.String() && storedFlat.Moderator != userID {
		wrongModErr := errors.New("already taken by another moderator")
		return FlatDTO{}, wrongModErr
	}
	return s.repo.Update(ctx, f)
}

func NewService(r Repository, l logging.Logger) *Service {
	return &Service{repo: r, logger: l}
}
