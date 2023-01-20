package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type ErrorMessage struct {
	string
}

type ErrorObject struct {
	Message string
	Code int
}

func badRequest(c *gin.Context, errorMessage ErrorMessage) {
	errorReturn := map[string]string{
		"error": "Bad request.",
	}
	if len(errorMessage.string) > 0 {
		errorReturn["error"] = errorMessage.string
	}
	c.AbortWithStatusJSON(400, errorReturn)
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
