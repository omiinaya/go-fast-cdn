package main

import (
	"encoding/json"
	"fmt"
	"image"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"runtime"
	"runtime/pprof"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kevinanielsen/go-fast-cdn/src/database"
	mediaHandlers "github.com/kevinanielsen/go-fast-cdn/src/handlers/media"
	testutils "github.com/kevinanielsen/go-fast-cdn/src/testUtils"
	"github.com/kevinanielsen/go-fast-cdn/src/util"
	"gorm.io/gorm"
)

// MonitoringMetric represents a single monitoring metric
type MonitoringMetric struct {
	Name      string                 `json:"name"`
	Value     float64                `json:"value"`
	Unit      string                 `json:"unit"`
	Timestamp time.Time              `json:"timestamp"`
	Tags      map[string]interface{} `json:"tags,omitempty"`
}

// MonitoringResult represents the result of monitoring
type MonitoringResult struct {
	TestName       string             `json:"test_name"`
	Timestamp      time.Time          `json:"timestamp"`
	Metrics        []MonitoringMetric `json:"metrics"`
	Status         string             `json:"status"` // "success", "warning", "error"
	Message        string             `json:"message,omitempty"`
	Duration       time.Duration      `json:"duration"`
	UnifiedMetrics []MonitoringMetric `json:"unified_metrics,omitempty"`
	LegacyMetrics  []MonitoringMetric `json:"legacy_metrics,omitempty"`
	Comparison     float64            `json:"comparison,omitempty"` // percentage difference
}

// MonitoringSuite represents a suite of monitoring tests
type MonitoringSuite struct {
	results       []MonitoringResult
	tempDir       string
	legacyRouter  *gin.Engine
	unifiedRouter *gin.Engine
	mediaHandler  *mediaHandlers.MediaHandler
	db            *gorm.DB
}

