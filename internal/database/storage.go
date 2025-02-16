package database

import (
	"RedisService/internal/config"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectDB(cfg *config.Config) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf("postgres://%s/%s?sslmode=disable", cfg.PostgresData.Address, cfg.PostgresData.Name)
	db, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}
