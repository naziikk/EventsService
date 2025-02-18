package events

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
)

type EventData struct {
	eventName   string `json:"event_name" binding:"required"`
	placesCount int    `json:"places_count" binding:"required"`
	price       int64  `json:"price" binding:"required"`
}

func CreateEventRequest(context *gin.Context, db *pgxpool.Pool) {
	var req EventData
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат запроса"})
		return
	}

	if !CreateEvent(db, req) {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании события"})
		return
	}
	context.JSON(http.StatusCreated, gin.H{"message": "Событие успешно создано"})
}

func CreateEvent(db *pgxpool.Pool, req EventData) bool {
	ctx := context.Background()

	tx, err := db.Begin(ctx)
	if err != nil {
		return false
	}
	defer tx.Rollback(ctx)

	query := "INSERT INTO events_service_data.events (event_name, places_count, price) VALUES ($1, $2, $3)"
	_, err = tx.Exec(ctx, query, req.eventName, req.placesCount, req.price)
	if err != nil {
		return false
	}
	err = tx.Commit(ctx)
	if err != nil {
		return false
	}
	return true
}
