package routes

import (
	"CVC_ragh/controllers"
	"CVC_ragh/utils"

	"github.com/gin-gonic/gin"
)

func RegisterDynamicScanningRoutes(router *gin.Engine) {
	dynamicScanGroup := router.Group("/api/dynamic-scan")
	dynamicScanGroup.Use((utils.AuthMiddleware()))
	{
		// Fetch Falco logs
		dynamicScanGroup.GET("/logs", controllers.StartFalcoInBackground)
		dynamicScanGroup.POST("/stop-falco", controllers.StopFalco)
		dynamicScanGroup.GET("getUserAlerts", controllers.GetFalcoAlerts)
	}
}

func RegisterFalcoWebHook(router *gin.Engine) {
	falcoWebHookGroup := router.Group("/api")
	// falcoWebHookGroup.Use((utils.AuthMiddleware()))
	{
		falcoWebHookGroup.POST("/falco-webhook", controllers.FalcoWebhookHandler)
		router.GET("/metrics", controllers.PrometheusMetricsHandler)
	}
}
