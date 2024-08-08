package middleware

import "context"

type Repository interface {
	GetRoleByToken(ctx context.Context, token Token) (UserIDRoleDTO, error)
}
