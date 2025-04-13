package controllers

import (
	"CVC_ragh/config"
	"CVC_ragh/models"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TrivyScanResult struct {
	ImageName     string         `json:"image_name"`
	SeverityCount map[string]int `json:"severity_count"`
	ScanTime      string         `json:"scan_time"`
	UserID        string         `bson:"userId" json:"userId"`
}

// POST /pull-image/:imageName
func PullImageAndSave(c *gin.Context) {
	imageName := c.Param("imageName")
	imageName = strings.TrimPrefix(imageName, "/")

	if imageName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Image name is required"})
		return
	}

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

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Docker client init failed", "details": err.Error()})
		return
	}
	cli.NegotiateAPIVersion(context.Background())

	// Check if image already in DB for this user
	existing, err := GetImageByLabelAndUserID(imageName, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB check failed"})
		return
	}
	if existing != nil {
		c.JSON(http.StatusConflict, gin.H{"message": "Image has already been pulled"})
		return
	}

	// Try inspecting the image locally
	imageDetails, _, err := cli.ImageInspectWithRaw(context.Background(), imageName)
	if err != nil {
		// If not found locally, pull from Docker Hub
		if client.IsErrNotFound(err) || strings.Contains(err.Error(), "No such image") {
			out, err := cli.ImagePull(context.Background(), imageName, types.ImagePullOptions{})
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to pull image", "details": err.Error()})
				return
			}
			defer out.Close()

			// Optional: Stream pull output or wait a bit for pull to complete
			time.Sleep(10 * time.Second)

			// Try inspecting again after pulling
			imageDetails, _, err = cli.ImageInspectWithRaw(context.Background(), imageName)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to inspect image after pull", "details": err.Error()})
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to inspect image", "details": err.Error()})
			return
		}
	}

	repository, tag := parseRepositoryAndTag(imageName)

	err = SaveImageDetails(repository, tag, imageDetails.Created, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image to DB", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Image pulled (or found locally) and saved successfully"})
}

func ScanImageOnly(c *gin.Context) {
	imageName := c.Param("imageName")
	imageName = strings.TrimPrefix(imageName, "/")
	fmt.Println(imageName)

	if imageName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Image name is required"})
		return
	}

	userIdRaw, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userId, ok := userIdRaw.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid userId type"})
		return
	}
	// Add this above the Trivy scan command
	existingScan, err := GetScanByImageNameAndUserID(imageName, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB check failed"})
		return
	}
	if existingScan != nil {
		c.JSON(http.StatusConflict, gin.H{"message": "Image has already been scanned"})
		return
	}

	cmd := exec.Command("trivy", "image", "--quiet", "--format", "json", imageName)
	var output bytes.Buffer
	cmd.Stdout = &output
	if err := cmd.Run(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Trivy scan failed",
			"details": err.Error(),
		})
		return
	}

	var trivyResults map[string]interface{}
	if err := json.Unmarshal(output.Bytes(), &trivyResults); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse Trivy results"})
		return
	}

	severityCount := map[string]int{"CRITICAL": 0, "HIGH": 0, "MEDIUM": 0, "LOW": 0}
	if results, ok := trivyResults["Results"].([]interface{}); ok {
		for _, result := range results {
			if resultMap, ok := result.(map[string]interface{}); ok {
				if vulnerabilities, ok := resultMap["Vulnerabilities"].([]interface{}); ok {
					for _, vuln := range vulnerabilities {
						if vulnMap, ok := vuln.(map[string]interface{}); ok {
							if severity, ok := vulnMap["Severity"].(string); ok {
								if _, exists := severityCount[severity]; exists {
									severityCount[severity]++
								}
							}
						}
					}
				}
			}
		}
	}

	scanResult := TrivyScanResult{
		ImageName:     imageName,
		SeverityCount: severityCount,
		ScanTime:      time.Now().UTC().Format(time.RFC3339),
		UserID:        userId,
	}

	if err := SaveTrivyScan(scanResult.ImageName, scanResult.SeverityCount, scanResult.UserID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save scan result"})
		return
	}

	// Ensure the PDF folder exists
	if err := os.MkdirAll("pdf_reports", os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create PDF folder"})
		return
	}

	// Sanitize filename
	sanitizedImageName := strings.ReplaceAll(strings.ReplaceAll(imageName, "/", "_"), ":", "_")
	pdfFileName := fmt.Sprintf("trivy_report_%s_%d.pdf", sanitizedImageName, time.Now().Unix())
	pdfFilePath := filepath.Join("pdf_reports", pdfFileName)

	if err := GenerateUserFriendlyPDF(trivyResults, scanResult, pdfFilePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "PDF generation failed",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "Scan successful",
		"scan_result":   scanResult,
		"pdf_file_path": pdfFilePath,
	})
}

// ------------------ Helpers ------------------

func parseRepositoryAndTag(image string) (string, string) {
	tag := "latest"
	lastColon := strings.LastIndex(image, ":")
	if lastColon > -1 && !strings.Contains(image[lastColon:], "/") {
		return image[:lastColon], image[lastColon+1:]
	}
	return image, tag
}

