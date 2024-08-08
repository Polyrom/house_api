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

func (r *repository) GetRoleByToken(ctx context.Context, token Token) (UserIDRoleDTO, error) {
	q := `SELECT u.id, u.role
				FROM tokens t
  			JOIN users u ON t.user_id = u.id
				WHERE t.token = $1;`
	var userIDRole UserIDRoleDTO
	err := r.client.QueryRow(ctx, q, token).Scan(&userIDRole.ID, &userIDRole.Role)
	if err != nil {
		return UserIDRoleDTO{}, err
	}
	return userIDRole, nil
}

func NewRepository(c postgres.Client, l logging.Logger) Repository {
	return &repository{client: c, l: l}
}
