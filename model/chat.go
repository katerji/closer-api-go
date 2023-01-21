package model

type Chat struct {
	Id          int    `json:"id"`
	Name        int    `json:"name"`
	Description int    `json:"description"`
	Users       []User `json:"users"`
}

func (c *Chat) SetNewUser(user User) {
	c.Users = append(c.Users, user)
}