package main

import (
	"RedisService/internal/config"
	"RedisService/internal/database"
	_ "RedisService/internal/database"
	"RedisService/src/api/events"
	"RedisService/src/api/redis"
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

	server.POST("/users/register", func(context *gin.Context) {
		// TODO: Добавить обработку запроса на регистрацию пользователя
	})

	server.POST("/users/authorize", func(context *gin.Context) {
		// TODO: Добавить обработку запроса на авторизацию пользователя
	})

	server.GET("/users/:id/info", func(context *gin.Context) {
		// TODO: Добавить обработку запроса на получение информации о пользователе
	})

	server.PUT("/users/:id/update", func(context *gin.Context) {
		// TODO: Добавить обработку запроса на обновление информации о пользователе
	})

	server.POST("/user/:id/reset_password", func(context *gin.Context) {
		// TODO: Добавить обработку запроса на сброс пароля пользователя
	})

	server.POST("/redis/user/:id/waiting_list", redis_worker.AddUserToWaitingListRequest)

	server.POST("/redis/next_user", redis_worker.ProcessNextUserRequest)

	server.POST("/events/create", func(context *gin.Context) {
		events.CreateEventRequest(context, db)
	})

	server.GET("/users/:id/events", func(context *gin.Context) {
		events.GetUserEventsRequest(context, db)
	})

	server.POST("/event/:id/register", func(context *gin.Context) {
		events.AttendEventRequest(context, db)
	})

	server.DELETE("/event/:id/cancel", func(context *gin.Context) {
		events.CancelVisitRequest(context, db)
	})

	server.GET("/events/list", func(context *gin.Context) {
		// TODO: Добавить обработку запроса на получение списка событий
	})

	server.GET("/events/:id/participants", func(context *gin.Context) {
		// TODO: Добавить обработку запроса на получение списка участников события
	})

	err := server.Run(cfg.HTTPServer.Address)
	if err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}
