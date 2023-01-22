package controller

import (
	"closer-api-go/model"
	"closer-api-go/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

const InvitationGroupRoute = "/invitation"

const InviteRoute = "/send/:phone_number"

func InviteController(c *gin.Context) {
	contactPhoneNumber := c.Param("phone_number")
	user := GetCurrentUser(c)
	if contactPhoneNumber == user.PhoneNumber {
		SendBadRequestResponse(c, ErrorMessage{})
		return
	}
	contact, err := service.GetUserByPhoneNumber(contactPhoneNumber)
	if err != nil {
		SendBadRequestResponse(c, ErrorMessage{"Phone number does not exist"})
		return
	}

	err = service.Invite(user.Id, contact.Id)
	if err != nil {
		SendBadRequestResponse(c, ErrorMessage{"There is a pending invitation"})
		return
	}
	invitations, err := getSentAndReceivedInvitations(user)
	if err != nil {
		fmt.Println(err)
		ErrorResponse(c, ErrorObject{})
		return
	}
	c.JSON(http.StatusOK, invitations)
	return
}

const AcceptInvitationRoute = "/accept/:inviter_id"

func AcceptInvitationController(c *gin.Context) {
	inviterId, _ := strconv.Atoi(c.Param("inviter_id"))
	user := GetCurrentUser(c)
	if !service.IsAuthorizedToAcceptOrRejectInvitation(user.Id, inviterId) {
		SendUnauthorizedResponse(c)
		return
	}
	inviter, err := service.GetUserById(inviterId)
	if err != nil {
		ErrorResponse(c, ErrorObject{})
		return
	}
	err = service.AddContact(user, inviter)
	if err != nil {
		SendBadRequestResponse(c, ErrorMessage{"Already contacts"})
		return
	}
	service.DeleteInvitationByUserIds(user.Id, inviterId)
	invitations, err := getSentAndReceivedInvitations(user)
	if err != nil {
		ErrorResponse(c, ErrorObject{"Error fetching invitations", 500})
		return
	}
	c.JSON(http.StatusOK, invitations)
	return
}

const RejectInvitationRoute = "/reject/:inviter_id"

func RejectInvitationController(c *gin.Context) {
	inviterId, _ := strconv.Atoi(c.Param("inviter_id"))
	user := GetCurrentUser(c)
	if !service.IsAuthorizedToAcceptOrRejectInvitation(user.Id, inviterId) {
		SendUnauthorizedResponse(c)
		return
	}
	service.DeleteInvitationByUserIds(user.Id, inviterId)
	invitations, err := getSentAndReceivedInvitations(user)
	if err != nil {
		ErrorResponse(c, ErrorObject{"Error fetching invitations", 500})
		return
	}
	c.JSON(http.StatusOK, invitations)
	return
}

const DeleteInvitationRoute = "/delete/:invitation_id"

func DeleteInvitationController(c *gin.Context) {
	invitationId, err := strconv.Atoi(c.Param("invitation_id"))
	if err != nil {
		SendBadRequestResponse(c, ErrorMessage{})
	}
	user := GetCurrentUser(c)
	if !service.IsAuthorizedToDeleteInvitation(user.Id, invitationId) {
		SendUnauthorizedResponse(c)
		return
	}
	service.DeleteInvitation(invitationId)
	invitations, err := getSentAndReceivedInvitations(user)
	if err != nil {
		ErrorResponse(c, ErrorObject{"Error fetching invitations", 500})
		return
	}
	c.JSON(http.StatusOK, invitations)
	return
}

const GetInvitationsRoute = "/invitations"

func GetInvitationsController(c *gin.Context) {
	user := GetCurrentUser(c)
	invitations, err := getSentAndReceivedInvitations(user)
	if err != nil {
		fmt.Println(err)
		ErrorResponse(c, ErrorObject{})
		return
	}
	c.JSON(http.StatusOK, invitations)
	return
}

func getSentAndReceivedInvitations(user model.User) (map[string][]model.Invitation, error) {
	invitations := make(map[string][]model.Invitation)
	sentInvitations, err := service.GetSentInvitations(user)
	if err != nil {
		return invitations, err
	}
	if sentInvitations == nil {
		sentInvitations = []model.Invitation{}
	}
	receivedInvitations, err := service.GetReceivedInvitations(user)
	if err != nil {
		return invitations, err
	}
	if receivedInvitations == nil {
		receivedInvitations = []model.Invitation{}
	}
	invitations["sent_invitations"] = sentInvitations
	invitations["received_invitations"] = receivedInvitations
	return invitations, nil
}
