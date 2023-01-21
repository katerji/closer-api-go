package model

type Message struct {
	Id          int         `json:"id"`
	SenderId    int         `json:"sender_user_id"`
	ChatId      int         `json:"chat_id"`
	Message     string      `json:"message"`
	MessageType MessageType `json:"message_type"`
	S3Path      string      `json:"s3_path"`
}

type MessageType int

const (
	MessageTypeText  MessageType = 1
	MessageTypeImage MessageType = 2
)
