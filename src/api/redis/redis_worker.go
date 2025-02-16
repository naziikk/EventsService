package redis_worker

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetWaitingList(eventId string) string {
	return "waiting_list:" + eventId
}

type WaitingListRequest struct {
	eventId string `json:"event_id" binding:"required"`
}

func AddUserToWaitingListRequest(context *gin.Context) {
	userId := context.Param("id")

	var req WaitingListRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат запроса"})
		return
	}

	key := GetWaitingList(req.eventId)

	if err := Rdb.LPush(Ctx, key, userId).Err(); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при добавлении пользователя в очередь"})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message":  "Пользователь добавлен в очередь",
		"event_id": req.eventId,
		"user_id":  userId,
	})
}

type NextUserRequest struct {
	eventId string `json:"event_id" binding:"required"`
}

func ProcessNextUserRequest(context *gin.Context) {
	var req NextUserRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат запроса"})
		return
	}
	key := GetWaitingList(req.eventId)

	userId, err := Rdb.RPop(Ctx, key).Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			context.JSON(http.StatusOK, gin.H{"message": "Очередь пуста"})
		} else {
			context.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при извлечении пользователя"})
		}
		return
	}

	context.JSON(http.StatusOK, gin.H{"user_id": userId})
}
