package main

import (
	"RedisService/internal/config"
	"RedisService/internal/database"
	_ "RedisService/internal/database"
	"RedisService/src/handlers/authorization"
	"RedisService/src/handlers/events"
	"RedisService/src/handlers/middleware"
	"RedisService/src/handlers/redis"
	"RedisService/src/handlers/user"
	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5"
	"log"
)

func main() {
	cfg := config.MustLoadConfig()   // Получили конфиги
	db, _ := database.ConnectDB(cfg) // Подключились к базе данных
	defer db.Close()                 // Не забываем закрыть соединение с базой данных по окончании
	if err := database.InitDB(db); err != nil {
		log.Fatalf("Ошибка инициализации БД: %v", err)
	}
	redis_worker.InitRedis(cfg)    // Подключились к Redis
	defer redis_worker.Rdb.Close() // Не забываем закрыть соединение с Redis по окончании

	//notifications.SendEmail("zacnazar20@gmail.com", "Test", "Test")
	router := gin.Default()

	router.Use(middleware.LoggingMiddleware())
	router.Use(middleware.AuthMiddleware(cfg))

	userRoutes := router.Group("/user")
	{
		userRoutes.POST("/register", func(ctx *gin.Context) {
			authorization.LoginRequest(ctx, db)
		})
		userRoutes.POST("/authorize", func(ctx *gin.Context) {
			authorization.UserAuthorizationRequest(ctx, db, cfg)
		})
		userRoutes.GET("/me", func(ctx *gin.Context) {
			user.GetUserInfoRequest(ctx, db)
		})
		userRoutes.PUT("/update", func(ctx *gin.Context) {
			user.UpdateUserInfoRequest(ctx, db)
		})
		userRoutes.POST("/reset_password", func(ctx *gin.Context) {
			authorization.ResetPasswordRequest(ctx, db)
		})
		userRoutes.GET("/events", func(ctx *gin.Context) {
			events.GetUserEventsRequest(ctx, db)
		})
	}

	redisRoutes := router.Group("/redis")
	{
		redisRoutes.POST("/user/waiting_list", redis_worker.AddUserToWaitingListRequest)
		redisRoutes.POST("/next_user", redis_worker.ProcessNextUserRequest)
	}

	eventRoutes := router.Group("/events")
	{
		eventRoutes.POST("/create", func(ctx *gin.Context) {
			events.CreateEventRequest(ctx, db)
		})
		eventRoutes.GET("/list", func(ctx *gin.Context) {
			events.GetAllEventsRequest(ctx, db)
		})
		eventRoutes.GET("/:id/participants", func(ctx *gin.Context) {
			events.GetEventsParticipantsRequest(ctx, db)
		})
	}

	attendanceRoutes := router.Group("/event")
	{
		attendanceRoutes.POST("/:id/register", func(ctx *gin.Context) {
			events.AttendEventRequest(ctx, db)
		})
		attendanceRoutes.DELETE("/:id/cancel", func(ctx *gin.Context) {
			events.CancelVisitRequest(ctx, db)
		})
	}

	err := router.Run(cfg.HTTPServer.Address)
	if err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}
