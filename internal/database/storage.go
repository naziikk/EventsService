package database

import (
	"RedisService/internal/config"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
	"time"
)

func ConnectDB(cfg *config.Config) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf("postgres://%s/%s?sslmode=disable", cfg.PostgresData.Address, cfg.PostgresData.Name)
	db, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func InitDB(db *pgxpool.Pool) error {
	sqlFile := "internal/database/migrations/init_db.sql"
	query, err := os.ReadFile(sqlFile)
	if err != nil {
		return fmt.Errorf("ошибка чтения файла миграции: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = db.Exec(ctx, string(query))
	if err != nil {
		return fmt.Errorf("ошибка выполнения миграции: %w", err)
	}

	fmt.Println("База данных успешно проинициализирована!")
	return nil
}
