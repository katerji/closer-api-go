package service

import (
	"closer-api-go/dbclient"
	"closer-api-go/model"
	"errors"
)

func AddContact(user model.User, contact model.User) error{
	insertId := dbclient.GetDbInstance().Insert(addContactQuery, user.Id, contact.Id, contact.Name)
	if insertId == 0 {
		return errors.New("already contacts")
	}
	insertId = dbclient.GetDbInstance().Insert(addContactQuery, contact.Id, user.Id, user.Name)
	if insertId == 0 {
		return errors.New("already contacts")
	}
	return nil
}


const addContactQuery = "INSERT INTO contacts_go (user_id, contact_user_id, contact_name) VALUES (?, ?, ?)"