package controller

import (
	"closer-api-go/closerjwt"
	"closer-api-go/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type Controller struct{}

type LoginRequest struct {
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}
type RegisterRequest struct {
	Name                 string `json:"name"`
	PhoneNumber          int    `json:"phone_number"`
	CountryCode          string `json:"country_code"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"password_confirmation"`
}

func Login(c *gin.Context) {
	var loginRequest LoginRequest
	if err := c.BindJSON(&loginRequest); err != nil {
		badRequest(c, ErrorMessage{})
		return
	}
	user, err := service.LoginService(loginRequest.PhoneNumber, loginRequest.Password)
	if err != nil {
		ErrorResponse(c, ErrorObject{
			err.Error(),
			403,
		})
		return
	}
	jwtToken, err := closerjwt.CreateJwt(user)
	if err != nil {
		ErrorResponse(c, ErrorObject{})
		return
	}
	response := map[string]any{
		"user":         user,
		"access_token": jwtToken,
	}
	c.JSON(http.StatusOK, response)
}

func Register(c *gin.Context) {
	var registerRequest RegisterRequest
	if err := c.BindJSON(&registerRequest); err != nil {
		fmt.Println(err)
		badRequest(c, ErrorMessage{})
		return
	}
	if registerRequest.PasswordConfirmation != registerRequest.Password {
		badRequest(c, ErrorMessage{"Password confirmation does not match."})
		return
	}

	phoneNumberFull := registerRequest.CountryCode + strconv.Itoa(registerRequest.PhoneNumber)
	user, err := service.RegisterUserService(registerRequest.Name, phoneNumberFull, registerRequest.Password)
	if err != nil {
		badRequest(c, ErrorMessage{
			"Phone number already exists",
		})
		return
	}
	c.JSON(http.StatusOK, user)
	return
}
