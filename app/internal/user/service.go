package user

import (
	"context"
	"math/rand"
	"strings"
	"time"

	"github.com/Polyrom/houses_api/pkg/logging"
)

const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

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

func (s *Service) GenerateRandomEmailPrefix(ctx context.Context, length int) string {
	var builder strings.Builder
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		builder.WriteString(string(chars[r.Intn(len(chars))]))
	}
	return builder.String()
}

func NewService(r Repository, l logging.Logger) *Service {
	return &Service{repo: r, logger: l}
}
