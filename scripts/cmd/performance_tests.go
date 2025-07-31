package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"runtime"
	"runtime/pprof"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kevinanielsen/go-fast-cdn/src/database"
	docHandlers "github.com/kevinanielsen/go-fast-cdn/src/handlers/docs"
	imageHandlers "github.com/kevinanielsen/go-fast-cdn/src/handlers/image"
	mediaHandlers "github.com/kevinanielsen/go-fast-cdn/src/handlers/media"
	testutils "github.com/kevinanielsen/go-fast-cdn/src/testUtils"
	"github.com/kevinanielsen/go-fast-cdn/src/util"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// PerformanceTestResult represents the result of a performance test
type PerformanceTestResult struct {
	TestName          string        `json:"test_name"`
	LegacyTime        time.Duration `json:"legacy_time"`
	UnifiedTime       time.Duration `json:"unified_time"`
	Difference        float64       `json:"difference_percent"`
	LegacyMemory      uint64        `json:"legacy_memory_bytes"`
	UnifiedMemory     uint64        `json:"unified_memory_bytes"`
	MemoryDiff        float64       `json:"memory_difference_percent"`
	LegacyThroughput  float64       `json:"legacy_throughput"`
	UnifiedThroughput float64       `json:"unified_throughput"`
	ThroughputDiff    float64       `json:"throughput_difference_percent"`
}

// PerformanceTestSuite represents a suite of performance tests
type PerformanceTestSuite struct {
	results       []PerformanceTestResult
	tempDir       string
	legacyRouter  *gin.Engine
	unifiedRouter *gin.Engine
	mediaHandler  *mediaHandlers.MediaHandler
	imageHandler  *imageHandlers.ImageHandler
	docHandler    *docHandlers.DocHandler
	db            *gorm.DB
}

// NewPerformanceTestSuite creates a new performance test suite
func NewPerformanceTestSuite() *PerformanceTestSuite {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "performance-test")
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
	imageRepo := database.NewImageRepo(db)
	docRepo := database.NewDocRepo(db)

	mediaHandler := mediaHandlers.NewMediaHandler(mediaRepo)
	imageHandler := imageHandlers.NewImageHandler(imageRepo)
	docHandler := docHandlers.NewDocHandler(docRepo)

	// Create routers
	legacyRouter := gin.Default()
	unifiedRouter := gin.Default()

	// Setup legacy routes
	legacyRouter.GET("/api/cdn/images", imageHandler.HandleAllImages)
	legacyRouter.GET("/api/cdn/docs", docHandler.HandleAllDocs)
	legacyRouter.POST("/api/cdn/upload/image", imageHandler.HandleImageUpload)
	legacyRouter.POST("/api/cdn/upload/doc", docHandler.HandleDocUpload)

	// Setup unified routes
	unifiedRouter.GET("/api/cdn/media", mediaHandler.HandleAllMedia)
	unifiedRouter.POST("/api/cdn/upload/media", mediaHandler.HandleMediaUpload)

	return &PerformanceTestSuite{
		tempDir:       tempDir,
		legacyRouter:  legacyRouter,
		unifiedRouter: unifiedRouter,
		mediaHandler:  mediaHandler,
		imageHandler:  imageHandler,
		docHandler:    docHandler,
		db:            db,
	}
}

// Cleanup cleans up resources used by the test suite
func (suite *PerformanceTestSuite) Cleanup() {
	// Remove temporary directory
	os.RemoveAll(suite.tempDir)
}

// MeasureMemoryUsage measures the memory usage before and after a function call
func (suite *PerformanceTestSuite) MeasureMemoryUsage(fn func()) uint64 {
	// Force garbage collection before measuring
	runtime.GC()

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	before := m.Alloc

	fn()

	// Force garbage collection after measuring
	runtime.GC()
	runtime.ReadMemStats(&m)
	after := m.Alloc

	return after - before
}

