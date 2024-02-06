package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jay-bhogayata/blogapi/logger"
)

func Init(ctx context.Context, url string) (*pgxpool.Pool, error) {

	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		logger.Log.Error("error while creating db connection pool", "err", err)
		return nil, err
	}

	err = pool.Ping(ctx)
	if err != nil {
		logger.Log.Error("error while piing the db", "err", err)
		return nil, err
	}

	logger.Log.Info("connected to database successfully")

	return pool, err

}
