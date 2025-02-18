package events

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"net/http"
)

type Ticket struct {
	ID     int    `json:"id"`
	UserID int    `json:"user_id"`
	Event  string `json:"event_id"`
	Price  string `json:"price"`
}

func GetUserEventsRequest(context *gin.Context, db *pgxpool.Pool) {
	userId := context.Param("id")
	rows, err := GetEventsFromDB(db, userId)
	if err != nil {
		log.Println("Ошибка при запросе:", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "database error"})
		return
	}
	defer rows.Close()

	events, err := ConvertRowsToJSON(rows)
	if err != nil {
		log.Println("Ошибка при преобразовании записей в JSON", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "data processing error"})
		return
	}

	context.JSON(http.StatusOK, events)
}

func GetEventsFromDB(db *pgxpool.Pool, userId string) (pgx.Rows, error) {
	query := "SELECT id, user_id, event FROM events_service_data.tickets WHERE user_id = $1"
	return db.Query(context.Background(), query, userId)
}

func ConvertRowsToJSON(rows pgx.Rows) ([]Ticket, error) {
	var events []Ticket

	for rows.Next() {
		var ticket Ticket
		err := rows.Scan(&ticket.ID, &ticket.UserID, &ticket.Event)

		if err != nil {
			return nil, err
		}

		events = append(events, ticket)
	}
	return events, nil
}
