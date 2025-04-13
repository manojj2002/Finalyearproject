package routes

import (
	"CVC_ragh/controllers"
	"CVC_ragh/utils"

	"github.com/gin-gonic/gin"
)

func RegisterScanningRoutes(router *gin.Engine) {
	scanGroup := router.Group("/api/static-scan")
	scanGroup.Use((utils.AuthMiddleware()))
	{
		scanGroup.POST("pull-image/*imageName", controllers.PullImageAndSave)
		scanGroup.POST("/scan-image/*imageName", controllers.ScanImageOnly)
		scanGroup.GET("/getImageDetails", controllers.GetUserImages)
		scanGroup.GET("/getScanResults", controllers.GetUserScanResults)
	}
}
