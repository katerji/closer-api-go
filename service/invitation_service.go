package service

import (
	"closer-api-go/dbclient"
	"closer-api-go/model"
	"errors"
	"fmt"
)

func Invite(userId int, contactUserId int) error {
	invitationId := dbclient.GetDbInstance().Insert(insertInvitationQuery, userId, contactUserId)
	if invitationId == 0 {
		return errors.New("There is an existing pending invitation")
	}
	return nil
}

func GetSentInvitations(user model.User) ([]model.User, error) {
	invitations := make([]model.User, 0)
	rows, err := dbclient.GetDbInstance().Query(getSentInvitationsQuery, user.Id)
	if err != nil {
		fmt.Println(err)
		return invitations, errors.New("error fetching from db")
	}
	for rows.Next() {
		var contactUser model.User
		var invitation model.Invitation
		err := rows.Scan(&invitation.Id, &contactUser.Id, &contactUser.Name, &contactUser.PhoneNumber)
		if err != nil {
			var invitations []model.User
			return invitations, err
		}
		invitation.Inviter = user
		invitation.Contact = contactUser
		invitations = append(invitations, invitation.ToOutput())
	}
	return invitations, nil
}

func GetReceivedInvitations(user model.User) ([]model.User, error) {
	var invitations []model.User
	rows, err := dbclient.GetDbInstance().Query(getReceivedInvitationsQuery, user.Id)
	if err != nil {
		fmt.Println(err)
		return invitations, errors.New("error fetching from db")
	}
	for rows.Next() {
		var inviter model.User
		var invitation model.Invitation
		err := rows.Scan(&invitation.Id, &inviter.Id, &inviter.Name, &inviter.PhoneNumber)
		if err != nil {
			var invitations []model.User
			return invitations, err
		}
		invitation.Inviter = inviter
		invitation.Contact = user
		invitations = append(invitations, invitation.ToOutput())
	}
	return invitations, nil
}

func IsAuthorizedToAcceptOrRejectInvitation(userId int, inviterId int) bool {
	var isAuthorized int
	err := dbclient.GetDbInstance().QueryRow(isAuthorizedToAcceptOrRejectInvitationQuery, inviterId, userId).Scan(&isAuthorized)
	if err != nil {
		return false
	}
	return isAuthorized == 1
}

func IsAuthorizedToDeleteInvitation(userId int, invitationId int) bool {
	var isAuthorized int
	err := dbclient.GetDbInstance().QueryRow(isAuthorizedToDeleteInvitationQuery, invitationId, userId).Scan(&isAuthorized)
	if err != nil {
		return false
	}
	return isAuthorized == 1
}

func DeleteInvitation(invitationId int) bool {
	return dbclient.GetDbInstance().Exec(deleteInvitationQuery, invitationId)
}
func DeleteInvitationByUserIds(userId int, inviterId int) bool {
	return dbclient.GetDbInstance().Exec(deleteInvitationQuery, inviterId, userId)
}

const insertInvitationQuery = "INSERT INTO invitations_go (user_id, contact_user_id) VALUES (?, ?)"
const getSentInvitationsQuery = "SELECT i.id, i.contact_user_id, u.name as contact_user_name, u.phone_number as contact_phone_number " +
	"FROM invitations_go i " +
	"JOIN users_go u ON u.id = i.contact_user_id " +
	"WHERE i.user_id = ? " +
	"ORDER BY id desc"
const getReceivedInvitationsQuery = "SELECT i.id, i.user_id, u.name as inviter_name, u.phone_number as inviter_phone_number " +
	"FROM invitations_go i " +
	"JOIN users_go u ON u.id = i.user_id " +
	"WHERE i.contact_user_id = ? " +
	"ORDER BY id desc "
const getInviterFromInvitationIdQuery = "SELECT u.id, u.name, u.phone_number " +
	"FROM invitations_go i " +
	"JOIN users_go u ON i.user_id = u.id " +
	"WHERE i.id = ?"
const isAuthorizedToAcceptOrRejectInvitationQuery = "SELECT 1 FROM invitations_go WHERE user_id = ? AND contact_user_id = ?"
const isAuthorizedToDeleteInvitationQuery = "SELECT 1 FROM invitations_go WHERE id = ? AND user_id = ?"
const deleteInvitationQuery = "DELETE FROM invitations_go WHERE id = ?"
const deleteInvitationByUserIdsQuery = "DELETE FROM invitations_go WHERE user_id = ? AND contact_user_id = ?"
