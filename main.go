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
	//r.Use(middleware.AuthMiddleware())
	//r.GET("/benchmark", MyBenchLogger(), benchEndpoint)
	//
	//// Authorization group
	//// authorized := r.Group("/", AuthRequired())
	//// exactly the same as:
	//authorized := r.Group("/")
	//// per group middleware! in this case we use the custom created
	//// AuthRequired() middleware just in the "authorized" group.
	//authorized.Use(AuthRequired())
	//{
	//	authorized.POST("/login", loginEndpoint)
	//	authorized.POST("/submit", submitEndpoint)
	//	authorized.POST("/read", readEndpoint)
	//
	//	// nested group
	//	testing := authorized.Group("testing")
	//	testing.GET("/analytics", analyticsEndpoint)
	//}
	api := r.Group("/api")
	auth := api.Group("/auth")
	auth.POST("/login", controller.Login)
	auth.POST("/register", controller.Register)


	api.Use(middleware.AuthMiddleware())

	api.GET("/chat", func (c *gin.Context) {
		fmt.Println(c.Get("user"))
		c.JSON(http.StatusOK, "yes")
	})
	r.Run(":9999")
}