// BenchmarkDatabaseQueries benchmarks database queries between legacy and unified implementations
func (suite *PerformanceTestSuite) BenchmarkDatabaseQueries() {
	log.Println("Running database query performance tests...")

	// Create test data
	img, err := testutils.CreateDummyImage(200, 200)
	require.NoError(nil, err)

	docContent := testutils.CreateDummyDocument()

	// Calculate checksums (not used in current test implementation)
	_, _ = testutils.CalculateImageChecksum(img)
	_ = testutils.CalculateDocumentChecksum(docContent)

	// Test 1: Get all images vs Get all media with type=image
	suite.runDatabaseTest("Get All Images vs Get All Media (type=image)", func() {
		suite.imageHandler.HandleAllImages(&gin.Context{})
	}, func() {
		suite.mediaHandler.HandleAllMedia(&gin.Context{Request: &http.Request{URL: &url.URL{RawQuery: "type=image"}}})
	})

	// Test 2: Get all docs vs Get all media with type=document
	suite.runDatabaseTest("Get All Docs vs Get All Media (type=document)", func() {
		suite.docHandler.HandleAllDocs(&gin.Context{})
	}, func() {
		req := &http.Request{URL: &url.URL{}}
		req.URL.RawQuery = "type=document"
		suite.mediaHandler.HandleAllMedia(&gin.Context{Request: req})
	})

	// Test 3: Get image by checksum vs Get media by checksum
	suite.runDatabaseTest("Get Image by Checksum vs Get Media by Checksum", func() {
		// Create a test context for the handler
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		suite.imageHandler.HandleAllImages(c)
	}, func() {
		// Create a test context for the handler
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		suite.mediaHandler.HandleAllMedia(c)
	})

	// Test 4: Get doc by checksum vs Get media by checksum
	suite.runDatabaseTest("Get Doc by Checksum vs Get Media by Checksum", func() {
		// Create a test context for the handler
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		suite.docHandler.HandleAllDocs(c)
	}, func() {
		// Create a test context for the handler
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		suite.mediaHandler.HandleAllMedia(c)
	})
}

// runDatabaseTest runs a single database performance test
func (suite *PerformanceTestSuite) runDatabaseTest(testName string, legacyFn, unifiedFn func()) {
	// Measure legacy performance
	var legacyTime time.Duration
	var legacyMemory uint64

	// Warm up
	for i := 0; i < 10; i++ {
		legacyFn()
	}

	// Measure time
	start := time.Now()
	for i := 0; i < 100; i++ {
		legacyFn()
	}
	legacyTime = time.Since(start)

	// Measure memory
	legacyMemory = suite.MeasureMemoryUsage(func() {
		for i := 0; i < 100; i++ {
			legacyFn()
		}
	})

	// Measure unified performance
	var unifiedTime time.Duration
	var unifiedMemory uint64

	// Warm up
	for i := 0; i < 10; i++ {
		unifiedFn()
	}

	// Measure time
	start = time.Now()
	for i := 0; i < 100; i++ {
		unifiedFn()
	}
	unifiedTime = time.Since(start)

	// Measure memory
	unifiedMemory = suite.MeasureMemoryUsage(func() {
		for i := 0; i < 100; i++ {
			unifiedFn()
		}
	})

	// Calculate differences
	timeDiff := float64(legacyTime-unifiedTime) / float64(legacyTime) * 100
	memoryDiff := float64(int64(legacyMemory)-int64(unifiedMemory)) / float64(legacyMemory) * 100

	// Store result
	result := PerformanceTestResult{
		TestName:      testName,
		LegacyTime:    legacyTime,
		UnifiedTime:   unifiedTime,
		Difference:    timeDiff,
		LegacyMemory:  legacyMemory,
		UnifiedMemory: unifiedMemory,
		MemoryDiff:    memoryDiff,
	}

	suite.results = append(suite.results, result)

	log.Printf("Database Test: %s - Legacy: %v, Unified: %v, Difference: %.2f%%", testName, legacyTime, unifiedTime, timeDiff)
}

// BenchmarkAPIEndpoints benchmarks API endpoints between legacy and unified implementations
func (suite *PerformanceTestSuite) BenchmarkAPIEndpoints() {
	log.Println("Running API endpoint performance tests...")

	// Test 1: GET all images vs GET all media with type=image
	suite.runAPITest("GET /api/cdn/images vs GET /api/cdn/media?type=image",
		"GET", "/api/cdn/images",
		"GET", "/api/cdn/media?type=image")

	// Test 2: GET all docs vs GET all media with type=document
	suite.runAPITest("GET /api/cdn/docs vs GET /api/cdn/media?type=document",
		"GET", "/api/cdn/docs",
		"GET", "/api/cdn/media?type=document")

	// Test 3: POST upload image vs POST upload media (image)
	suite.runAPITest("POST /api/cdn/upload/image vs POST /api/cdn/upload/media (image)",
		"POST", "/api/cdn/upload/image",
		"POST", "/api/cdn/upload/media", "image")

	// Test 4: POST upload doc vs POST upload media (document)
	suite.runAPITest("POST /api/cdn/upload/doc vs POST /api/cdn/upload/media (document)",
		"POST", "/api/cdn/upload/doc",
		"POST", "/api/cdn/upload/media", "document")
}

