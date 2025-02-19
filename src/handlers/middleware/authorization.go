package middleware

import (
	"RedisService/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
)

//var jwtKey = []byte("your-secret-key")

type Claims struct {
	Username string `json:"username"`
	UserID   string `json:"user_id"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}

func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(context *gin.Context) {
		token, err := context.Cookie("jwt_token")

		if err != nil {
			context.JSON(http.StatusUnauthorized, gin.H{"message": "Необходимо авторизоваться"})
			context.Abort()
			return
		}

		claims := &Claims{}
		jwtToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
			return cfg.JWTSecret, nil
		})
		if err != nil || !jwtToken.Valid {
			context.JSON(http.StatusUnauthorized, gin.H{"message": "Недействительный токен"})
			context.Abort()
			return
		}
		context.Set("userID", claims.UserID)
		context.Set("username", claims.Username)
		context.Set("email", claims.Email)
		context.Next()
	}
}
