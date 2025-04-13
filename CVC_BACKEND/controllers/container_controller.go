package controllers

import (
	"CVC_ragh/config"
	"CVC_ragh/models"
	"CVC_ragh/utils"
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func initDockerClient() (*client.Client, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, fmt.Errorf("Docker client error: %v", err)
	}
	return cli, nil
}
func getContainerCollection() *mongo.Collection {
	return config.GetDB().Collection("containers")
}

func getImageCollection() *mongo.Collection {
	return config.GetDB().Collection("images")
}

func isContainerRunning(cli *client.Client, containerName string) (bool, error) {
	containerJSON, err := cli.ContainerInspect(context.Background(), containerName)
	if err != nil {
		return false, fmt.Errorf("Failed to inspect container: %v", err)
	}
	return containerJSON.State.Running, nil
}

func CreateContainerFromImage(c *gin.Context) {
	imageName := c.Param("imageName")              // returns "/iron/node"
	imageName = strings.TrimPrefix(imageName, "/") // remove the leading slash
	fmt.Println("🚀 Creating container from image:", imageName)
	userId, ok := utils.GetUserIdFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Initialize Docker client
	cli, err := initDockerClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Docker client initialization failed", "details": err.Error()})
		return
	}
	defer cli.Close()

	fmt.Println("✅ Docker client initialized successfully!")

	// Create a new container (without manually setting a name)
	resp, err := cli.ContainerCreate(
		context.Background(),
		&container.Config{
			Image: imageName,
			Cmd:   []string{"tail", "-f", "/dev/null"}, // Keep container idle
		},
		nil, nil, nil, "",
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create container", "details": err.Error()})
		return
	}

	fmt.Println("✅ Container created successfully! ID:", resp.ID)

	// Get the default name assigned by Docker
	containerJSON, err := cli.ContainerInspect(context.Background(), resp.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to inspect container", "details": err.Error()})
		return
	}
	containerName := containerJSON.Name // Docker includes a leading "/", remove it
	if len(containerName) > 0 {
		containerName = containerName[1:]
	}

	fmt.Println("📝 Assigned container name:", containerName)

	// Create a Container object to store in DB
	container := models.Container{
		ContainerID: resp.ID,
		Name:        containerName,
		Status:      "created", // Set initial status as created
		BaseImage:   imageName,
		UserID:      userId,
	}

	// Store container details in the database using upsertContainerInDB
	err = insertContainerInDB(container)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save container info", "details": err.Error()})
		return
	}

	fmt.Println("✅ Container details stored in database!")

	// Respond with container ID & name
	c.JSON(http.StatusOK, gin.H{
		"message":        "Container created successfully",
		"container_id":   resp.ID,
		"container_name": containerName,
		"base_image":     imageName,
	})
}

func StartContainer(c *gin.Context) {
	containerName := c.Param("name")
	userId, ok := utils.GetUserIdFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Initialize Docker client
	cli, err := initDockerClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Docker client initialization failed", "details": err.Error()})
		return
	}
	defer cli.Close()

	// Check if the container is running
	isRunning, err := isContainerRunning(cli, containerName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to inspect container", "details": err.Error()})
		return
	}

	// If the container is already running, notify the user
	if isRunning {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Container is already running"})
		return
	}

	// Start the container
	ctx := context.Background()
	if err := cli.ContainerStart(ctx, containerName, types.ContainerStartOptions{}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start container"})
		return
	}

	// Use the upsert function to update the status to "running"
	container := models.Container{
		Name:   containerName,
		UserID: userId,
		Status: "running", // Set the status to "running" since it was started
	}

	err = updateContainerInDB(container)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update container status", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Container started", "container_id": containerName})
}

// StopContainer stops a running container
func StopContainer(c *gin.Context) {
	containerName := c.Param("name")
	userId, ok := utils.GetUserIdFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	cli, err := initDockerClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Docker client initialization failed", "details": err.Error()})
		return
	}
	defer cli.Close()

	// Check if container is running using the helper function
	isRunning, err := isContainerRunning(cli, containerName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to inspect container", "details": err.Error()})
		return
	}

	// If the container isn't running, notify the user
	if !isRunning {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Container is not running"})
		return
	}

	// Stop the container
	ctx := context.Background()
	if err := cli.ContainerStop(ctx, containerName, container.StopOptions{}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to stop container"})
		return
	}

	container := models.Container{
		Name:   containerName,
		UserID: userId,
		Status: "exited", // Set the status to "running" since it was started
	}

	// Update the container status in the database
	err = updateContainerInDB(container)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update container status", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Container exited", "container_id": containerName})
}

