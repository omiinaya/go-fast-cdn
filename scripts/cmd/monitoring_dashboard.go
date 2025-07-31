package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

// MonitoringData represents the structure for monitoring data
type MonitoringData struct {
	TestName       string             `json:"test_name"`
	Timestamp      time.Time          `json:"timestamp"`
	Metrics        []MonitoringMetric `json:"metrics"`
	Status         string             `json:"status"`
	Message        string             `json:"message,omitempty"`
	Duration       time.Duration      `json:"duration"`
	UnifiedMetrics []MonitoringMetric `json:"unified_metrics,omitempty"`
	LegacyMetrics  []MonitoringMetric `json:"legacy_metrics,omitempty"`
	Comparison     float64            `json:"comparison,omitempty"`
}

// MonitoringMetric represents a single monitoring metric
type MonitoringMetric struct {
	Name      string                 `json:"name"`
	Value     float64                `json:"value"`
	Unit      string                 `json:"unit"`
	Timestamp time.Time              `json:"timestamp"`
	Tags      map[string]interface{} `json:"tags,omitempty"`
}

// DashboardServer represents the monitoring dashboard server
type DashboardServer struct {
	router         *gin.Engine
	reportsDir     string
	port           string
	monitoringData []MonitoringData
}

// NewDashboardServer creates a new dashboard server
func NewDashboardServer(reportsDir, port string) *DashboardServer {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	server := &DashboardServer{
		router:     router,
		reportsDir: reportsDir,
		port:       port,
	}

	server.setupRoutes()
	return server
}

// setupRoutes sets up the routes for the dashboard server
func (s *DashboardServer) setupRoutes() {
	// Serve static files
	s.router.Static("/static", "./static")

	// Serve the main HTML file
	s.router.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})

	// API endpoints
	s.router.GET("/api/monitoring-data", s.getMonitoringData)
	s.router.GET("/api/monitoring-summary", s.getMonitoringSummary)
	s.router.GET("/api/system-health", s.getSystemHealth)

	// WebSocket for real-time updates (simplified implementation)
	s.router.GET("/ws/monitoring-updates", s.handleWebSocketUpdates)
}

// getMonitoringData returns the monitoring data
func (s *DashboardServer) getMonitoringData(c *gin.Context) {
	// Try to load the latest monitoring data
	data, err := s.loadMonitoringData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}

// getMonitoringSummary returns a summary of the monitoring data
func (s *DashboardServer) getMonitoringSummary(c *gin.Context) {
	data, err := s.loadMonitoringData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	summary := s.generateSummary(data)
	c.JSON(http.StatusOK, summary)
}

// getSystemHealth returns the system health status
func (s *DashboardServer) getSystemHealth(c *gin.Context) {
	// In a real implementation, this would check actual system health
	health := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now(),
		"checks": map[string]interface{}{
			"database": "healthy",
			"api":      "healthy",
			"storage":  "healthy",
		},
		"uptime": "24h 15m 30s",
	}

	c.JSON(http.StatusOK, health)
}

// handleWebSocketUpdates handles WebSocket connections for real-time updates
func (s *DashboardServer) handleWebSocketUpdates(c *gin.Context) {
	// This is a simplified implementation
	// In a real-world scenario, you would use a proper WebSocket library
	c.JSON(http.StatusNotImplemented, gin.H{"message": "WebSocket updates not implemented in this demo"})
}

// loadMonitoringData loads the monitoring data from the reports directory
func (s *DashboardServer) loadMonitoringData() ([]MonitoringData, error) {
	// Try to find the latest monitoring results file
	resultsFile := filepath.Join(s.reportsDir, "post-deployment-monitoring-results.json")

	// If the file doesn't exist, generate sample data
	if _, err := os.Stat(resultsFile); os.IsNotExist(err) {
		return s.generateSampleData(), nil
	}

	// Read the file
	data, err := os.ReadFile(resultsFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read monitoring results file: %v", err)
	}

	// Parse the JSON
	var monitoringData []MonitoringData
	if err := json.Unmarshal(data, &monitoringData); err != nil {
		return nil, fmt.Errorf("failed to parse monitoring results: %v", err)
	}

	return monitoringData, nil
}

