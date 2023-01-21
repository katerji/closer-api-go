package controller

import (
	"closer-api-go/model"
	"closer-api-go/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type MessageRequest struct {
	ChatId      int               `json:"chat_id"`
	Message     string            `json:"message"`
	MessageType model.MessageType `json:"message_type"`
}

const CreateMessageRoute = "/message"
func CreateMessageController(c *gin.Context) {
	var messageRequest MessageRequest
	if err := c.BindJSON(&messageRequest); err != nil {
		fmt.Println(err)
		SendBadRequestResponse(c, ErrorMessage{})
		return
	}
	user := GetCurrentUser(c)
	if !service.IsUserInChat(messageRequest.ChatId, user.Id) {
		SendUnauthorizedResponse(c)
		return
	}
	message := model.Message{
		SenderId:    user.Id,
		ChatId:      messageRequest.ChatId,
		Message:     messageRequest.Message,
		MessageType: messageRequest.MessageType,
	}
	service.InsertMessage(message)
	SendEmptyOkayResponse(c)
	return
}

const GetChatMessagesRoute = "/messages/chat/:chat_id"
func GetChatMessageController(c *gin.Context) {
	chatId, err := strconv.Atoi(c.Param("chat_id"))
	if err != nil {
		SendBadRequestResponse(c, ErrorMessage{})
		return
	}
	user:=GetCurrentUser(c)
	if !service.IsUserInChat(chatId, user.Id) {
		SendUnauthorizedResponse(c)
		return
	}
	messages := service.GetChatMessages(chatId)
	response := make(map[string][]model.Message)
	response["messages"] = messages
	c.JSON(http.StatusOK, response)
}