package service

import (
	"closer-api-go/dbclient"
	"closer-api-go/model"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUserService(name string, phoneNumber string, password string) (model.User, error) {
	password, err := HashPassword(password)
	if err != nil {
		return model.User{}, err
	}
	db := dbclient.GetDbInstance()

	args := []any{
		name,
		phoneNumber,
		password,
	}
	userId := db.Insert(registerQuery, args...)
	if userId == 0 {
		return model.User{}, errors.New("phone Number already exists")
	}
	return model.User{
		Id:          userId,
		Name:        name,
		PhoneNumber: phoneNumber,
	}, nil
}

func LoginService(phoneNumber string, password string) (model.User, error) {
	var loginRow loginRow
	db := dbclient.GetDbInstance()
	row := db.QueryRow(loginQuery, phoneNumber)
	err := row.Scan(&loginRow.Id, &loginRow.PhoneNumber, &loginRow.Name, &loginRow.Password)
	if err != nil {
		fmt.Println(err)
		return model.User{}, errors.New("Phone number does not exist")
	}

	authed := CheckPasswordHash(password, loginRow.Password)
	if !authed {

		return model.User{}, errors.New("Incorrect password")
	}

	return model.User{
		Id:          loginRow.Id,
		Name:        loginRow.Name,
		PhoneNumber: loginRow.PhoneNumber,
	}, nil
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

const registerQuery = "INSERT INTO users_go (name, phone_number, password) VALUES (?, ?, ?)"
const loginQuery = "SELECT id, phone_number, name, password FROM users_go WHERE phone_number = ?"

type loginRow struct {
	Id          int    `json:"id"`
	PhoneNumber string `json:"phone_number"`
	Name        string `json:"name"`
	Password    string `json:"password"`
}
