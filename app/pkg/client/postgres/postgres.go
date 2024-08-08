package postgres

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Polyrom/houses_api/internal/config"
	"github.com/Polyrom/houses_api/pkg/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Client interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
}

func NewClient(ctx context.Context, sc config.StorageConfig) (*pgxpool.Pool, error) {
	var pool *pgxpool.Pool
	var err error

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", sc.Username, sc.Password, sc.Host, sc.Port, sc.Database)
	err = utils.Repeat(func() error {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		pool, err = pgxpool.New(ctx, dsn)
		if err != nil {
			return err
		}
		return nil
	}, sc.MaxAttempts, 5*time.Second)

	if err != nil {
		log.Fatal("repeat connect to postgresql")
	}

	return pool, nil
}
