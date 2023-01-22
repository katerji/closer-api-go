package controller

import (
	"closer-api-go/awsclient"
	"closer-api-go/model"
	"closer-api-go/service"
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nfnt/resize"
	"image"
	"image/jpeg"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
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
	user := GetCurrentUser(c)
	if !service.IsUserInChat(chatId, user.Id) {
		SendUnauthorizedResponse(c)
		return
	}
	messages := service.GetChatMessages(chatId)
	response := make(map[string][]model.Message)
	response["messages"] = messages
	c.JSON(http.StatusOK, response)
}

const UploadImageRoute = "/message/upload"

func UploadImageController(c *gin.Context) {
	formFile, err := c.FormFile("file")

	if err != nil {
		fmt.Println(err)
		return
	}
	if !isFileValidated(formFile) {
		SendBadRequestResponse(c, ErrorMessage{"Only jpg/jpeg/png files allowed with a maximum size of 15mb"})
		return
	}
	user := GetCurrentUser(c)
	message := c.PostForm("message")
	chatId, err := strconv.Atoi(c.PostForm("chat_id"))
	if err != nil {
		SendBadRequestResponse(c, ErrorMessage{"No chat found"})
		return
	}
	if !service.IsUserInChat(chatId, user.Id) {
		SendUnauthorizedResponse(c)
		return
	}

	timeNow := time.Now().Unix()
	fileExtension := filepath.Ext(formFile.Filename)
	fileName := strconv.FormatInt(timeNow, 10) + fileExtension
	filePathOnS3 := fmt.Sprintf("%d/%d/%s", chatId, user.Id, fileName)

	imagePath, _ := formFile.Open()
	defer func(imagePath multipart.File) {
		err = imagePath.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(imagePath)

	blurredImageBase64 := blurAndResizeImage(imagePath)
	if len(blurredImageBase64) == 0 {
		SendBadRequestResponse(c, ErrorMessage{})
		return
	}
	go awsclient.UploadToS3(imagePath, filePathOnS3)
	messageObject := model.Message{
		SenderId:    user.Id,
		ChatId:      chatId,
		Message:     message,
		MessageType: model.MessageTypeImage,
		S3Path:      filePathOnS3,
		Base64EncodedBlur: blurredImageBase64,
	}
	service.InsertMessageImage(messageObject)
	SendEmptyOkayResponse(c)
	return
}

func isFileValidated(file *multipart.FileHeader) bool {
	ext := filepath.Ext(file.Filename)
	size := float64(file.Size) / (1 << 20)
	validExtension := ext == ".jpg" || ext == ".jpeg" || ext == ".png"
	return validExtension && size < 15
}

func blurAndResizeImage(file multipart.File) string {
	imageToResize, _, err := image.Decode(file)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	resizedImage := resize.Resize(9, 0, imageToResize, resize.Lanczos3)
	tempFile, _ := os.CreateTemp("", "")
	defer func(name string) {
		err = os.Remove(name)
		if err != nil {
			fmt.Println(err)
		}
	}(tempFile.Name())
	defer tempFile.Close()
	err = jpeg.Encode(tempFile, resizedImage, nil)
	readFile, err := os.ReadFile(tempFile.Name())
	if err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(readFile)
}
