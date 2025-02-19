package events

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"net/http"
)

type PrivateEvent struct {
	EventName        string `json:"event_name" binding:"required"`
	EventDescription string `json:"event_description" binding:"required"`
	PlacesCount      int    `json:"places_count" binding:"required"`
	Price            int64  `json:"price" binding:"required"`
	Venue            string `json:"venue" binding:"required"`
	InvitationCode   string `json:"code" binding:"required"`
}

func CreatePrivateEventRequest(ctx *gin.Context, db *pgxpool.Pool) {
	organizerId := ctx.GetString("userId")

	var req PrivateEvent
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"message": "Неверный формат запроса"})
		return
	}
	if !checkCodeUniqueness(db, req.InvitationCode) {
		log.Printf("Пригласительный код уже занят, отказано в создании события")
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Пригласительный код уже занят, создайте другой"})
	}

	if !createPrivateEvent(db, req, organizerId) {
		log.Printf("Ошибка при создании события")
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка при создании события"})
	}
	ctx.JSON(http.StatusCreated, gin.H{"message": "Событие успешно создано"})
}

func createPrivateEvent(db *pgxpool.Pool, e PrivateEvent, organizerId string) bool {
	query := "INSERT INTO events_service_data.events (event_name, event_description, places_count, " +
		"price, organizer_id, venue, is_private, invitation_code) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)"

	_, err := db.Exec(context.Background(), query, e.EventName,
		e.EventDescription, e.PlacesCount, e.Price, organizerId, e.Venue, true, e.InvitationCode)
	if err != nil {
		return false
	}
	return true
}

func checkCodeUniqueness(db *pgxpool.Pool, code string) bool {
	query := "SELECT invitation_code FROM events_service_data.events WHERE invitation_code = $1"
	var invitationCode string
	err := db.QueryRow(context.Background(), query, code).Scan(&invitationCode)
	if err != nil {
		return true
	}
	return false
}
