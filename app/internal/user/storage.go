package user

import "context"

type Repository interface {
	Create(ctx context.Context, u User) (UserID, error)
	GetByID(ctx context.Context, uid UserID) (User, error)
	AddToken(ctx context.Context, uid UserID, token Token) error
}
