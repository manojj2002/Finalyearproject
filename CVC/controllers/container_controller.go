package controllers

import (
	"CVC_ragh/models"
	"context"
	"fmt"
	"net/http"
	"time"

	"CVC_ragh/config"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func CreateContainerFromImage(c *gin.Context) {
	imageName := c.Param("imageName")              // Get image name from request
	containerName := c.DefaultPostForm("name", "") // Get container name from request, default to empty string if not provided

	// Initialize Docker client
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Docker connection failed"})
		return
	}
	defer cli.Close()

	// Create a new container from the image (but do NOT start it)
	resp, err := cli.ContainerCreate(
		context.Background(),
		&container.Config{
			Image: imageName,                           // Set image name
			Cmd:   []string{"tail", "-f", "/dev/null"}, // Keep container idle
			// If a name is provided, use it, otherwise Docker will assign one
			Labels: map[string]string{
				"name": containerName,
			},
		},
		nil, nil, nil, containerName, // Pass the container name if provided
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create container"})
		return
	}

	// Store container details in the database (including the name)
	err = StoreContainerInDB(resp.ID, imageName, containerName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save container info"})
		return
	}

	// Respond with container ID and name (even if it's not started yet)
	c.JSON(http.StatusOK, gin.H{"message": "Container created successfully", "container_id": resp.ID, "container_name": containerName})
}

// StartContainer starts a new container and updates its status in the database
func StartContainer(c *gin.Context) {
	// Get container ID from URL params
	containerID := c.Param("id")

	// Initialize Docker client
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Docker client error"})
		return
	}
	defer cli.Close()

	ctx := context.Background()

	// Check if the container exists and is stopped
	containerJSON, err := cli.ContainerInspect(ctx, containerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to inspect container"})
		return
	}

	// If the container is already running, notify the user
	if containerJSON.State.Running {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Container is already running"})
		return
	}

	// Start the container
	if err := cli.ContainerStart(ctx, containerID, types.ContainerStartOptions{}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start container"})
		return
	}

	// Update the container status in the database to "running"
	err = UpdateContainerStatusInDB(containerID, "running")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update container status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Container started", "container_id": containerID})
}

// StopContainer stops a running container
func StopContainer(c *gin.Context) {
	containerID := c.Param("id")

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Docker client error"})
		return
	}
	defer cli.Close()

	ctx := context.Background()
	if err := cli.ContainerStop(ctx, containerID, container.StopOptions{}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to stop container"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Container stopped", "container_id": containerID})
}

// InspectContainer retrieves container details
func InspectContainer(c *gin.Context) {
	containerID := c.Param("id")

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Docker client error"})
		return
	}
	defer cli.Close()

	ctx := context.Background()
	containerJSON, err := cli.ContainerInspect(ctx, containerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to inspect container"})
		return
	}

	c.JSON(http.StatusOK, containerJSON)
}

// DeleteContainer removes a container and updates its status in the database
func DeleteContainer(c *gin.Context) {
	containerID := c.Param("id")

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Docker client error"})
		return
	}
	defer cli.Close()

	ctx := context.Background()

	// Stop the container before deleting
	if err := cli.ContainerStop(ctx, containerID, container.StopOptions{}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to stop container"})
		return
	}

	// Remove the container
	if err := cli.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{Force: true}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete container"})
		return
	}

	// Update the container status in the database to "deleted"
	err = UpdateContainerStatusInDB(containerID, "deleted")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update container status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Container deleted", "container_id": containerID})
}

func GetContainerDetails(c *gin.Context) {
	containerID := c.Param("id")

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Docker client"})
		return
	}

	containerJSON, err := cli.ContainerInspect(context.Background(), containerID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Container %s not found", containerID)})
		return
	}

	// Prepare response data
	containerDetails := gin.H{
		"Name":    containerJSON.Name,
		"Image":   containerJSON.Config.Image,
		"Created": containerJSON.Created,
		"State":   containerJSON.State.Status,
	}

	c.JSON(http.StatusOK, containerDetails)
}

func StoreContainerInDB(containerID, imageName, containerName string) error {
	// Create a container object
	container := models.Container{
		ContainerID: containerID,
		Name:        containerName,
		Status:      "created", // Set initial status as created
		BaseImage:   imageName,
		CreatedAt:   time.Now(),
	}

	// Get the MongoDB collection
	collection := config.GetDB().Collection("containers")

	// Insert the container into the database
	_, err := collection.InsertOne(context.Background(), container)
	if err != nil {
		return fmt.Errorf("failed to insert container in DB: %v", err)
	}

	return nil
}

func UpdateContainerStatusInDB(containerID, status string) error {
	// Get the MongoDB collection
	collection := config.GetDB().Collection("containers")

	// Define the filter and update data
	filter := bson.M{"container_id": containerID}
	update := bson.M{"$set": bson.M{"status": status}}

	// Update the container status in the DB
	_, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return fmt.Errorf("failed to update container status in DB: %v", err)
	}

	return nil
}