// runAPITest runs a single API performance test
func (suite *PerformanceTestSuite) runAPITest(testName, legacyMethod, legacyPath, unifiedMethod, unifiedPath string, mediaType ...string) {
	// Create test data
	var img image.Image
	var docContent []byte
	var err error

	if len(mediaType) > 0 && mediaType[0] == "image" {
		img, err = testutils.CreateDummyImage(200, 200)
		require.NoError(nil, err)
	} else {
		docContent = testutils.CreateDummyDocument()
	}

	// Measure legacy performance
	var legacyTime time.Duration
	var legacyMemory uint64

	// Warm up
	for i := 0; i < 5; i++ {
		suite.makeRequest(suite.legacyRouter, legacyMethod, legacyPath, img, docContent)
	}

	// Measure time
	start := time.Now()
	for i := 0; i < 20; i++ {
		suite.makeRequest(suite.legacyRouter, legacyMethod, legacyPath, img, docContent)
	}
	legacyTime = time.Since(start)

	// Measure memory
	legacyMemory = suite.MeasureMemoryUsage(func() {
		for i := 0; i < 20; i++ {
			suite.makeRequest(suite.legacyRouter, legacyMethod, legacyPath, img, docContent)
		}
	})

	// Measure unified performance
	var unifiedTime time.Duration
	var unifiedMemory uint64

	// Warm up
	for i := 0; i < 5; i++ {
		suite.makeRequest(suite.unifiedRouter, unifiedMethod, unifiedPath, img, docContent)
	}

	// Measure time
	start = time.Now()
	for i := 0; i < 20; i++ {
		suite.makeRequest(suite.unifiedRouter, unifiedMethod, unifiedPath, img, docContent)
	}
	unifiedTime = time.Since(start)

	// Measure memory
	unifiedMemory = suite.MeasureMemoryUsage(func() {
		for i := 0; i < 20; i++ {
			suite.makeRequest(suite.unifiedRouter, unifiedMethod, unifiedPath, img, docContent)
		}
	})

	// Calculate differences
	timeDiff := float64(legacyTime-unifiedTime) / float64(legacyTime) * 100
	memoryDiff := float64(int64(legacyMemory)-int64(unifiedMemory)) / float64(legacyMemory) * 100

	// Store result
	result := PerformanceTestResult{
		TestName:      testName,
		LegacyTime:    legacyTime,
		UnifiedTime:   unifiedTime,
		Difference:    timeDiff,
		LegacyMemory:  legacyMemory,
		UnifiedMemory: unifiedMemory,
		MemoryDiff:    memoryDiff,
	}

	suite.results = append(suite.results, result)

	log.Printf("API Test: %s - Legacy: %v, Unified: %v, Difference: %.2f%%", testName, legacyTime, unifiedTime, timeDiff)
}

// makeRequest makes an HTTP request to the specified router
func (suite *PerformanceTestSuite) makeRequest(router *gin.Engine, method, path string, img image.Image, docContent []byte) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	if method == "GET" {
		c.Request = httptest.NewRequest(method, path, nil)
	} else {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		var part io.Writer
		var err error

		if img != nil {
			part, err = writer.CreateFormFile("image", "test-image.png")
			if err != nil {
				log.Fatal("Failed to create form file:", err)
			}
			err = testutils.EncodeImage(part, img)
			if err != nil {
				log.Fatal("Failed to encode image:", err)
			}
		} else {
			part, err = writer.CreateFormFile("doc", "test-document.txt")
			if err != nil {
				log.Fatal("Failed to create form file:", err)
			}
			_, err = part.Write(docContent)
			if err != nil {
				log.Fatal("Failed to write document content:", err)
			}
		}

		err = writer.Close()
		if err != nil {
			log.Fatal("Failed to close writer:", err)
		}

		c.Request = httptest.NewRequest(method, path, body)
		c.Request.Header.Add("Content-Type", writer.FormDataContentType())
	}

	router.ServeHTTP(w, c.Request)
}

