package model

type Chat struct {
	Id          int       `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Users       []User    `json:"users"`
	Messages    []Message `json:"messages"`
}

func (c *Chat) SetNewUser(user User) {
	c.Users = append(c.Users, user)
}
func (c *Chat) SetMessages(messages []Message) {
	c.Messages = messages
}
