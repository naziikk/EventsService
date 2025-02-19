package authorization

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

type UserData struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func LoginRequest(context *gin.Context, db *pgxpool.Pool) {
	var req UserData

	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Неверный формат запроса"})
		return
	}

	exists, err := checkUsernameUniqueness(db, req.Username)
	if err != nil {
		log.Printf("Ошибка при проверке уникальности имени пользователя: %v", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка при проверке уникальности имени пользователя"})
		return
	}
	if !exists {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Имя пользователя уже занято"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Ошибка при хешировании пароля: %v", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка при хешировании пароля"})
		return
	}

	if !savePasswordInDB(db, UserData{Username: req.Username, Email: req.Email, Password: string(hashedPassword)}) {
		log.Printf("Ошибка при сохранении пароля")
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка при сохранении пароля"})
		return
	}
	context.JSON(http.StatusCreated, gin.H{"message": "Пароль успешно сохранен"})
}

func savePasswordInDB(db *pgxpool.Pool, req UserData) bool {
	query := "INSERT INTO events_service_data.users (username, email, password) VALUES ($1, $2, $3)"
	_, err := db.Exec(context.Background(), query, req.Username, req.Email, req.Password)
	if err != nil {
		log.Printf("Ошибка при сохранении пароля в базе данных: %v", err)
		return false
	}
	return true
}

func checkUsernameUniqueness(db *pgxpool.Pool, username string) (bool, error) {
	query := "SELECT 1 FROM events_service_data.users WHERE username = $1 LIMIT 1"
	var exists int
	err := db.QueryRow(context.Background(), query, username).Scan(&exists)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return true, nil
		}
		return false, err
	}
	return false, nil
}
