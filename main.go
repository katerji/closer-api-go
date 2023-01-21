package main

import (
	"closer-api-go/controller"
	"closer-api-go/middleware"
	"github.com/gin-gonic/gin"
)

func main() {
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
	//	Route::post('/message', [MessageController::class, 'create']);
	//	Route::post('/upload', [MessageController::class, 'upload']);
	//	Route::get('/messages/chat/{id}', [MessageController::class, 'index']);
	//
	//	Route::post('/chat', [ChatController::class, 'create']);
	//	Route::get('/chat/{id}', [ChatController::class, 'getChat']);
	//	Route::get('/chats', [ChatController::class, 'index']);

	r.Run(":9999")
}
