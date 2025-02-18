package middleware

import (
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

func LoggingMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		start := time.Now()
		context.Next()
		finish := time.Since(start)

		log.Println("Обработан запрос ", context.Request.URL.Path, ". Время выполнения: ", finish)
	}
}
