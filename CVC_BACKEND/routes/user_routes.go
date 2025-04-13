package routes

import (
	"CVC_ragh/controllers"
	"CVC_ragh/models"

	"github.com/gin-gonic/gin"
)

type ImageWithContainers struct {
	Image      models.Image       `json:"image"`
	Containers []models.Container `json:"containers"`
}

func RegisterUserRoutes(router *gin.Engine) {
	userGroup := router.Group("/api/user")
	{
		userGroup.POST("/register", controllers.RegisterUser)
		userGroup.POST("/login", controllers.LoginUser)
	}
}

func RegisterGitHubAuthRoutes(router *gin.Engine) {
	auth := router.Group("/api/auth")
	{
		auth.GET("/login", controllers.HandleGitHubLogin)
		auth.GET("/github/callback", controllers.HandleGitHubCallback)
	}
}