// DeleteContainer removes a container and updates its status in the database
func DeleteContainer(c *gin.Context) {
	containerName := c.Param("name") // Ensure this matches the route parameter
	userId, ok := utils.GetUserIdFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	cli, err := initDockerClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Docker client initialization failed", "details": err.Error()})
		return
	}
	defer cli.Close()

	ctx := context.Background()

	// Attempt to stop the container (ignore errors if already stopped or not found)
	_ = cli.ContainerStop(ctx, containerName, container.StopOptions{})

	// Remove the container
	if err := cli.ContainerRemove(ctx, containerName, types.ContainerRemoveOptions{Force: true}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete container"})
		return
	}

	// Update the container status in the database
	collection := getContainerCollection()

	// Define the filter and delete operation
	filter := bson.M{"name": containerName, "userId": userId}
	_, err = collection.DeleteOne(context.Background(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete container from DB"})
		return
	}

	// Respond to the client
	c.JSON(http.StatusOK, gin.H{"message": "Container deleted successfully", "container_name": containerName})
}

func insertContainerInDB(container models.Container) error {
	collection := getContainerCollection()

	_, err := collection.InsertOne(context.Background(), container)
	return err
}

func updateContainerInDB(container models.Container) error {
	collection := getContainerCollection()

	filter := bson.M{"name": container.Name, "userId": container.UserID}

	updateFields := bson.M{}

	if container.ContainerID != "" {
		updateFields["container_id"] = container.ContainerID
	}
	if container.Status != "" {
		updateFields["status"] = container.Status
	}
	if container.BaseImage != "" {
		updateFields["image"] = container.BaseImage
	}

	if len(updateFields) == 0 {
		return nil // nothing to update
	}

	update := bson.M{"$set": updateFields}

	_, err := collection.UpdateOne(context.Background(), filter, update)
	return err
}

func GetImageByContainerName(containerName string) (string, string, error) {
	collection := getContainerCollection() // Get MongoDB collection

	// Define query filter
	filter := bson.M{"name": containerName}

	// Define a variable of type models.Container
	var container models.Container

	// Execute the query
	err := collection.FindOne(context.Background(), filter).Decode(&container)
	if err == mongo.ErrNoDocuments {
		return "", "", fmt.Errorf("no container found with name: %s", containerName)
	} else if err != nil {
		return "", "", fmt.Errorf("database error: %v", err)
	}

	return container.BaseImage, container.UserID, nil
}

func GetUserImages(c *gin.Context) {
	// Get userId from JWT middleware context
	userId, ok := utils.GetUserIdFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	collection := getImageCollection()

	// Query images based on userId, not username
	cursor, err := collection.Find(context.Background(), bson.M{"userId": userId})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch images"})
		return
	}

	var images []models.Image
	if err := cursor.All(context.Background(), &images); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode images"})
		return
	}

	c.JSON(http.StatusOK, images)
}
func GetUserContainers(c *gin.Context) {
	userId, ok := utils.GetUserIdFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Query MongoDB for container *names* owned by the user
	collection := getContainerCollection()
	var userContainers []models.Container
	cursor, err := collection.Find(context.TODO(), bson.M{"userId": userId})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error"})
		return
	}
	defer cursor.Close(context.TODO())
	if err := cursor.All(context.TODO(), &userContainers); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Decode error"})
		return
	}

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Docker client error"})
		return
	}

	var results []gin.H
	for _, userC := range userContainers {
		inspect, err := cli.ContainerInspect(context.Background(), userC.Name)
		if err != nil {
			continue // if deleted recently
		}

		results = append(results, gin.H{
			"Name":    inspect.Name,
			"ID":      inspect.ID,
			"Image":   inspect.Config.Image,
			"State":   inspect.State.Status,
			"Created": inspect.Created,
		})
	}

	c.JSON(http.StatusOK, gin.H{"containers": results})
}
