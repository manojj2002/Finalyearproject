package main

import (
	"log"
	"net/http"
	"os/exec"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize the Gin router
	r := gin.Default()

	// Define a default GET endpoint for the root
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to the Container Scanner API!",
		})
	})

	// Define the endpoint
	r.POST("/scan", func(c *gin.Context) {
		// Parse the request body
		var request struct {
			ImageName string `json:"imageName"`
		}
		if err := c.ShouldBindJSON(&request); err != nil || request.ImageName == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Container image name is required"})
			return
		}

		// Command to run Trivy scanner
		cmd := exec.Command("docker", "run", "--rm", "aquasec/trivy:latest", "image", request.ImageName)

		// Execute the command and capture the output
		output, err := cmd.CombinedOutput()
		if err != nil {
			log.Printf("Error running Trivy: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan image", "details": string(output)})
			return
		}

		// Send the Trivy output as response
		c.JSON(http.StatusOK, gin.H{"success": true, "report": strings.TrimSpace(string(output))})
	})

	// Start the server
	log.Println("Server running on port 8080")
	r.Run(":8080")
}