// BenchmarkConcurrentRequests benchmarks concurrent request handling
func (suite *PerformanceTestSuite) BenchmarkConcurrentRequests() {
	log.Println("Running concurrent request performance tests...")

	// Test 1: Concurrent GET requests
	suite.runConcurrentTest("Concurrent GET Requests", "GET", "/api/cdn/images", "GET", "/api/cdn/media?type=image", 50, 10)

	// Test 2: Concurrent POST requests
	suite.runConcurrentTest("Concurrent POST Requests", "POST", "/api/cdn/upload/image", "POST", "/api/cdn/upload/media", 20, 5)
}

// runConcurrentTest runs a concurrent performance test
func (suite *PerformanceTestSuite) runConcurrentTest(testName, legacyMethod, legacyPath, unifiedMethod, unifiedPath string, numRequests, concurrent int) {
	// Create test data
	img, err := testutils.CreateDummyImage(200, 200)
	require.NoError(nil, err)

	// Measure legacy throughput
	start := time.Now()
	var wg sync.WaitGroup
	requestsPerGoroutine := numRequests / concurrent

	for i := 0; i < concurrent; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < requestsPerGoroutine; j++ {
				suite.makeRequest(suite.legacyRouter, legacyMethod, legacyPath, img, nil)
			}
		}()
	}
	wg.Wait()
	legacyDuration := time.Since(start)
	legacyThroughput := float64(numRequests) / legacyDuration.Seconds()

	// Measure unified throughput
	start = time.Now()
	for i := 0; i < concurrent; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < requestsPerGoroutine; j++ {
				suite.makeRequest(suite.unifiedRouter, unifiedMethod, unifiedPath, img, nil)
			}
		}()
	}
	wg.Wait()
	unifiedDuration := time.Since(start)
	unifiedThroughput := float64(numRequests) / unifiedDuration.Seconds()

	// Calculate differences
	throughputDiff := float64(legacyThroughput-unifiedThroughput) / float64(legacyThroughput) * 100

	// Store result
	result := PerformanceTestResult{
		TestName:          testName,
		LegacyThroughput:  legacyThroughput,
		UnifiedThroughput: unifiedThroughput,
		ThroughputDiff:    throughputDiff,
	}

	suite.results = append(suite.results, result)

	log.Printf("Concurrent Test: %s - Legacy: %.2f req/s, Unified: %.2f req/s, Difference: %.2f%%", testName, legacyThroughput, unifiedThroughput, throughputDiff)
}

// BenchmarkDifferentFileSizes benchmarks performance with different file sizes
func (suite *PerformanceTestSuite) BenchmarkDifferentFileSizes() {
	log.Println("Running different file size performance tests...")

	// Test with different image sizes
	sizes := []struct {
		name   string
		width  int
		height int
	}{
		{"Small Image (100x100)", 100, 100},
		{"Medium Image (500x500)", 500, 500},
		{"Large Image (1000x1000)", 1000, 1000},
	}

	for _, size := range sizes {
		suite.runFileSizeTest(size.name, size.width, size.height)
	}
}

