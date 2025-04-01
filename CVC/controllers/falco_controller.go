package controllers

import (
	"fmt"
	"net/http"
	"regexp"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Prometheus Metrics
var (
	// Existing Metrics
	activeRuleAlerts = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "container_security_rule_alerts_active",
			Help: "Number of currently active security alerts by rule",
		},
		[]string{"rule_name", "priority", "container_id"},
	)

	// New Metrics
	ruleAlertsByImage = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "container_security_alerts_by_image_total",
			Help: "Number of security alerts grouped by container image",
		},
		[]string{"image", "image_tag", "registry"},
	)

	ruleAlertsByCategory = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "container_security_alerts_by_category_total",
			Help: "Number of security alerts grouped by category",
		},
		[]string{"category", "priority"},
	)

	containerViolationsCount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "container_security_violations_count",
			Help: "Number of security violations per container",
		},
		[]string{"container_id", "container_name", "namespace", "pod_name"},
	)

	hostAlertFrequency = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "container_security_host_alerts_total",
			Help: "Number of security alerts by host",
		},
		[]string{"host_id", "host_name", "cluster"},
	)

	ruleEffectiveness = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "container_security_rule_effectiveness",
			Help: "Ratio of true positives to total alerts for a rule (0-1)",
		},
		[]string{"rule_name"},
	)
	metricsMutex               sync.Mutex
	vulnerabilityExposureScore = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "container_vulnerability_exposure_score",
			Help: "Exposure score for container based on vulnerabilities (0-10)",
		},
		[]string{"container_id", "container_name", "image"},
	)
)

func init() {
	prometheus.MustRegister(
		activeRuleAlerts, ruleAlertsByImage, ruleAlertsByCategory, containerViolationsCount,
		hostAlertFrequency, ruleEffectiveness,
		vulnerabilityExposureScore)

}

// FalcoAlert represents the structure of a Falco alert log
type FalcoAlert struct {
	Rule         string `json:"rule"`
	Priority     string `json:"priority"`
	Output       string `json:"output"`
	Time         string `json:"time"`
	Source       string `json:"source"`
	Hostname     string `json:"hostname"`
	Category     string `json:"category"` // Add this field
	OutputFields struct {
		ContainerID   string `json:"container.id"`
		ContainerName string `json:"container.name"`
		ProcessName   string `json:"proc.name"`
	} `json:"output_fields"`
}

// extractImageFromOutput extracts the container image name from the output string
func extractImageFromOutput(output string) string {
	re := regexp.MustCompile(`image=([\w\d:/.-]+)`)
	match := re.FindStringSubmatch(output)
	if len(match) > 1 {
		return match[1]
	}
	return "unknown"
}

// FalcoWebhookHandler processes incoming alerts and updates Prometheus metrics
func FalcoWebhookHandler(c *gin.Context) {
	var falcoEvent FalcoAlert

	// Bind JSON payload to FalcoAlert struct
	if err := c.BindJSON(&falcoEvent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload"})
		return
	}

	// Log received alert
	fmt.Println("🔥 Received Falco Alert:")
	fmt.Printf("  🔥 Rule: %s\n", falcoEvent.Rule)
	fmt.Printf("  🚨 Priority: %s\n", falcoEvent.Priority)
	fmt.Printf("  📝 Output: %s\n", falcoEvent.Output)

	containerID := falcoEvent.OutputFields.ContainerID
	//containerName := falcoEvent.OutputFields.ContainerName
	image := extractImageFromOutput(falcoEvent.Output) // Extract image from output
	metricsMutex.Lock()
	defer metricsMutex.Unlock()

	activeRuleAlerts.WithLabelValues(falcoEvent.Rule, falcoEvent.Priority, containerID).Inc()
	fmt.Println("🔥 Incremented activeRuleAlerts for:", falcoEvent.Rule, falcoEvent.Priority, containerID)

	ruleAlertsByImage.WithLabelValues(image, "latest", "docker.io").Inc()
	ruleAlertsByCategory.WithLabelValues(falcoEvent.OutputFields.ProcessName, falcoEvent.Priority).Inc()
	activeCount := 0
	metricFamilies, err := prometheus.DefaultGatherer.Gather()
	if err == nil {
		for _, mf := range metricFamilies {
			if mf.GetName() == "container_security_rule_alerts_active" {
				for _, m := range mf.GetMetric() {
					activeCount += int(m.GetGauge().GetValue())
				}
			}
		}
	}

	// ✅ Trigger alert if active alerts exceed 10
	if activeCount > 10 {
		fmt.Println("🚨 Alert: More than 10 active security rule alerts detected! 🚨")
		// You can also send this alert to an external system, email, or webhook
	}

	c.JSON(http.StatusOK, gin.H{"message": "Falco alert processed successfully"})

}

// PrometheusMetricsHandler exposes metrics for Prometheus
func PrometheusMetricsHandler(c *gin.Context) {
	promhttp.Handler().ServeHTTP(c.Writer, c.Request)
}
