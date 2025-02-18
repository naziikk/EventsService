package main

import (
	"RedisService/internal/config"
	"RedisService/internal/database"
	_ "RedisService/internal/database"
	"RedisService/src/api/events"
	"RedisService/src/api/redis"
	"RedisService/src/api/user"
	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5"
	"log"
)

func main() {
	cfg := config.MustLoadConfig()   // Получили конфиги
	db, _ := database.ConnectDB(cfg) // Подключились к базе данных
	defer db.Close()                 // Не забываем закрыть соединение с базой данных по окончании
	redis_worker.InitRedis(cfg)      // Подключились к Redis
	defer redis_worker.Rdb.Close()   // Не забываем закрыть соединение с Redis по окончании

	server := gin.Default()

	server.POST("/redis/user/:id/waiting_list", redis_worker.AddUserToWaitingListRequest)

	server.POST("/redis/next_user", redis_worker.ProcessNextUserRequest)

	server.POST("/create_event", func(context *gin.Context) {
		events.CreateEventRequest(context, db)
	})

	server.GET("/user/:id/events", func(context *gin.Context) {
		events.GetUserEventsRequest(context, db)
	})

	server.POST("/user/:id/visit_event", func(context *gin.Context) {
		user.AttendEventRequest(context, db)
	})

	server.DELETE("/user/:id/cancel_visit", func(context *gin.Context) {
		user.CancelVisitRequest(context, db)
	})
	err := server.Run(cfg.HTTPServer.Address)
	if err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}
