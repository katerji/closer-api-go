package controller

import (
	"closer-api-go/model"
	"closer-api-go/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
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

const CreateChatRoute = "/chat/:contact_id"

func CreateChatController(c *gin.Context) {
	user := GetCurrentUser(c)
	contactId, err := strconv.Atoi(c.Param("contact_id"))
	if err != nil {
		SendBadRequestResponse(c, ErrorMessage{})
		return
	}
	if !service.AreUsersContacts(user.Id, contactId) {
		SendUnauthorizedResponse(c)
		return
	}
	chatId := service.GetChatIdByUserIds(user.Id, contactId)
	response := make(map[string]model.Chat)
	if chatId > 0 {
		response["chat"], err = service.GetChatById(chatId, user.Id)
		if err != nil {
			ErrorResponse(c, ErrorObject{})
			return
		}
	} else {
		chatId = service.CreateChatWithoutAnyInfo(user.Id, contactId)
		response["chat"], err = service.GetChatById(chatId, user.Id)
	}
	c.JSON(http.StatusOK, response)
	return
}

const GetChatRoute = "/chat/:chat_id"

func GetChatController(c *gin.Context) {
	chatId, err := strconv.Atoi(c.Param("chat_id"))
	if err != nil {
		fmt.Println(chatId)
		fmt.Println(err)
		SendBadRequestResponse(c, ErrorMessage{})
		return
	}
	user := GetCurrentUser(c)
	if !service.IsUserInChat(chatId, user.Id) {
		SendUnauthorizedResponse(c)
		return
	}
	chat, err := service.GetChatById(chatId, user.Id)
	if err != nil {
		ErrorResponse(c, ErrorObject{})
		return
	}
	response := make(map[string]model.Chat)
	response["chat"] = chat
	c.JSON(http.StatusOK, response)
	return
}