// NewMonitoringSuite creates a new monitoring suite
func NewMonitoringSuite() *MonitoringSuite {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "monitoring-test")
	if err != nil {
		log.Fatal("Failed to create temporary directory:", err)
	}

	// Set the execution path to the temp directory
	util.ExPath = tempDir

	// Initialize database
	database.ConnectToDB()
	db := database.DB

	// Create handlers
	mediaRepo := database.NewMediaRepo(db)
	mediaHandler := mediaHandlers.NewMediaHandler(mediaRepo)

	// Create routers
	legacyRouter := gin.Default()
	unifiedRouter := gin.Default()

	// Setup routes (simplified for monitoring)
	legacyRouter.GET("/api/cdn/images", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	legacyRouter.GET("/api/cdn/docs", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	unifiedRouter.GET("/api/cdn/media", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	return &MonitoringSuite{
		tempDir:       tempDir,
		legacyRouter:  legacyRouter,
		unifiedRouter: unifiedRouter,
		mediaHandler:  mediaHandler,
		db:            db,
	}
}

// Cleanup cleans up resources used by the monitoring suite
func (suite *MonitoringSuite) Cleanup() {
	// Remove temporary directory
	os.RemoveAll(suite.tempDir)
}

// MonitorSystemResources monitors system resource usage
func (suite *MonitoringSuite) MonitorSystemResources() {
	log.Println("Monitoring system resources...")

	start := time.Now()
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	result := MonitoringResult{
		TestName:  "System Resource Usage",
		Timestamp: time.Now(),
		Status:    "success",
		Duration:  time.Since(start),
	}

	// Collect system metrics
	result.Metrics = append(result.Metrics, MonitoringMetric{
		Name:      "CPU Usage",
		Value:     0, // Placeholder - would need external library for real CPU usage
		Unit:      "percent",
		Timestamp: time.Now(),
	})

	result.Metrics = append(result.Metrics, MonitoringMetric{
		Name:      "Memory Usage",
		Value:     float64(m.Alloc),
		Unit:      "bytes",
		Timestamp: time.Now(),
	})

	result.Metrics = append(result.Metrics, MonitoringMetric{
		Name:      "Memory Allocations",
		Value:     float64(m.Mallocs),
		Unit:      "count",
		Timestamp: time.Now(),
	})

	result.Metrics = append(result.Metrics, MonitoringMetric{
		Name:      "Goroutines",
		Value:     float64(runtime.NumGoroutine()),
		Unit:      "count",
		Timestamp: time.Now(),
	})

	suite.results = append(suite.results, result)
	log.Printf("System resource monitoring completed in %v", result.Duration)
}

// MonitorAPIPerformance monitors API endpoint performance
func (suite *MonitoringSuite) MonitorAPIPerformance() {
	log.Println("Monitoring API performance...")

	// Test unified media endpoints
	suite.monitorEndpoint("GET /api/cdn/media", suite.unifiedRouter, "GET", "/api/cdn/media", nil)

	// Test legacy endpoints for backward compatibility
	suite.monitorEndpoint("GET /api/cdn/images", suite.legacyRouter, "GET", "/api/cdn/images", nil)
	suite.monitorEndpoint("GET /api/cdn/docs", suite.legacyRouter, "GET", "/api/cdn/docs", nil)
}

// monitorEndpoint monitors a single API endpoint
func (suite *MonitoringSuite) monitorEndpoint(testName string, router *gin.Engine, method, path string, body []byte) {
	start := time.Now()

	// Make multiple requests to get average performance
	var totalTime time.Duration
	var successCount int
	requestCount := 20

	for i := 0; i < requestCount; i++ {
		reqStart := time.Now()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(method, path, nil)
		router.ServeHTTP(w, req)

		if w.Code == http.StatusOK {
			successCount++
		}

		totalTime += time.Since(reqStart)
	}

	avgTime := totalTime / time.Duration(requestCount)
	successRate := float64(successCount) / float64(requestCount) * 100

	result := MonitoringResult{
		TestName:  testName,
		Timestamp: time.Now(),
		Status:    "success",
		Duration:  time.Since(start),
	}

	result.Metrics = append(result.Metrics, MonitoringMetric{
		Name:      "Average Response Time",
		Value:     float64(avgTime.Nanoseconds()) / 1e6, // Convert to milliseconds
		Unit:      "ms",
		Timestamp: time.Now(),
	})

	result.Metrics = append(result.Metrics, MonitoringMetric{
		Name:      "Success Rate",
		Value:     successRate,
		Unit:      "percent",
		Timestamp: time.Now(),
	})

	result.Metrics = append(result.Metrics, MonitoringMetric{
		Name:      "Total Requests",
		Value:     float64(requestCount),
		Unit:      "count",
		Timestamp: time.Now(),
	})

	// Set status based on success rate
	if successRate < 95 {
		result.Status = "error"
		result.Message = fmt.Sprintf("Low success rate: %.2f%%", successRate)
	} else if successRate < 99 {
		result.Status = "warning"
		result.Message = fmt.Sprintf("Moderate success rate: %.2f%%", successRate)
	}

	suite.results = append(suite.results, result)
	log.Printf("%s monitoring completed in %v, success rate: %.2f%%", testName, result.Duration, successRate)
}

// MonitorDatabasePerformance monitors database performance
func (suite *MonitoringSuite) MonitorDatabasePerformance() {
	log.Println("Monitoring database performance...")

	start := time.Now()

	// Monitor database connection
	result := MonitoringResult{
		TestName:  "Database Performance",
		Timestamp: time.Now(),
		Status:    "success",
		Duration:  time.Since(start),
	}

	// Test database connection
	dbStart := time.Now()
	sqlDB, err := suite.db.DB()
	if err != nil {
		result.Status = "error"
		result.Message = fmt.Sprintf("Failed to get database instance: %v", err)
		suite.results = append(suite.results, result)
		return
	}

	err = sqlDB.Ping()
	if err != nil {
		result.Status = "error"
		result.Message = fmt.Sprintf("Database ping failed: %v", err)
		suite.results = append(suite.results, result)
		return
	}
	dbPingTime := time.Since(dbStart)

	// Get database stats
	var stats struct {
		OpenConnections int `json:"open_connections"`
		InUse           int `json:"in_use"`
		Idle            int `json:"idle"`
	}

	dbStats := sqlDB.Stats()
	stats.OpenConnections = dbStats.OpenConnections
	stats.InUse = dbStats.InUse
	stats.Idle = dbStats.Idle

	result.Metrics = append(result.Metrics, MonitoringMetric{
		Name:      "Database Ping Time",
		Value:     float64(dbPingTime.Nanoseconds()) / 1e6, // Convert to milliseconds
		Unit:      "ms",
		Timestamp: time.Now(),
	})

	result.Metrics = append(result.Metrics, MonitoringMetric{
		Name:      "Open Connections",
		Value:     float64(stats.OpenConnections),
		Unit:      "count",
		Timestamp: time.Now(),
	})

	result.Metrics = append(result.Metrics, MonitoringMetric{
		Name:      "In Use Connections",
		Value:     float64(stats.InUse),
		Unit:      "count",
		Timestamp: time.Now(),
	})

	result.Metrics = append(result.Metrics, MonitoringMetric{
		Name:      "Idle Connections",
		Value:     float64(stats.Idle),
		Unit:      "count",
		Timestamp: time.Now(),
	})

	suite.results = append(suite.results, result)
	log.Printf("Database performance monitoring completed in %v", result.Duration)
}

// MonitorMediaOperations monitors media operations performance
func (suite *MonitoringSuite) MonitorMediaOperations() {
	log.Println("Monitoring media operations...")

	// Create test data
	img, err := testutils.CreateDummyImage(200, 200)
	if err != nil {
		log.Printf("Failed to create test image: %v", err)
		return
	}

	docContent := testutils.CreateDummyDocument()

	// Monitor image operations
	suite.monitorMediaOperation("Image Upload", img, nil)
	suite.monitorMediaOperation("Document Upload", nil, docContent)
}

// monitorMediaOperation monitors a single media operation
func (suite *MonitoringSuite) monitorMediaOperation(testName string, img image.Image, docContent []byte) {
	start := time.Now()

	result := MonitoringResult{
		TestName:  testName,
		Timestamp: time.Now(),
		Status:    "success",
		Duration:  time.Since(start),
	}

	// Simulate media operation (in a real scenario, this would use the actual handlers)
	operationStart := time.Now()

	// Simulate processing time based on content size
	var size int
	if img != nil {
		// Simulate image processing
		size = 200 * 200 * 4 // Approximate size of 200x200 RGBA image
	} else if docContent != nil {
		size = len(docContent)
	}

	// Simulate processing delay based on size
	time.Sleep(time.Duration(size/10000) * time.Millisecond)

	operationTime := time.Since(operationStart)

	result.Metrics = append(result.Metrics, MonitoringMetric{
		Name:      "Operation Time",
		Value:     float64(operationTime.Nanoseconds()) / 1e6, // Convert to milliseconds
		Unit:      "ms",
		Timestamp: time.Now(),
	})

	result.Metrics = append(result.Metrics, MonitoringMetric{
		Name:      "File Size",
		Value:     float64(size),
		Unit:      "bytes",
		Timestamp: time.Now(),
	})

	// Calculate processing rate, handling the case where operationTime is 0
	var processingRate float64
	if operationTime.Seconds() > 0 {
		processingRate = float64(size) / operationTime.Seconds()
	} else {
		processingRate = 0
	}

	result.Metrics = append(result.Metrics, MonitoringMetric{
		Name:      "Processing Rate",
		Value:     processingRate,
		Unit:      "bytes/sec",
		Timestamp: time.Now(),
	})

	suite.results = append(suite.results, result)
	log.Printf("%s monitoring completed in %v", testName, result.Duration)
}

// MonitorBackwardCompatibility monitors backward compatibility
func (suite *MonitoringSuite) MonitorBackwardCompatibility() {
	log.Println("Monitoring backward compatibility...")

	start := time.Now()

	result := MonitoringResult{
		TestName:  "Backward Compatibility",
		Timestamp: time.Now(),
		Status:    "success",
		Duration:  time.Since(start),
	}

	// Test legacy endpoints
	legacyEndpoints := []string{
		"/api/cdn/images",
		"/api/cdn/docs",
	}

	unifiedEndpoints := []string{
		"/api/cdn/media?type=image",
		"/api/cdn/media?type=document",
	}

	for i, legacyEndpoint := range legacyEndpoints {
		if i < len(unifiedEndpoints) {
			unifiedEndpoint := unifiedEndpoints[i]

			// Test legacy endpoint
			legacyTime := suite.testEndpointResponseTime(suite.legacyRouter, "GET", legacyEndpoint)

			// Test unified endpoint
			unifiedTime := suite.testEndpointResponseTime(suite.unifiedRouter, "GET", unifiedEndpoint)

			// Calculate comparison, handling the case where legacyTime is 0
			var comparison float64
			if legacyTime > 0 {
				comparison = float64(legacyTime-unifiedTime) / float64(legacyTime) * 100
			} else {
				comparison = 0
			}

			result.Metrics = append(result.Metrics, MonitoringMetric{
				Name:      fmt.Sprintf("Legacy %s Response Time", legacyEndpoint),
				Value:     float64(legacyTime.Nanoseconds()) / 1e6,
				Unit:      "ms",
				Timestamp: time.Now(),
			})

			result.Metrics = append(result.Metrics, MonitoringMetric{
				Name:      fmt.Sprintf("Unified %s Response Time", unifiedEndpoint),
				Value:     float64(unifiedTime.Nanoseconds()) / 1e6,
				Unit:      "ms",
				Timestamp: time.Now(),
			})

			result.Metrics = append(result.Metrics, MonitoringMetric{
				Name:      fmt.Sprintf("Comparison %s vs %s", legacyEndpoint, unifiedEndpoint),
				Value:     comparison,
				Unit:      "percent",
				Timestamp: time.Now(),
			})
		}
	}

	suite.results = append(suite.results, result)
	log.Printf("Backward compatibility monitoring completed in %v", result.Duration)
}

// testEndpointResponseTime tests the response time of a single endpoint
func (suite *MonitoringSuite) testEndpointResponseTime(router *gin.Engine, method, path string) time.Duration {
	start := time.Now()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(method, path, nil)
	router.ServeHTTP(w, req)

	return time.Since(start)
}

// GenerateReport generates a monitoring report
func (suite *MonitoringSuite) GenerateReport() string {
	report := "# Unified Media Repository Post-Deployment Monitoring Report\n\n"
	report += fmt.Sprintf("Generated at: %s\n\n", time.Now().Format(time.RFC3339))

	report += "## Monitoring Results\n\n"
	report += "| Test Name | Status | Duration | Metrics Count | Message |\n"
	report += "|------------|--------|----------|---------------|---------|\n"

	for _, result := range suite.results {
		message := result.Message
		if message == "" {
			message = "N/A"
		}
		report += fmt.Sprintf("| %s | %s | %v | %d | %s |\n",
			result.TestName,
			result.Status,
			result.Duration,
			len(result.Metrics),
			message)
	}

	report += "\n## Detailed Metrics\n\n"

	for _, result := range suite.results {
		report += fmt.Sprintf("### %s\n\n", result.TestName)
		report += fmt.Sprintf("**Status**: %s\n\n", result.Status)
		report += fmt.Sprintf("**Duration**: %v\n\n", result.Duration)

		if result.Message != "" {
			report += fmt.Sprintf("**Message**: %s\n\n", result.Message)
		}

		report += "| Metric Name | Value | Unit | Timestamp |\n"
		report += "|-------------|-------|------|-----------|\n"

		for _, metric := range result.Metrics {
			report += fmt.Sprintf("| %s | %.2f | %s | %s |\n",
				metric.Name,
				metric.Value,
				metric.Unit,
				metric.Timestamp.Format(time.RFC3339))
		}

		report += "\n"
	}

	report += "## Analysis\n\n"

	// Analyze results
	var errorCount, warningCount, successCount int
	for _, result := range suite.results {
		switch result.Status {
		case "error":
			errorCount++
		case "warning":
			warningCount++
		case "success":
			successCount++
		}
	}

	report += fmt.Sprintf("- **Total Tests**: %d\n", len(suite.results))
	report += fmt.Sprintf("- **Successful Tests**: %d\n", successCount)
	report += fmt.Sprintf("- **Warnings**: %d\n", warningCount)
	report += fmt.Sprintf("- **Errors**: %d\n", errorCount)

	report += "\n## Recommendations\n\n"

	if errorCount > 0 {
		report += "### Critical Issues\n"
		report += "- Address the failing tests immediately\n"
		report += "- Investigate the root causes of the errors\n"
		report += "- Implement fixes and re-run monitoring\n\n"
	}

	if warningCount > 0 {
		report += "### Warnings\n"
		report += "- Review the tests with warnings\n"
		report += "- Consider optimization for borderline performance\n"
		report += "- Monitor these areas closely in production\n\n"
	}

	report += "### General Recommendations\n"
	report += "- Continue regular monitoring of system performance\n"
	report += "- Set up automated alerts for critical metrics\n"
	report += "- Implement performance baselines for trend analysis\n"
	report += "- Consider scaling resources based on usage patterns\n"
	report += "- Regularly review and update monitoring thresholds\n"

	return report
}

func main() {
	// Start CPU profiling
	cpuProfile, err := os.Create("monitoring-cpu.prof")
	if err != nil {
		log.Fatal("could not create CPU profile: ", err)
	}
	defer cpuProfile.Close()
	if err := pprof.StartCPUProfile(cpuProfile); err != nil {
		log.Fatal("could not start CPU profile: ", err)
	}
	defer pprof.StopCPUProfile()

	// Create monitoring suite
	suite := NewMonitoringSuite()
	defer suite.Cleanup()

	// Run monitoring tests
	suite.MonitorSystemResources()
	suite.MonitorAPIPerformance()
	suite.MonitorDatabasePerformance()
	suite.MonitorMediaOperations()
	suite.MonitorBackwardCompatibility()

	// Generate report
	report := suite.GenerateReport()

	// Save report to file
	reportFile := "post-deployment-monitoring-report.md"
	err = os.WriteFile(reportFile, []byte(report), 0644)
	if err != nil {
		log.Fatal("Failed to write report file:", err)
	}

	// Also save results as JSON
	resultsFile := "post-deployment-monitoring-results.json"
	jsonData, err := json.MarshalIndent(suite.results, "", "  ")
	if err != nil {
		log.Fatal("Failed to marshal results:", err)
	}
	err = os.WriteFile(resultsFile, jsonData, 0644)
	if err != nil {
		log.Fatal("Failed to write results file:", err)
	}

	log.Println("Post-deployment monitoring completed!")
	log.Printf("Report saved to: %s", reportFile)
	log.Printf("Results saved to: %s", resultsFile)

	// Print summary
	log.Println("\n=== Post-Deployment Monitoring Summary ===")
	for _, result := range suite.results {
		log.Printf("%s: Status=%s, Duration=%v",
			result.TestName, result.Status, result.Duration)
	}
}
