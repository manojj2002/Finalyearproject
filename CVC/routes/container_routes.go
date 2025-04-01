package routes

import (
	"CVC_ragh/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterContainerRoutes(router *gin.Engine) {
	containerGroup := router.Group("/api/container")
	{
		containerGroup.POST("/createContainer/:imageName")
		containerGroup.POST("/startContainer/:id", controllers.StartContainer)     // Start a new container
		containerGroup.POST("/stopContainer/:id", controllers.StopContainer)       // Stop a container
		containerGroup.GET("/inspectContainer/:id", controllers.InspectContainer)  // Inspect a container
		containerGroup.DELETE("/deleteContainer/:id", controllers.DeleteContainer) // Delete a container
		containerGroup.GET("/getContainerDetails/:id", controllers.GetContainerDetails)
	}
}
