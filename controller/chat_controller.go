package controller

import (
	"closer-api-go/model"
	"closer-api-go/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

const GetChatsRoute = "/chats"
func GetChatsController(c *gin.Context) {
	user := GetCurrentUser(c)
	chats, err := service.GetUserChats(user.Id)
	if err != nil {
		ErrorResponse(c, ErrorObject{})
		return
	}
	response := make(map[string][]model.Chat)
	response["chats"] = chats
	c.JSON(http.StatusOK, response)
	return
}
