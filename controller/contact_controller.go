package controller

import (
	"closer-api-go/model"
	"closer-api-go/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

const GetContactsRoute = "/contacts"

func GetContactsController(c *gin.Context) {
	user := GetCurrentUser(c)
	contacts, err := service.GetContacts(user.Id)
	if err != nil {
		ErrorResponse(c, ErrorObject{})
		return
	}
	result := make(map[string][]model.User)
	if contacts == nil {
		contacts = []model.User{}
	}
	result["contacts"] = contacts
	c.JSON(http.StatusOK, result)
	return
}
