package main

import (
	"context"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"time"
	"crypto/sha256"
	"encoding/hex"

	
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/gin-gonic/gin"
)

// ScanResult represents the MongoDB document
type ScanResult struct {
	ImageName string `bson:"imageName"`
	Report    string `bson:"report"`
	Hash      string `bson:"hash"`
	Timestamp string `bson:"timestamp"`
}

// Initialize MongoDB client
func initMongoDB() (*mongo.Client, *mongo.Collection, error) {
	clientOptions := options.Client().ApplyURI("mongodb+srv://manojjanasale2002:8LGLIukoAYINte1H@cluster0.jrm7r.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0") // Replace with your MongoDB URI
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, nil, err
	}

	// Ping MongoDB to verify connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, nil, err
	}

	log.Println("Connected to MongoDB!")
	collection := client.Database("Sample").Collection("Scan_results") // Database and Collection
	return client, collection, nil
}

func main() {
	// Initialize Gin router and MongoDB client
	r := gin.Default()
	client, collection, err := initMongoDB()
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(context.TODO())

	// Define the scan endpoint
	r.POST("/scan", func(c *gin.Context) {
		var request struct {
			ImageName string `json:"imageName"`
		}
		if err := c.ShouldBindJSON(&request); err != nil || request.ImageName == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Container image name is required"})
			return
		}

		// Run Trivy scan
		cmd := exec.Command("docker", "run", "--rm", "aquasec/trivy:latest", "image", request.ImageName)
		output, err := cmd.CombinedOutput()
		if err != nil {
			log.Printf("Error running Trivy: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan image", "details": string(output)})
			return
		}


		// Extract only the vulnerability report
		filteredOutput := strings.Split(string(output), "\n\n")
		var report string
		if len(filteredOutput) > 1 {
			report := filteredOutput[len(filteredOutput)-1] // Assume the last part contains the report
			log.Printf("Filtered Report: %s", report)
		} else {
			report := string(output)
			log.Printf("Filtered Report: %s", report)
		}
		
		// Save scan result in MongoDB
		scanResult := ScanResult{
			ImageName: request.ImageName,
			Report:    strings.TrimSpace(report),
			Hash:      generateSHA256Hash(report),
			Timestamp: time.Now().Format(time.RFC3339),
		}
		_, err = collection.InsertOne(context.TODO(), scanResult)
		if err != nil {
			log.Printf("Error inserting to MongoDB: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save scan result"})
			return
		}

		// Respond to the client
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"report":  scanResult.Report,
			"hash":    scanResult.Hash,
		})
	})

	// Start the server
	log.Println("Server running on port 8080")
	r.Run(":8080")
}

// Generate SHA256 hash for a string
func generateSHA256Hash(data string) string {
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}
