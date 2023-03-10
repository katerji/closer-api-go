package service

import (
	"closer-api-go/dbclient"
	"closer-api-go/model"
	"database/sql"
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
		var base64EncodedString sql.NullString
		err = rows.Scan(&message.Id, &message.Message, &message.MessageType, &message.SenderId, &base64EncodedString, &message.CreatedAt)
		if err != nil {
			fmt.Println(err)
			return []model.Message{}
		}
		message.ChatId = chatId
		message.Base64EncodedBlur = base64EncodedString.String
		messages = append(messages, message.ToOutput())
	}
	return messages
}

func GetMessageById(messageId int) (model.Message, error) {
	var message model.Message
	err := dbclient.GetDbInstance().QueryRow(getMessageByIdQuery, messageId).Scan(
		&message.Message,
		&message.MessageType,
		&message.SenderId,
		&message.ChatId,
		&message.S3Path,
		&message.CreatedAt,
	)
	if err != nil {
		fmt.Println(err)
		return model.Message{}, err
	}
	return message, nil
}

const insertMessageQuery = "insert into messages_go (sender_user_id, chat_id, message, message_type) values (?, ?, ?, ?)"
const insertImageMessageQuery = "insert into messages_go (sender_user_id, chat_id, message, message_type, s3_path, blurred_image_base64) values (?, ?, ?, ?, ?, ?)"
const getChatMessagesQuery = "SELECT id, message, message_type, sender_user_id, blurred_image_base64, created_at FROM messages_go WHERE chat_id = ? ORDER BY created_at DESC LIMIT 50"
const getMessageByIdQuery = "SELECT message, message_type, sender_user_id, chat_id, s3_path, created_at FROM messages_go WHERE id = ?"
