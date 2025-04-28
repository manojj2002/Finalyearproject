package routes

import (
	"CVC/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(router *gin.Engine) {
	userGroup := router.Group("/api/user")
	{
		userGroup.POST("/register", controllers.RegisterUser)
	}
}
