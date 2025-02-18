package authorization

import (
	"context"
	"github.com/gin-gonic/gin"
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

func loginRequest(context *gin.Context, db *pgxpool.Pool) {
	var req UserData

	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Неверный формат запроса"})
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