func SaveImageDetails(repository, tag string, created string, userId string) error {
	createdTime, err := time.Parse(time.RFC3339, created)
	if err != nil {
		return fmt.Errorf("parse creation time failed: %v", err)
	}

	collection := config.GetDB().Collection("images")
	image := models.Image{
		ID:         primitive.NewObjectID(),
		Label:      fmt.Sprintf("%s:%s", repository, tag),
		Repository: repository,
		Tag:        tag,
		CreateddAt: createdTime.Unix(),
		UserId:     userId,
	}
	_, err = collection.InsertOne(context.Background(), image)
	return err
}

func SaveTrivyScan(imageName string, severityCount map[string]int, userId string) error {
	collection := config.GetDB().Collection("trivy_scans")
	scan := models.TrivyScan{
		ImageName:     imageName,
		SeverityCount: severityCount,
		ScanTime:      time.Now(),
		UserId:        userId,
	}
	_, err := collection.InsertOne(context.Background(), scan)
	return err
}

// ----------------- PDF Report -------------------

func GenerateUserFriendlyPDF(scanResult map[string]interface{}, summary TrivyScanResult, filePath string) error {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(10, 10, 10)
	pdf.AddPage()

	// Title
	pdf.SetFont("Arial", "B", 20)
	pdf.Cell(190, 10, "Trivy Vulnerability Scan Report")
	pdf.Ln(12)
	// Metadata
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(190, 8, fmt.Sprintf("Image: %s", summary.ImageName))
	pdf.Ln(6)
	//pdf.Cell(190, 8, fmt.Sprintf("Scanned By: %s"))
	pdf.Ln(6)
	pdf.Cell(190, 8, fmt.Sprintf("Scan Time: %s", summary.ScanTime))
	pdf.Ln(12)

	// Severity Summary
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(190, 10, "Vulnerability Severity Summary")
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 12)
	pdf.CellFormat(60, 8, "Severity", "1", 0, "", false, 0, "")
	pdf.CellFormat(60, 8, "Count", "1", 1, "", false, 0, "")

	for severity, count := range summary.SeverityCount {
		pdf.CellFormat(60, 8, severity, "1", 0, "", false, 0, "")
		pdf.CellFormat(60, 8, fmt.Sprintf("%d", count), "1", 1, "", false, 0, "")
	}
	pdf.Ln(10)

	// Top Vulnerabilities
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(190, 10, "Sample Vulnerabilities")
	pdf.Ln(10)
	pdf.SetFont("Arial", "", 11)

	// Extract and list a few vulnerabilities
	if results, ok := scanResult["Results"].([]interface{}); ok {
		for _, result := range results {
			if resultMap, ok := result.(map[string]interface{}); ok {
				if vulns, ok := resultMap["Vulnerabilities"].([]interface{}); ok {
					count := 0
					for _, v := range vulns {
						if vmap, ok := v.(map[string]interface{}); ok && count < 5 {
							title := vmap["Title"]
							cve := vmap["VulnerabilityID"]
							severity := vmap["Severity"]
							desc := vmap["Description"]
							fixedVersion := vmap["FixedVersion"]

							pdf.SetFont("Arial", "B", 11)
							pdf.Cell(190, 8, fmt.Sprintf("%s [%s]", title, severity))
							pdf.Ln(6)

							pdf.SetFont("Arial", "", 11)
							pdf.MultiCell(0, 6, fmt.Sprintf(
								"CVE: %s\n%s\n Fix Available: %s",
								cve,
								desc,
								fixedVersion,
							), "", "", false)
							pdf.Ln(4)

							count++
						}
					}
				}
			}
		}
	}

	// Save the file
	return pdf.OutputFileAndClose(filePath)
}
func GetScanByImageNameAndUserID(imageName, userId string) (*TrivyScanResult, error) {
	collection := config.GetDB().Collection("trivy_scans")
	filter := bson.M{"image_name": imageName, "userId": userId}

	var result TrivyScanResult
	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // not found, safe to continue
		}
		return nil, err
	}

	return &result, nil
}
func GetImageByLabelAndUserID(label, userId string) (*models.Image, error) {
	collection := config.GetDB().Collection("images")
	filter := bson.M{"label": label, "userId": userId}

	var result models.Image
	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &result, nil
}
func GetUserScanResults(c *gin.Context) {
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

	collection := config.GetDB().Collection("trivy_scans") // Collection name should match your DB

	// Define a filter to find documents matching the userId
	filter := bson.M{"userId": userId}

	// Find documents
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch scan results"})
		return
	}
	defer cursor.Close(context.TODO())

	// Decode into slice of TrivyScan
	var scanResults []models.TrivyScan
	if err := cursor.All(context.TODO(), &scanResults); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode scan results"})
		return
	}

	// Return results
	c.JSON(http.StatusOK, scanResults)
}
