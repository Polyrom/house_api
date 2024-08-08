package house

import (
	"context"

	"errors"

	"github.com/Polyrom/houses_api/pkg/client/postgres"
	"github.com/Polyrom/houses_api/pkg/logging"
	"github.com/jackc/pgx/v5/pgconn"
)

type repository struct {
	client postgres.Client
	logger logging.Logger
}

func (r *repository) Create(ctx context.Context, h CreateHouseDTO) (House, error) {
	q := `INSERT INTO houses 
					(address, year, developer) 
				VALUES 
					($1, $2, $3)
				RETURNING 
					id, address, year, developer, created_at, update_at`
	var nh House
	err := r.client.QueryRow(ctx, q, h.Address, h.Year, h.Developer).
		Scan(&nh.ID, &nh.Address, &nh.Year, &nh.Developer, &nh.CreatedAt, &nh.UpdateAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.Is(err, pgErr) {
			pgErr = err.(*pgconn.PgError)
			r.logger.Errorf("SQL Error: %s, Detail: %s, Where: %s", pgErr.Message, pgErr.Detail, pgErr.Where)
			return House{}, pgErr
		}
		return House{}, err
	}
	return nh, nil
}

func NewRepository(c postgres.Client, l logging.Logger) Repository {
	return &repository{
		client: c,
		logger: l,
	}
}
