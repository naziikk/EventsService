package authorization

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

type NewUserData struct {
	Password string `json:"password" binding:"required"`
}

func ResetPasswordRequest(context *gin.Context, db *pgxpool.Pool) {
	userID, exists := context.Get("userID")
	if exists == false {
		log.Printf("Ошибка при получении ID пользователя")
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Вы не зарегистрированы"})
		return
	}

	var req NewUserData
	if err := context.ShouldBindJSON(&req); err != nil {
		log.Printf("Неверные данные в json %v", err)
		context.JSON(http.StatusBadRequest, gin.H{"message": "Неверный формат запроса"})
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Ошибка при хешировании пароля: %v", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка при хешировании пароля"})
		return
	}

	if !updatePasswordInDB(db, userID.(int), string(passwordHash)) {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка при обновлении пароля"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Пароль успешно обновлен"})
}

func updatePasswordInDB(db *pgxpool.Pool, userID int, passwordHash string) bool {
	query := "UPDATE events_service.users SET password = $1 WHERE user_id = $2"

	_, err := db.Exec(context.Background(), query, passwordHash, userID)
	if err != nil {
		log.Printf("Ошибка при обновлении пароля в базе: %v", err)
		return false
	}
	return true
}
