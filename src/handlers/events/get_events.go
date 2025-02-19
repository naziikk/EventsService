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

type EventInfo struct {
	ID               int     `json:"id"`
	EventName        string  `json:"event_name"`
	EventDescription string  `json:"event_description"`
	PlacesCount      int     `json:"places_count"`
	Price            float64 `json:"price"`
	OrganizerID      int     `json:"organizer_id"`
	Venue            string  `json:"venue"`
	IsPrivate        bool    `json:"is_private"`
}

func GetUserEventsRequest(context *gin.Context, db *pgxpool.Pool) {
	userId := context.Param("id")
	rows, err := getEventsFromDB(db, userId)
	if err != nil {
		log.Println("Ошибка при запросе:", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "database error"})
		return
	}
	defer rows.Close()

	events, err := convertRowsToJSON(rows)
	if err != nil {
		log.Println("Ошибка при преобразовании записей в JSON", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "data processing error"})
		return
	}

	context.JSON(http.StatusOK, events)
}

func getEventsFromDB(db *pgxpool.Pool, userId string) (pgx.Rows, error) {
	query := "SELECT id, user_id, event FROM events_service_data.tickets WHERE user_id = $1"
	return db.Query(context.Background(), query, userId)
}

func convertRowsToJSON(rows pgx.Rows) ([]Ticket, error) {
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

func GetAllEventsRequest(context *gin.Context, db *pgxpool.Pool) {
	rows, err := getAllEventsFromDB(db)
	if err != nil {
		log.Printf("Ошибка при получении списка событий: %v", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка при получении списка событий"})
		return
	}
	defer rows.Close()

	events, err := convertEventsToJSON(rows)
	if err != nil {
		log.Printf("Ошибка при преобразовании данных в JSON: %v", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка при преобразовании данных в JSON"})
		return
	}
	context.JSON(http.StatusOK, events)
}

func getAllEventsFromDB(db *pgxpool.Pool) (pgx.Rows, error) {
	query := "SELECT * FROM events_service.events"
	return db.Query(context.Background(), query)
}

func convertEventsToJSON(rows pgx.Rows) ([]EventInfo, error) {
	var events []EventInfo
	for rows.Next() {
		var event EventInfo
		err := rows.Scan(&event.ID, &event.EventName, &event.EventDescription, &event.PlacesCount, &event.Price,
			&event.OrganizerID, &event.Venue, &event.IsPrivate)
		if err != nil {
			return nil, err
		}

		events = append(events, event)
	}
	return events, nil
}