// runFileSizeTest runs a file size performance test
func (suite *PerformanceTestSuite) runFileSizeTest(testName string, width, height int) {
	// Create test image
	img, err := testutils.CreateDummyImage(width, height)
	require.NoError(nil, err)

	// Measure legacy performance
	var legacyTime time.Duration
	var legacyMemory uint64

	// Warm up
	for i := 0; i < 3; i++ {
		suite.makeRequest(suite.legacyRouter, "POST", "/api/cdn/upload/image", img, nil)
	}

	// Measure time
	start := time.Now()
	for i := 0; i < 10; i++ {
		suite.makeRequest(suite.legacyRouter, "POST", "/api/cdn/upload/image", img, nil)
	}
	legacyTime = time.Since(start)

	// Measure memory
	legacyMemory = suite.MeasureMemoryUsage(func() {
		for i := 0; i < 10; i++ {
			suite.makeRequest(suite.legacyRouter, "POST", "/api/cdn/upload/image", img, nil)
		}
	})

	// Measure unified performance
	var unifiedTime time.Duration
	var unifiedMemory uint64

	// Warm up
	for i := 0; i < 3; i++ {
		suite.makeRequest(suite.unifiedRouter, "POST", "/api/cdn/upload/media", img, nil)
	}

	// Measure time
	start = time.Now()
	for i := 0; i < 10; i++ {
		suite.makeRequest(suite.unifiedRouter, "POST", "/api/cdn/upload/media", img, nil)
	}
	unifiedTime = time.Since(start)

	// Measure memory
	unifiedMemory = suite.MeasureMemoryUsage(func() {
		for i := 0; i < 10; i++ {
			suite.makeRequest(suite.unifiedRouter, "POST", "/api/cdn/upload/media", img, nil)
		}
	})

	// Calculate differences
	timeDiff := float64(legacyTime-unifiedTime) / float64(legacyTime) * 100
	memoryDiff := float64(int64(legacyMemory)-int64(unifiedMemory)) / float64(legacyMemory) * 100

	// Store result
	result := PerformanceTestResult{
		TestName:      testName,
		LegacyTime:    legacyTime,
		UnifiedTime:   unifiedTime,
		Difference:    timeDiff,
		LegacyMemory:  legacyMemory,
		UnifiedMemory: unifiedMemory,
		MemoryDiff:    memoryDiff,
	}

	suite.results = append(suite.results, result)

	log.Printf("File Size Test: %s - Legacy: %v, Unified: %v, Difference: %.2f%%", testName, legacyTime, unifiedTime, timeDiff)
}

// GenerateReport generates a performance test report
func (suite *PerformanceTestSuite) GenerateReport() string {
	report := "# Unified Media Repository Performance Test Report\n\n"
	report += fmt.Sprintf("Generated at: %s\n\n", time.Now().Format(time.RFC3339))

	report += "## Test Results\n\n"
	report += "| Test Name | Legacy Time | Unified Time | Difference (%) | Legacy Memory | Unified Memory | Memory Diff (%) | Legacy Throughput | Unified Throughput | Throughput Diff (%) |\n"
	report += "|------------|-------------|-------------|----------------|--------------|---------------|-----------------|-------------------|-------------------|-------------------|\n"

	for _, result := range suite.results {
		legacyTimeStr := result.LegacyTime.String()
		unifiedTimeStr := result.UnifiedTime.String()

		if result.LegacyTime == 0 {
			legacyTimeStr = "N/A"
		}

		if result.UnifiedTime == 0 {
			unifiedTimeStr = "N/A"
		}

		legacyMemoryStr := fmt.Sprintf("%d", result.LegacyMemory)
		unifiedMemoryStr := fmt.Sprintf("%d", result.UnifiedMemory)

		if result.LegacyMemory == 0 {
			legacyMemoryStr = "N/A"
		}

		if result.UnifiedMemory == 0 {
			unifiedMemoryStr = "N/A"
		}

		legacyThroughputStr := fmt.Sprintf("%.2f", result.LegacyThroughput)
		unifiedThroughputStr := fmt.Sprintf("%.2f", result.UnifiedThroughput)

		if result.LegacyThroughput == 0 {
			legacyThroughputStr = "N/A"
		}

		if result.UnifiedThroughput == 0 {
			unifiedThroughputStr = "N/A"
		}

		report += fmt.Sprintf("| %s | %s | %s | %.2f | %s | %s | %.2f | %s | %s | %.2f |\n",
			result.TestName,
			legacyTimeStr,
			unifiedTimeStr,
			result.Difference,
			legacyMemoryStr,
			unifiedMemoryStr,
			result.MemoryDiff,
			legacyThroughputStr,
			unifiedThroughputStr,
			result.ThroughputDiff)
	}

	report += "\n## Analysis\n\n"

	// Calculate average performance difference
	var totalTimeDiff, totalMemoryDiff, totalThroughputDiff float64
	var timeCount, memoryCount, throughputCount int

	for _, result := range suite.results {
		if result.LegacyTime > 0 && result.UnifiedTime > 0 {
			totalTimeDiff += result.Difference
			timeCount++
		}

		if result.LegacyMemory > 0 && result.UnifiedMemory > 0 {
			totalMemoryDiff += result.MemoryDiff
			memoryCount++
		}

		if result.LegacyThroughput > 0 && result.UnifiedThroughput > 0 {
			totalThroughputDiff += result.ThroughputDiff
			throughputCount++
		}
	}

	if timeCount > 0 {
		avgTimeDiff := totalTimeDiff / float64(timeCount)
		report += fmt.Sprintf("- **Average Time Difference**: %.2f%% (%s)\n", avgTimeDiff,
			func() string {
				if avgTimeDiff > 0 {
					return "Unified is faster"
				}
				return "Legacy is faster"
			}())
	}

	if memoryCount > 0 {
		avgMemoryDiff := totalMemoryDiff / float64(memoryCount)
		report += fmt.Sprintf("- **Average Memory Difference**: %.2f%% (%s)\n", avgMemoryDiff,
			func() string {
				if avgMemoryDiff > 0 {
					return "Unified uses less memory"
				}
				return "Legacy uses less memory"
			}())
	}

	if throughputCount > 0 {
		avgThroughputDiff := totalThroughputDiff / float64(throughputCount)
		report += fmt.Sprintf("- **Average Throughput Difference**: %.2f%% (%s)\n", avgThroughputDiff,
			func() string {
				if avgThroughputDiff > 0 {
					return "Legacy has higher throughput"
				}
				return "Unified has higher throughput"
			}())
	}

	report += "\n## Recommendations\n\n"

	// Add recommendations based on results
	for _, result := range suite.results {
		if result.LegacyTime > 0 && result.UnifiedTime > 0 {
			if result.Difference < -10 {
				report += fmt.Sprintf("- **Optimize %s**: Unified implementation is significantly slower. Consider reviewing database queries and implementing caching.\n", result.TestName)
			} else if result.Difference > 10 {
				report += fmt.Sprintf("- **%s**: Unified implementation shows good performance improvement.\n", result.TestName)
			}
		}

		if result.LegacyMemory > 0 && result.UnifiedMemory > 0 {
			if result.MemoryDiff < -10 {
				report += fmt.Sprintf("- **Memory Optimization for %s**: Unified implementation uses significantly more memory. Review memory allocation patterns.\n", result.TestName)
			}
		}
	}

	report += "\n## General Recommendations\n\n"
	report += "- Continue to monitor performance metrics in production environments\n"
	report += "- Consider implementing caching mechanisms for frequently accessed media\n"
	report += "- Optimize database queries and indexes for better performance\n"
	report += "- Implement load balancing for high-traffic scenarios\n"
	report += "- Consider using a content delivery network (CDN) for media files\n"
	report += "- Regularly review and optimize the performance of both unified and legacy endpoints\n"

	return report
}

