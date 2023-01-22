package model

type Message struct {
	Id                int         `json:"id"`
	SenderId          int         `json:"sender_user_id"`
	ChatId            int         `json:"chat_id"`
	Message           string      `json:"message"`
	MessageType       MessageType `json:"message_type"`
	S3Path            string      `json:"s3_path"`
	Base64EncodedBlur string      `json:"base64_encoded_blur"`
	CreatedAt         string      `json:"created_at"`
}

func (m *Message) ToOutput() Message {
	return Message{
		Id:                m.Id,
		SenderId:          m.SenderId,
		ChatId:            m.ChatId,
		Message:           m.Message,
		MessageType:       m.MessageType,
		Base64EncodedBlur: m.Base64EncodedBlur,
		CreatedAt:         m.CreatedAt,
	}
}

type MessageType int

const (
	MessageTypeText  MessageType = 1
	MessageTypeImage MessageType = 2
)
