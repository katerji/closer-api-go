package main

import (
	"closer-api-go/controller"
	"closer-api-go/middleware"
	"closer-api-go/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
	}
	r := gin.Default()

	api := r.Group("/api")
	auth := api.Group(controller.AuthGroupRoute)
	auth.POST(controller.LoginRoute, controller.Login)
	auth.POST(controller.RegisterRoute, controller.Register)

	api.Use(middleware.AuthMiddleware())

	api.GET(controller.GetInvitationsRoute, controller.GetInvitationsController)

	invitationGroup := api.Group(controller.InvitationGroupRoute)
	invitationGroup.POST(controller.InviteRoute, controller.InviteController)
	invitationGroup.POST(controller.AcceptInvitationRoute, controller.AcceptInvitationController)
	invitationGroup.DELETE(controller.RejectInvitationRoute, controller.RejectInvitationController)
	invitationGroup.DELETE(controller.DeleteInvitationRoute, controller.DeleteInvitationController)

	api.GET(controller.GetContactsRoute, controller.GetContactsController)

	api.GET(controller.GetChatsRoute, controller.GetChatsController)
	api.POST(controller.CreateChatRoute, controller.CreateChatController)
	api.GET(controller.GetChatRoute, controller.GetChatController)

	api.POST(controller.CreateMessageRoute, controller.CreateMessageController)
	api.GET(controller.GetChatMessagesRoute, controller.GetChatMessageController)
	api.POST(controller.UploadImageRoute, controller.UploadImageController)
	s := service.Server{}
	go s.InitGrpc()

	err = r.Run(":85")
	if err != nil {
		fmt.Println(err)
	}
}
