package user

import (
	"context"
	"time"

	"errors"

	"github.com/Polyrom/houses_api/pkg/client/postgres"
	"github.com/Polyrom/houses_api/pkg/logging"
	"github.com/jackc/pgx/v5/pgconn"
)

type repository struct {
	client postgres.Client
	logger logging.Logger
}

func (r *repository) Create(ctx context.Context, u User) (UserID, error) {
	q := `INSERT INTO users 
					(email, password, role) 
				VALUES 
					($1, $2, $3)
				RETURNING 
					id`
	var userid UserID
	err := r.client.QueryRow(ctx, q, u.Email, u.Password, u.Role).Scan(&userid)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.Is(err, pgErr) {
			pgErr = err.(*pgconn.PgError)
			r.logger.Errorf("SQL Error: %s, Detail: %s, Where: %s", pgErr.Message, pgErr.Detail, pgErr.Where)
			return "", pgErr
		}
		return "", err
	}
	return userid, nil
}

func (r *repository) GetByID(ctx context.Context, uid UserID) (User, error) {
	q := `SELECT 
						id, email, password, role
					FROM 
						users 
					WHERE 
						id = $1`
	var u User
	err := r.client.QueryRow(ctx, q, uid).Scan(&u.ID, &u.Email, &u.Password, &u.Role)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.Is(err, pgErr) {
			pgErr = err.(*pgconn.PgError)
			r.logger.Errorf("SQL Error: %s, Detail: %s, Where: %s", pgErr.Message, pgErr.Detail, pgErr.Where)
			return User{}, pgErr
		}
		return User{}, err
	}
	return u, nil
}

func (r *repository) AddToken(ctx context.Context, uid UserID, token Token) error {
	q := `INSERT INTO tokens 
					(user_id, token, expires_at)
				VALUES
					($1, $2, $3) 
				ON CONFLICT 
					(user_id) 
				DO UPDATE SET
					token = $2,
					expires_at = $3`
	_, err := r.client.Exec(ctx, q, uid, token, time.Now().Add(1*time.Hour))
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.Is(err, pgErr) {
			pgErr = err.(*pgconn.PgError)
			r.logger.Errorf("SQL Error: %s, Detail: %s, Where: %s", pgErr.Message, pgErr.Detail, pgErr.Where)
			return pgErr
		}
		return err
	}
	return nil
}

func NewRepository(c postgres.Client, l logging.Logger) Repository {
	return &repository{
		client: c,
		logger: l,
	}
}
