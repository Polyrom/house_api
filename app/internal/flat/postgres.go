package flat

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

func (r *repository) GetByHouseIDClient(ctx context.Context, fl FlatID) ([]FlatDTO, error) {
	q := `SELECT 
					id, house_id, price, rooms, status 
				FROM 
					flats 
				WHERE house_id = $1
				AND status = 'approved'`

	rows, err := r.client.Query(ctx, q, fl)
	if err != nil {
		return nil, err
	}
	fls := make([]FlatDTO, 0)
	for rows.Next() {
		var f FlatDTO
		err = rows.Scan(&f.ID, &f.HouseID, &f.Price, &f.Rooms, &f.Status)
		if err != nil {
			return nil, err
		}
		fls = append(fls, f)
	}
	return fls, nil
}

func (r *repository) GetByHouseIDModerator(ctx context.Context, fl FlatID) ([]FlatDTO, error) {
	q := `SELECT 
					id, house_id, price, rooms, status 
				FROM 
					flats 
				WHERE house_id = $1`

	rows, err := r.client.Query(ctx, q, fl)
	if err != nil {
		return nil, err
	}
	fls := make([]FlatDTO, 0)
	for rows.Next() {
		var f FlatDTO
		err = rows.Scan(&f.ID, &f.HouseID, &f.Price, &f.Rooms, &f.Status)
		if err != nil {
			return nil, err
		}
		fls = append(fls, f)
	}
	return fls, nil
}

func (r *repository) Create(ctx context.Context, fl CreateFlatDTO) (FlatDTO, error) {
	q := `INSERT INTO flats 
					(house_id, price, rooms) 
				VALUES 
					($1, $2, $3)
				RETURNING 
					id, house_id, price, rooms, status`
	var f FlatDTO
	err := r.client.QueryRow(ctx, q, fl.HouseID, fl.Price, fl.Rooms).
		Scan(&f.ID, &f.HouseID, &f.Price, &f.Rooms, &f.Status)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.Is(err, pgErr) {
			pgErr = err.(*pgconn.PgError)
			r.logger.Errorf("SQL Error: %s, Detail: %s, Where: %s", pgErr.Message, pgErr.Detail, pgErr.Where)
			return f, pgErr
		}
		return f, err
	}
	return f, nil
}

func (r *repository) Update(ctx context.Context, fl UpdateFlatStatusDTO) (FlatDTO, error) {
	q := `UPDATE flats 
				SET status = $1
				WHERE id = $2
				AND house_id = $3
				RETURNING 
					id, house_id, price, rooms, status`
	var f FlatDTO
	err := r.client.QueryRow(ctx, q, fl.Status, fl.ID, fl.HouseID).Scan(&f.ID, &f.HouseID, &f.Price, &f.Rooms, &f.Status)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.Is(err, pgErr) {
			pgErr = err.(*pgconn.PgError)
			r.logger.Errorf("SQL Error: %s, Detail: %s, Where: %s", pgErr.Message, pgErr.Detail, pgErr.Where)
			return f, pgErr
		}
		return f, err
	}
	return f, nil
}

func NewRepository(c postgres.Client, l logging.Logger) Repository {
	return &repository{
		client: c,
		logger: l,
	}
}
