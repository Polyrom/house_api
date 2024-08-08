package middleware

import (
	"context"

	"github.com/Polyrom/houses_api/pkg/client/postgres"
	"github.com/Polyrom/houses_api/pkg/logging"
)

type repository struct {
	client postgres.Client
	l      logging.Logger
}

func (r *repository) GetRoleByToken(ctx context.Context, token Token) (Role, error) {
	q := `SELECT u.role
				FROM tokens t
  			JOIN users u ON t.user_id = u.id
				WHERE t.token = $1;`
	var role string
	err := r.client.QueryRow(ctx, q, token).Scan(&role)
	if err != nil {
		return Role(""), err
	}
	return Role(role), nil
}

func NewRepository(c postgres.Client, l logging.Logger) Repository {
	return &repository{client: c, l: l}
}
