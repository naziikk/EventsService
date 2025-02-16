package main

import (
	"RedisService/src/api/events"
	"RedisService/src/api/redis"
	"context"
	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

func connectDB() (*pgxpool.Pool, error) {
	dsn := "postgres://username:password@localhost:5432/events_service_data"
	db, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func main() {
	db, err := connectDB()
	redis_worker.InitRedis()
	defer redis_worker.Rdb.Close()

	server := gin.Default()

	server.POST("/redis/user/:id/waiting_list", redis_worker.AddUserToWaitingListRequest)

	server.POST("/redis/next_user", redis_worker.ProcessNextUserRequest)

	server.POST("/create_event", func(context *gin.Context) {
		events.CreateEventRequest(context, db)
	})

	server.PUT("/organizer/update_event")

	err = server.Run("0.0.0.0:8010")
	if err != nil {

		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}
