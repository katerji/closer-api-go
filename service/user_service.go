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

func GetUserById(userId int) (model.User, error) {
	var user model.User
	err := dbclient.GetDbInstance().QueryRow(getUserById, userId).Scan(&user.Id, &user.Name, &user.PhoneNumber)
	if err != nil {
		fmt.Println(err)
		return model.User{}, err
	}
	return user, nil
}

const userByPhoneNumberQuery = "SELECT id, name, phone_number FROM users_go WHERE phone_number = ?"
const getUserById = "SELECT id, name, phone_number FROM users_go WHERE id = ?"
