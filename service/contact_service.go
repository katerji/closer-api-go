package service

import (
	"closer-api-go/dbclient"
	"closer-api-go/model"
	"errors"
	"fmt"
)

func AddContact(user model.User, contact model.User) error {
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

func GetContacts(userId int) ([]model.User, error) {
	rows, err := dbclient.GetDbInstance().Query(getContactsQuery, userId)
	var contacts []model.User
	if err != nil {
		fmt.Println(err)
		return contacts, err
	}
	for rows.Next() {
		var contact model.User
		err = rows.Scan(&contact.Id, &contact.Name, &contact.PhoneNumber)
		if err != nil {
			return contacts, err
		}
		contacts = append(contacts, contact)
	}
	return contacts, nil
}

func AreUsersContacts(userId int, contactId int) bool {
	var queryResult int
	err := dbclient.GetDbInstance().QueryRow(areUserContactsQuery, userId, contactId).Scan(&queryResult)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return queryResult == 1
}

const addContactQuery = "INSERT INTO contacts_go (user_id, contact_user_id, contact_name) VALUES (?, ?, ?)"
const getContactsQuery = "SELECT u.id, u.name, u.phone_number FROM contacts_go c JOIN users_go u ON c.contact_user_id = u.id WHERE c.user_id = ?"
const areUserContactsQuery = "SELECT 1 FROM contacts_go WHERE user_id = ? AND contact_user_id = ?"
