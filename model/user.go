package model

type User struct {
	Id    int    `json:"id"`
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
}
