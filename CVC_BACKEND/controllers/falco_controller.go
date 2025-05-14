package controllers

import (
	"CVC_ragh/config"
	"CVC_ragh/models"
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Prometheus Metrics
var (
	// Existing Metrics
	activeRuleAlerts = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "container_security_rule_alerts_active",
			Help: "Number of currently active security alerts by rule",
		},
		[]string{"rule_name"},
	)
	ruleAlertsByCategory = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "container_security_alerts_by_category_total",
			Help: "Number of security alerts grouped by category",
		},
		[]string{"priority"},
	)

	ruleAlertsByImage = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "container_security_alerts_by_image_total",
			Help: "Number of security alerts grouped by container image",
		},
		[]string{"image", "container_name"},
	)

	containerViolationsCount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "container_security_violations_count",
			Help: "Number of security violations per container",
		},
		[]string{"container_name", "rule_name"},
	)
)

func init() {
	prometheus.MustRegister(
		activeRuleAlerts, ruleAlertsByImage, ruleAlertsByCategory, containerViolationsCount)

}

func FalcoWebhookHandler(c *gin.Context) {
	// Temporary struct to handle Falco's incoming JSON structure
	type IncomingFalcoEvent struct {
		Rule     string `json:"rule"`
		Priority string `json:"priority"`
		Source   string `json:"source"`
		Output   string `json:"output"`

		OutputFields struct {
			ContainerID   string `json:"container.id"`
			ContainerName string `json:"container.name"`
			ProcessName   string `json:"proc.name"`
		} `json:"output_fields"`
	}

	var incoming IncomingFalcoEvent
	if err := c.BindJSON(&incoming); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload"})
		return
	}

	// Logging the alert
	fmt.Println("🔥 Received Falco Alert:")
	fmt.Printf("  🔥 Rule: %s\n", incoming.Rule)
	fmt.Printf("  🚨 Priority: %s\n", incoming.Priority)
	fmt.Printf("  📝 Container: %s\n", incoming.OutputFields.ContainerName)

	containerName := incoming.OutputFields.ContainerName

	// Fetch image name and userId
	imageName, userId, err := GetImageByContainerName(containerName)
	if err != nil {
		fmt.Println("❌ Error getting image name:", err)
		imageName = "unknown"
		userId = "unknown_user"
	}

	// Build the flattened alert model
	alertToSave := models.FalcoAlert{
		Rule:          incoming.Rule,
		Priority:      incoming.Priority,
		Source:        incoming.Source,
		UserId:        userId,
		Count:         1, // Default count for new insert
		ContainerID:   incoming.OutputFields.ContainerID,
		ContainerName: incoming.OutputFields.ContainerName,
		ProcessName:   incoming.OutputFields.ProcessName,
	}

	// Save or update alert in MongoDB
	if err := SaveFalcoAlertToMongo(alertToSave, userId); err != nil {
		fmt.Println("❌ Failed to save Falco alert:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save alert"})
		return
	}

	// 📊 Prometheus metrics
	activeRuleAlerts.WithLabelValues(incoming.Rule).Inc()
	ruleAlertsByImage.WithLabelValues(imageName, containerName).Inc()
	ruleAlertsByCategory.WithLabelValues(incoming.Priority).Inc()
	containerViolationsCount.WithLabelValues(containerName, incoming.Rule).Inc()

	c.JSON(http.StatusOK, gin.H{"message": "Falco alert processed and stored"})
}

// PrometheusMetricsHandler exposes metrics for Prometheus
func PrometheusMetricsHandler(c *gin.Context) {
	promhttp.Handler().ServeHTTP(c.Writer, c.Request)
}
func SaveFalcoAlertToMongo(alert models.FalcoAlert, userId string) error {
	collection := config.GetDB().Collection("falco_alerts")

	// Set user ID and ensure count starts at 1 for new documents
	alert.UserId = userId
	alert.Count = 1

	// Generate a new ID only for new documents
	alert.ID = primitive.NewObjectID()
	filter := bson.M{
		"userId":      userId,
		"rule":        alert.Rule,
		"priority":    alert.Priority,
		"source":      alert.Source,
		"containerId": alert.ContainerID,
	}

	// First, try to find an existing document
	var existingDoc models.FalcoAlert
	err := collection.FindOne(context.Background(), filter).Decode(&existingDoc)

	if err == nil {
		// Document exists, increment the count
		// Only increment the count, don't update other fields
		update := bson.M{
			"$inc": bson.M{"count": 1},
		}

		_, updateErr := collection.UpdateOne(
			context.Background(),
			bson.M{"_id": existingDoc.ID}, // Use the _id for a precise update
			update,
		)

		fmt.Println("Document updated, count incremented")
		return updateErr
	} else {
		// Insert new document
		_, insertErr := collection.InsertOne(context.Background(), alert)
		fmt.Println("New document inserted")
		return insertErr
	}
}
func GetFalcoAlerts(c *gin.Context) {
	userIdRaw, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userId, ok := userIdRaw.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid userId"})
		return
	}
	var falcoCollection = config.GetDB().Collection("falco_alerts")

	filter := bson.M{}
	if userId != "" {
		filter["userId"] = userId
	}

	cur, err := falcoCollection.Find(context.TODO(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching alerts: " + err.Error()})
		return
	}
	defer cur.Close(context.TODO())

	var alerts []models.FalcoAlert
	if err = cur.All(context.TODO(), &alerts); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding alerts"})
		return
	}

	c.JSON(http.StatusOK, alerts)
}
