package service

import (
	"closer-api-go/dbclient"
	"closer-api-go/model"
	"errors"
	"fmt"
)

func GetUserByPhoneNumber(phoneNumber string) (model.User, error) {
	db := dbclient.GetDbInstance()
	var user model.User
	err := db.QueryRow(userByPhoneNumberQuery, phoneNumber).Scan(&user.Id, &user.Name, &user.PhoneNumber)
	if err != nil {
		fmt.Println(err)
		return model.User{}, errors.New("error fetching from db")
	}
	return user, nil
}

const userByPhoneNumberQuery = "SELECT id, name, phone_number FROM users_go WHERE phone_number = ?"
