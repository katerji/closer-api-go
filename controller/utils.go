package controller

import (
	"closer-api-go/model"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ErrorMessage struct {
	string
}

type ErrorObject struct {
	Message string
	Code    int
}

func SendBadRequestResponse(c *gin.Context, errorMessage ErrorMessage) {
	errorReturn := map[string]string{
		"error": "Bad request.",
	}
	if len(errorMessage.string) > 0 {
		errorReturn["error"] = errorMessage.string
	}
	c.AbortWithStatusJSON(http.StatusBadRequest, errorReturn)
}
func ErrorResponse(c *gin.Context, errorObject ErrorObject) {
	errorReturn := map[string]string{
		"error": "Something went wrong.",
	}
	code := http.StatusInternalServerError
	if len(errorObject.Message) > 0 {
		errorReturn["error"] = errorObject.Message
		code = errorObject.Code
	}
	c.AbortWithStatusJSON(code, errorReturn)
}

func SendUnauthorizedResponse(c *gin.Context) {
	errorMessage := map[string]string{
		"error": "Unauthorized",
	}
	c.AbortWithStatusJSON(http.StatusUnauthorized, errorMessage)
}
func SendEmptyOkayResponse(c *gin.Context) {
	response := map[string]string{
		"success": "true",
	}
	c.JSON(http.StatusOK, response)
}

func GetCurrentUser(c *gin.Context) model.User {
	var user model.User
	userEnc, _ := c.Get("user")
	jsonEncodedUser, _ := json.Marshal(userEnc)
	err := json.Unmarshal(jsonEncodedUser, &user)
	if err != nil {
		fmt.Println(err)
		return user
	}
	return user
}
