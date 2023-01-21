package controller

import (
	"closer-api-go/model"
	"closer-api-go/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func InviteController(c *gin.Context) {
	contactPhoneNumber := c.Param("phone_number")
	contact, err := service.GetUserByPhoneNumber(contactPhoneNumber)
	if err != nil {
		badRequest(c, ErrorMessage{"Phone number does not exist"})
		return
	}
	user := GetCurrentUser(c)
	err = service.Invite(user.Id, contact.Id)
	if err != nil {
		badRequest(c, ErrorMessage{"There is a pending invitation"})
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