// generateSampleData generates sample monitoring data for demonstration
func (s *DashboardServer) generateSampleData() []MonitoringData {
	now := time.Now()

	return []MonitoringData{
		{
			TestName:  "System Resource Usage",
			Timestamp: now,
			Status:    "success",
			Duration:  15 * time.Millisecond,
			Metrics: []MonitoringMetric{
				{
					Name:      "CPU Usage",
					Value:     25.5,
					Unit:      "percent",
					Timestamp: now,
				},
				{
					Name:      "Memory Usage",
					Value:     2048000,
					Unit:      "bytes",
					Timestamp: now,
				},
				{
					Name:      "Memory Allocations",
					Value:     1250,
					Unit:      "count",
					Timestamp: now,
				},
				{
					Name:      "Goroutines",
					Value:     12,
					Unit:      "count",
					Timestamp: now,
				},
			},
		},
		{
			TestName:  "API Performance",
			Timestamp: now,
			Status:    "success",
			Duration:  45 * time.Millisecond,
			Metrics: []MonitoringMetric{
				{
					Name:      "Average Response Time",
					Value:     2.5,
					Unit:      "ms",
					Timestamp: now,
				},
				{
					Name:      "Success Rate",
					Value:     100,
					Unit:      "percent",
					Timestamp: now,
				},
				{
					Name:      "Total Requests",
					Value:     20,
					Unit:      "count",
					Timestamp: now,
				},
			},
		},
		{
			TestName:  "Database Performance",
			Timestamp: now,
			Status:    "success",
			Duration:  5 * time.Millisecond,
			Metrics: []MonitoringMetric{
				{
					Name:      "Database Ping Time",
					Value:     0.5,
					Unit:      "ms",
					Timestamp: now,
				},
				{
					Name:      "Open Connections",
					Value:     1,
					Unit:      "count",
					Timestamp: now,
				},
				{
					Name:      "In Use Connections",
					Value:     0,
					Unit:      "count",
					Timestamp: now,
				},
				{
					Name:      "Idle Connections",
					Value:     1,
					Unit:      "count",
					Timestamp: now,
				},
			},
		},
		{
			TestName:  "Media Operations",
			Timestamp: now,
			Status:    "success",
			Duration:  25 * time.Millisecond,
			Metrics: []MonitoringMetric{
				{
					Name:      "Operation Time",
					Value:     10.2,
					Unit:      "ms",
					Timestamp: now,
				},
				{
					Name:      "File Size",
					Value:     160000,
					Unit:      "bytes",
					Timestamp: now,
				},
				{
					Name:      "Processing Rate",
					Value:     15686274.5,
					Unit:      "bytes/sec",
					Timestamp: now,
				},
			},
		},
		{
			TestName:  "Backward Compatibility",
			Timestamp: now,
			Status:    "success",
			Duration:  30 * time.Millisecond,
			Metrics: []MonitoringMetric{
				{
					Name:      "Legacy /api/cdn/images Response Time",
					Value:     1.8,
					Unit:      "ms",
					Timestamp: now,
				},
				{
					Name:      "Unified /api/cdn/media?type=image Response Time",
					Value:     1.9,
					Unit:      "ms",
					Timestamp: now,
				},
				{
					Name:      "Comparison /api/cdn/images vs /api/cdn/media?type=image",
					Value:     -5.56,
					Unit:      "percent",
					Timestamp: now,
				},
				{
					Name:      "Legacy /api/cdn/docs Response Time",
					Value:     1.7,
					Unit:      "ms",
					Timestamp: now,
				},
				{
					Name:      "Unified /api/cdn/media?type=document Response Time",
					Value:     1.8,
					Unit:      "ms",
					Timestamp: now,
				},
				{
					Name:      "Comparison /api/cdn/docs vs /api/cdn/media?type=document",
					Value:     -5.88,
					Unit:      "percent",
					Timestamp: now,
				},
			},
		},
	}
}

// generateSummary generates a summary from the monitoring data
func (s *DashboardServer) generateSummary(data []MonitoringData) map[string]interface{} {
	totalTests := len(data)
	successfulTests := 0
	warningTests := 0
	errorTests := 0

	for _, test := range data {
		switch test.Status {
		case "success":
			successfulTests++
		case "warning":
			warningTests++
		case "error":
			errorTests++
		}
	}

	successRate := float64(0)
	if totalTests > 0 {
		successRate = float64(successfulTests) / float64(totalTests) * 100
	}

	// Extract average response time
	var totalResponseTime float64
	var responseTimeCount int

	for _, test := range data {
		for _, metric := range test.Metrics {
			if metric.Name == "Average Response Time" {
				totalResponseTime += metric.Value
				responseTimeCount++
			}
		}
	}

	avgResponseTime := float64(0)
	if responseTimeCount > 0 {
		avgResponseTime = totalResponseTime / float64(responseTimeCount)
	}

	summary := map[string]interface{}{
		"total_tests":       totalTests,
		"successful_tests":  successfulTests,
		"warning_tests":     warningTests,
		"error_tests":       errorTests,
		"success_rate":      successRate,
		"avg_response_time": avgResponseTime,
		"timestamp":         time.Now(),
		"status":            "operational",
	}

	if errorTests > 0 {
		summary["status"] = "critical"
	} else if warningTests > 0 {
		summary["status"] = "warning"
	}

	return summary
}

// Run starts the dashboard server
func (s *DashboardServer) Run() error {
	// Ensure the static directory exists
	if _, err := os.Stat("./static"); os.IsNotExist(err) {
		log.Println("Static directory not found, creating it...")
		if err := os.Mkdir("./static", 0755); err != nil {
			return fmt.Errorf("failed to create static directory: %v", err)
		}
	}

	// Check if the index.html file exists
	if _, err := os.Stat("./static/index.html"); os.IsNotExist(err) {
		log.Println("index.html not found in static directory")
	}

	log.Printf("Starting monitoring dashboard server on port %s", s.port)
	log.Printf("Monitoring reports directory: %s", s.reportsDir)
	log.Printf("Access the dashboard at: http://localhost%s", s.port)

	return s.router.Run(s.port)
}

func main() {
	// Configuration
	reportsDir := "../post-deployment-monitoring-reports"
	port := ":8080"

	// Allow configuration via environment variables
	if envReportsDir := os.Getenv("MONITORING_REPORTS_DIR"); envReportsDir != "" {
		reportsDir = envReportsDir
	}

	if envPort := os.Getenv("MONITORING_DASHBOARD_PORT"); envPort != "" {
		port = ":" + envPort
	}

	// Create and run the server
	server := NewDashboardServer(reportsDir, port)
	if err := server.Run(); err != nil {
		log.Fatalf("Failed to start dashboard server: %v", err)
	}
}
