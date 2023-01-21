package model

type Chat struct {
	Id          int       `json:"id"`
	Name        int       `json:"name"`
	Description int       `json:"description"`
	Users       []User    `json:"users"`
	Messages    []Message `json:"messages"`
}

func (c *Chat) SetNewUser(user User) {
	c.Users = append(c.Users, user)
}
func (c *Chat) SetMessages(messages []Message) {
	c.Messages = messages
}
