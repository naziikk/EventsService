package user

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"net/http"
)

type AttendEventData struct {
	EventID string `json:"event_id" binding:"required"`
}

func AttendEventRequest(c *gin.Context, db *pgxpool.Pool) {
	userID := c.Param("id")

	var req AttendEventData
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Ошибка при парсинге JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат запроса"})
		return
	}

	ctx := c.Request.Context()

	eventPrice, err := GetEventPrice(ctx, db, req.EventID)
	if err != nil {
		log.Printf("Ошибка при получении цены события: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении цены события"})
		return
	}

	userBudget, err := GetUserBudget(ctx, db, userID)
	if err != nil {
		log.Printf("Ошибка при получении бюджета пользователя: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении бюджета пользователя"})
		return
	}

	if userBudget < eventPrice {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "Недостаточно средств на балансе",
			"budget":  userBudget,
			"price":   eventPrice,
		})
		return
	}

	freePlaces, err := CheckFreePlaces(ctx, db, req.EventID)
	if err != nil {
		log.Printf("Ошибка при проверке свободных мест: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при проверке свободных мест"})
		return
	}

	if !freePlaces {
		c.JSON(http.StatusForbidden, gin.H{"message": "Нет свободных мест, можете добавиться в лист ожидания или выбрать другое событие"})
		return
	}

	if err := UpdateDB(ctx, db, req.EventID, userID, eventPrice); err != nil {
		log.Printf("Ошибка при обновлении базы данных: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при обновлении базы данных"})
		return
	}

	log.Printf("User %s registered for event %s", userID, req.EventID)
	c.JSON(http.StatusOK, gin.H{"message": "Запись прошла успешно"})
}

func GetEventPrice(ctx context.Context, db *pgxpool.Pool, eventID string) (int64, error) {
	query := "SELECT price FROM events_service_data.events WHERE id = $1"

	var price int64
	err := db.QueryRow(ctx, query, eventID).Scan(&price)
	return price, err
}

func CheckFreePlaces(ctx context.Context, db *pgxpool.Pool, eventID string) (bool, error) {
	query := "SELECT places_count FROM events_service_data.events WHERE id = $1"

	var placesCount int
	err := db.QueryRow(ctx, query, eventID).Scan(&placesCount)
	return placesCount > 0, err
}

func GetUserBudget(ctx context.Context, db *pgxpool.Pool, userID string) (int64, error) {
	query := "SELECT budget FROM events_service_data.users WHERE id = $1"

	var budget int64
	err := db.QueryRow(ctx, query, userID).Scan(&budget)
	return budget, err
}

func UpdateDB(ctx context.Context, db *pgxpool.Pool, eventID, userID string, eventPrice int64) error {
	tx, err := db.BeginTx(context.Background(), pgx.TxOptions{})
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, "UPDATE events_service_data.events SET places_count = places_count - 1 WHERE id = $1", eventID)
	if err != nil {
		log.Printf("Ошибка обновления мест в событии: %v", err)
		return err
	}

	_, err = tx.Exec(ctx, "UPDATE events_service_data.users SET budget = budget - $1 WHERE id = $2", eventPrice, userID)
	if err != nil {
		log.Printf("Ошибка обновления бюджета пользователя: %v", err)
		return err
	}

	_, err = tx.Exec(ctx, "INSERT INTO events_service_data.tickets (user_id, event, price) VALUES ($1, $2, $3)", userID, eventID, eventPrice)
	if err != nil {
		log.Printf("Ошибка добавления билета в базу данных: %v", err)
	}

	return tx.Commit(ctx)
}
