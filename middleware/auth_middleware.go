package middleware

import (
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	//requiredToken := os.Getenv("API_TOKEN")
	//
	//if requiredToken == "" {
	//	log.Fatal("Por favor, defina a variavel API_TOKEN")
	//}

	return func(c *gin.Context) {
		//token := c.Request.FormValue("api_token")
		//
		//if token == "" {
		//	c.JSON(http.StatusBadRequest, gin.H{"message": "Token deve ser preenchido"})
		//	return
		//}
		name := c.Query("name")
		if name != "a" {
			errorMessage := map[string]string{
				"error": "Unauthorized",
			}

			c.AbortWithStatusJSON(403, errorMessage)
			return
		}
		c.Next()
	}

}
