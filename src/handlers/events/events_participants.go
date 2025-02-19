package events

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"net/http"
)

type EventID struct {
	EventID int `json:"event_id" binding:"required"`
}

func GetEventsParticipantsRequest(ctx *gin.Context, db *pgxpool.Pool) {
	var req EventID
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Неверный формат запроса"})
		return
	}

	participants, err := getEventsParticipants(db, req)
	if err != nil {
		log.Printf("Ошибка при получении участников события: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка при получении участников события"})
	}
	ctx.JSON(http.StatusOK, gin.H{"participants": participants})
}

func getEventsParticipants(db *pgxpool.Pool, req EventID) ([]string, error) {
	ctx := context.Background()

	queryUsers := "SELECT user_id FROM events_service.tickets WHERE event_id = $1"
	rows, err := db.Query(ctx, queryUsers, req.EventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	participants := make([]string, 0)

	queryUsername := "SELECT username FROM events_service_data.users WHERE id = $1"
	for rows.Next() {
		var userID int64
		if err := rows.Scan(&userID); err != nil {
			return nil, err
		}

		var username string
		if err := db.QueryRow(ctx, queryUsername, userID).Scan(&username); err != nil {
			return nil, err
		}
		participants = append(participants, username)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return participants, nil
}
