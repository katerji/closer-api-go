package service

import (
	"closer-api-go/dbclient"
	"closer-api-go/model"
	"fmt"
)

func InsertMessage(message model.Message) int {
	return dbclient.GetDbInstance().Insert(insertMessageQuery, message.SenderId, message.ChatId, message.Message, model.MessageTypeText)
}

func InsertMessageImage(message model.Message) int {
	return dbclient.GetDbInstance().Insert(insertImageMessageQuery, message.SenderId, message.ChatId, message.Message, model.MessageTypeImage, message.S3Path, message.Base64EncodedBlur)
}

func GetChatMessages(chatId int) []model.Message {
	rows, err := dbclient.GetDbInstance().Query(getChatMessagesQuery, chatId)
	if err != nil {
		fmt.Println(err)
		return []model.Message{}
	}
	var messages []model.Message
	for rows.Next() {
		var message model.Message
		err = rows.Scan(&message.Id, &message.Message, &message.MessageType, &message.SenderId)
		if err != nil {
			fmt.Println(err)
			return []model.Message{}
		}
		message.ChatId = chatId
		messages = append(messages, message)
	}
	return messages
}

const insertMessageQuery = "insert into messages_go (sender_user_id, chat_id, message, message_type) values (?, ?, ?, ?)"
const insertImageMessageQuery = "insert into messages_go (sender_user_id, chat_id, message, message_type, s3_path, blurred_image_base64) values (?, ?, ?, ?, ?, ?)"
const getChatMessagesQuery = "SELECT id, message, message_type, sender_user_id FROM messages_go WHERE chat_id = ? ORDER BY created_at DESC LIMIT 50"
