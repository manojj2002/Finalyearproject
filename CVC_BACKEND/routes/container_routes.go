package routes

import (
	"CVC_ragh/controllers"
	"CVC_ragh/utils"

	"github.com/gin-gonic/gin"
)

func RegisterContainerRoutes(router *gin.Engine) {
	containerGroup := router.Group("/api/container")
	containerGroup.Use((utils.AuthMiddleware()))
	{
		containerGroup.POST("/createContainer/*imageName", controllers.CreateContainerFromImage)
		containerGroup.POST("/startContainer/:name", controllers.StartContainer)     // Start a new container
		containerGroup.POST("/stopContainer/:name", controllers.StopContainer)       // Stop a container
		containerGroup.DELETE("/deleteContainer/:name", controllers.DeleteContainer) // Delete a container
		containerGroup.GET("/getContainerDetails", controllers.GetUserContainers)

	}
}
