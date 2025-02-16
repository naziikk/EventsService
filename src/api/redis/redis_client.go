package redis_worker

import (
	"RedisService/internal/config"
	"context"
	"github.com/go-redis/redis/v8"
	"log"
)

var Rdb *redis.Client
var Ctx = context.Background()

func InitRedis(cfg *config.Config) {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     cfg.RedisServer.Address,
		Password: cfg.RedisServer.Password,
		DB:       cfg.RedisServer.DB,
	})

	_, err := Rdb.Ping(Ctx).Result()
	if err != nil {
		log.Fatalf("Ошибка при подключении к Redis: %v", err)
	}
	log.Println("Подключение к Redis прошло успешно")
}
