package user

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"net/http"
)

func GetUserInfoRequest(context *gin.Context, db *pgxpool.Pool) {
	userID, exists := context.Get("userID")
	if exists == false {
		log.Printf("Ошибка при получении ID пользователя")
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Вы не зарегистрированы"})
		return
	}
	username, _ := context.Get("username")
	email, _ := context.Get("email")

	visitedEvents, err := GetVisitedEventsCount(db, userID.(string))
	if err != nil {
		log.Printf("Ошибка при получении количества посещенных событий: %v", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка при получении данных"})
		return
	}
	log.Printf("Данные отправлены пользователю %s", username)
	context.JSON(http.StatusOK, gin.H{"user_id": userID,
		"username":       username,
		"email":          email,
		"visited_events": visitedEvents})
}

func GetVisitedEventsCount(db *pgxpool.Pool, userID string) (int, error) {
	query := "SELECT COUNT(*) FROM events_service.tickets WHERE user_id = $1"
	var count int
	err := db.QueryRow(context.Background(), query, userID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
