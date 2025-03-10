package events

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"net/http"
)

type EventCreationData struct {
	EventName        string `json:"event_name" binding:"required"`
	EventDescription string `json:"event_description" binding:"required"`
	PlacesCount      int    `json:"places_count" binding:"required"`
	Price            int64  `json:"price" binding:"required"`
	Venue            string `json:"venue" binding:"required"`
}

func CreateEventRequest(ctx *gin.Context, db *pgxpool.Pool) {
	organizerId := ctx.GetString("userId")

	var req PrivateEvent
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"message": "Неверный формат запроса"})
		return
	}

	if !createEvent(db, req, organizerId) {
		log.Printf("Ошибка при создании события")
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка при создании события"})
	}
	ctx.JSON(http.StatusCreated, gin.H{"message": "Событие успешно создано"})
}

func createEvent(db *pgxpool.Pool, e PrivateEvent, organizerId string) bool {
	query := "INSERT INTO events_service_data.events (event_name, event_description, places_count, " +
		"price, organizer_id, venue, is_private) VALUES ($1, $2, $3, $4, $5, $6, $7)"

	_, err := db.Exec(context.Background(), query, e.EventName,
		e.EventDescription, e.PlacesCount, e.Price, organizerId, e.Venue, false)
	if err != nil {
		return false
	}
	return true
}
