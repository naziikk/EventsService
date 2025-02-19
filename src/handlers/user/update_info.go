package user

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"net/http"
)

type UpdateUserData struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
}

func UpdateUserInfoRequest(c *gin.Context, db *pgxpool.Pool) {
	userID, exists := c.Get("userID")
	if exists == false {
		c.JSON(401, gin.H{"message": "Вы не зарегистрированы"})
		return
	}
	username, _ := c.Get("username")

	var req UpdateUserData
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"message": "Неверный формат запроса"})
		return
	}
	if username != req.Username {
		if !CheckIfUsernameExists(db, req.Username) {
			log.Printf("Пользователь %s хотел изменить имя на %s, но такое имя уже существует", username, req.Username)
			c.JSON(409, gin.H{"message": "Пользователь с таким именем уже существует"})
			return
		} else {
			query := "UPDATE events_service_data.users SET username = $1 WHERE user_id = $2"
			err := db.QueryRow(context.Background(), query, req.Username, userID).Scan()
			if err != nil {
				log.Printf("Ошибка при обновлении имени пользователя: %v", err)
				c.JSON(500, gin.H{"message": "Ошибка при обновлении имени пользователя"})
				return
			}
		}
	}
	c.JSON(http.StatusOK, gin.H{"message": "Информация о пользователе успешно обновлена"})
}

func CheckIfUsernameExists(db *pgxpool.Pool, username string) bool {
	query := "SELECT COUNT(*) FROM events_service_data.users WHERE username = $1"
	var count int
	err := db.QueryRow(context.Background(), query, username).Scan(&count)
	if err != nil {
		return false
	}
	if count > 0 {
		return false
	}

	return true
}
