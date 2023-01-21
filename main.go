package main

import (
	"closer-api-go/controller"
	"closer-api-go/middleware"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	r := gin.Default()

	api := r.Group("/api")
	auth := api.Group("/auth")
	auth.POST("/login", controller.Login)
	auth.POST("/register", controller.Register)


	api.Use(middleware.AuthMiddleware())

	api.GET("/chat", func (c *gin.Context) {
		fmt.Println(c.Get("user"))
		c.JSON(http.StatusOK, "yes")
	})

	api.GET("/invitations", controller.GetInvitationsController)

	invitationGroup := api.Group("/invitation")
	invitationGroup.POST("/send/:phone_number", controller.InviteController)
	invitationGroup.POST("/accept/:invitation_id", controller.AcceptInvitationController)
	//invitationGroup.POST("/reject/:invitation_id", controller.InviteController)
	//invitationGroup.POST("/delete/:invitation_id", controller.InviteController)

	r.Run(":9999")
}
