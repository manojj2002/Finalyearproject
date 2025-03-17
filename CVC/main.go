package main

import (
	"CVC/config"
	"CVC/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Register user routes
	config.InitDB()
	routes.RegisterUserRoutes(r)

	r.Run(":5000") // Start server on port 5000
}
