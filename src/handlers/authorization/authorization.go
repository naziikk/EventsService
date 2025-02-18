package authorization

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"time"
)

// TODO: добавить хеширование

var jwtKey = []byte("your-secret-key")

type Claims struct {
	Username string `json:"username"`
	UserID   string `json:"user_id"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}

func authorizationRequest(context *gin.Context, db *pgxpool.Pool) {
	var req UserData
	if err := context.ShouldBindJSON(&req); err != nil {
		log.Printf("Ошибка при парсинге JSON: %v", err)
		context.JSON(http.StatusBadRequest, gin.H{"message": "Неверный формат запроса"})
		return
	}

	if !validatePassword(db, req.Username, req.Password) {
		log.Printf("Неверное имя пользователя или пароль")
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Неверное имя пользователя или пароль"})
		return
	}

	Email, userId, err := getUsersEmail(db, req.Username)
	if err != nil {
		log.Printf("Ошибка при получении данных пользователя: %v", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка при получении данных пользователя"})
		return
	}

	claims := &Claims{
		Username: req.Username,
		UserID:   userId,
		Email:    Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		log.Printf("Ошибка при генерации токена: %v", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка при генерации токена"})
		return
	}

	context.SetCookie("jwt_token", tokenString, 3600*24, "/", "", true, true)

	context.JSON(http.StatusOK, gin.H{"message": "Успешная авторизация"})
}

func validatePassword(db *pgxpool.Pool, username string, password string) bool {
	var passwordHash string
	query := "SELECT password FROM events_service_data.users WHERE username = $1"

	err := db.QueryRow(context.Background(), query, username).Scan(&passwordHash)
	if err != nil {
		log.Printf("Ошибка при извлечении пароля из базы данных: %v", err)
		return false
	}

	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	if err != nil {
		return false
	}
	return true
}

func getUsersEmail(db *pgxpool.Pool, username string) (string, string, error) {
	var Email string
	var userId string
	query := "SELECT email, id FROM events_service_data.users WHERE username = $1"

	err := db.QueryRow(context.Background(), query, username).Scan(&Email, &userId)
	if err != nil {
		log.Printf("Ошибка при извлечении почты из базы данных: %v", err)
		return "", "", err
	}

	return Email, userId, nil
}