func main() {
	// Start CPU profiling
	cpuProfile, err := os.Create("cpu.prof")
	if err != nil {
		log.Fatal("could not create CPU profile: ", err)
	}
	defer cpuProfile.Close()
	if err := pprof.StartCPUProfile(cpuProfile); err != nil {
		log.Fatal("could not start CPU profile: ", err)
	}
	defer pprof.StopCPUProfile()

	// Create performance test suite
	suite := NewPerformanceTestSuite()
	defer suite.Cleanup()

	// Run performance tests
	suite.BenchmarkDatabaseQueries()
	suite.BenchmarkAPIEndpoints()
	suite.BenchmarkConcurrentRequests()
	suite.BenchmarkDifferentFileSizes()

	// Generate report
	report := suite.GenerateReport()

	// Save report to file
	reportFile := "performance-report.md"
	err = os.WriteFile(reportFile, []byte(report), 0644)
	if err != nil {
		log.Fatal("Failed to write report file:", err)
	}

	// Also save results as JSON
	resultsFile := "performance-results.json"
	jsonData, err := json.MarshalIndent(suite.results, "", "  ")
	if err != nil {
		log.Fatal("Failed to marshal results:", err)
	}
	err = os.WriteFile(resultsFile, jsonData, 0644)
	if err != nil {
		log.Fatal("Failed to write results file:", err)
	}

	log.Println("Performance testing completed!")
	log.Printf("Report saved to: %s", reportFile)
	log.Printf("Results saved to: %s", resultsFile)

	// Print summary
	log.Println("\n=== Performance Test Summary ===")
	for _, result := range suite.results {
		if result.LegacyTime > 0 && result.UnifiedTime > 0 {
			log.Printf("%s: Legacy=%v, Unified=%v, Diff=%.2f%%",
				result.TestName, result.LegacyTime, result.UnifiedTime, result.Difference)
		}
	}
}
