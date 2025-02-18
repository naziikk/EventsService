package user

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"net/http"
)

type Event struct {
	EventID string `json:"event_id" binding:"required"`
}

func CancelVisitRequest(context *gin.Context, db *pgxpool.Pool) {
	userId := context.Param("id")

	var req Event
	if err := context.ShouldBindJSON(&req); err != nil {
		log.Printf("Ошибка при парсинге JSON: %v", err)
		context.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат запроса"})
		return
	}

	if !UpdateDatabase(db, userId, req.EventID) {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при отмене посещения"})
	}
}

func UpdateDatabase(db *pgxpool.Pool, userId string, eventID string) bool {
	tx, err := db.BeginTx(context.Background(), pgx.TxOptions{})
	defer tx.Rollback(context.Background())

	eventPrice, err := GetEventPriceToRefund(db, eventID)
	if err != nil {
		log.Printf("Ошибка при получении цены события: %v", err)
		return false
	}

	query := "UPDATE events_service_data.users SET budget = budget + $1 WHERE id = $2"
	_, err = db.Exec(context.Background(), query, eventPrice, userId)
	if err != nil {
		log.Printf("Ошибка при обновлении бюджета пользователя: %v", err)
		return false
	}

	query = "UPDATE events_service_data.events SET places_count = places_count + 1 WHERE id = $1"
	_, err = db.Exec(context.Background(), query, eventID)
	if err != nil {
		log.Printf("Ошибка при обновлении количества свободных мест: %v", err)
		return false
	}

	query = "DELETE FROM events_service_data.tickets WHERE user_id = $1 AND event_id = $2"
	_, err = db.Exec(context.Background(), query, userId, eventID)
	if err != nil {
		log.Printf("Ошибка при удалении билета: %v", err)
		return false
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return false
	}

	return true
}

func GetEventPriceToRefund(db *pgxpool.Pool, eventID string) (int64, error) {
	query := "SELECT price FROM events_service_data.events WHERE id = $1"

	var price int64
	err := db.QueryRow(context.Background(), query, eventID).Scan(&price)
	if err != nil {
		return 0, err
	}

	return price, nil
}
