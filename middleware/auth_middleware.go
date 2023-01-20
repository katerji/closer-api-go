package middleware

import (
	"closer-api-go/closerjwt"
	"closer-api-go/controller"
	"github.com/gin-gonic/gin"
	"strings"
)

func AuthMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {
		authorization := c.GetHeader("Authorization")
		token := strings.ReplaceAll(authorization, "Bearer ", "")
		userObject, err := closerjwt.VerifyToken(token)
		if err != nil {
			controller.ErrorResponse(c, controller.ErrorObject{
				Message: "Unauthorized",
				Code:    403,
			})
			return
		}
		c.Set("user", userObject)
		c.Next()
	}

}
