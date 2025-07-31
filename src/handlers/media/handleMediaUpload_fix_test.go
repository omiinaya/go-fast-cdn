package handlers

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/kevinanielsen/go-fast-cdn/src/database"
	"github.com/kevinanielsen/go-fast-cdn/src/models"
	testutils "github.com/kevinanielsen/go-fast-cdn/src/testUtils"
	"github.com/kevinanielsen/go-fast-cdn/src/util"
	"github.com/stretchr/testify/require"
)

// TestMultipleFilesWithSameNameDifferentContent tests that multiple files with the same filename but different content can be uploaded successfully
func TestMultipleFilesWithSameNameDifferentContent(t *testing.T) {
	fmt.Println("\n=== Test: Multiple files with same name but different content ===")

	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "media-test-same-name")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	util.ExPath = tempDir

	// Initialize database
	database.ConnectToDB()
	defer func() {
		dbPath := filepath.Join(util.ExPath, database.DbFolder, database.DbName)
		if err := os.Remove(dbPath); err != nil && !os.IsNotExist(err) {
			t.Error(err)
		}
	}()

	// Create media handler
	mediaHandler := NewMediaHandler(database.NewMediaRepo(database.DB))

	// Test case 1: Upload first image with name "test.png"
	fmt.Println("Uploading first image with name 'test.png'...")
	img1, err := testutils.CreateDummyImage(200, 200)
	require.NoError(t, err)

	body1 := &bytes.Buffer{}
	writer1 := multipart.NewWriter(body1)

	part1, err := writer1.CreateFormFile("file", "test.png")
	require.NoError(t, err)

	err = testutils.EncodeImage(part1, img1)
	require.NoError(t, err)

	err = writer1.Close()
	require.NoError(t, err)

	w1 := httptest.NewRecorder()
	c1, _ := gin.CreateTestContext(w1)

	c1.Request = httptest.NewRequest(http.MethodPost, "/api/cdn/upload/media", body1)
	c1.Request.Header.Add("Content-Type", writer1.FormDataContentType())

	mediaHandler.HandleMediaUpload(c1)

	// Assert successful response
	require.Equal(t, http.StatusOK, w1.Result().StatusCode)
	responseBody1 := w1.Body.String()
	fmt.Printf("First upload response: %s\n", responseBody1)
	require.Contains(t, responseBody1, "file_url")
	require.Contains(t, responseBody1, "filename")

	// Test case 2: Upload second image with same name "test.png" but different content
	fmt.Println("Uploading second image with same name 'test.png' but different content...")
	img2, err := testutils.CreateDummyImage(300, 300) // Different dimensions = different content
	require.NoError(t, err)

	body2 := &bytes.Buffer{}
	writer2 := multipart.NewWriter(body2)

	part2, err := writer2.CreateFormFile("file", "test.png")
	require.NoError(t, err)

	err = testutils.EncodeImage(part2, img2)
	require.NoError(t, err)

	err = writer2.Close()
	require.NoError(t, err)

	w2 := httptest.NewRecorder()
	c2, _ := gin.CreateTestContext(w2)

	c2.Request = httptest.NewRequest(http.MethodPost, "/api/cdn/upload/media", body2)
	c2.Request.Header.Add("Content-Type", writer2.FormDataContentType())

	mediaHandler.HandleMediaUpload(c2)

	// Assert successful response (this should work now that UNIQUE constraint is removed)
	require.Equal(t, http.StatusOK, w2.Result().StatusCode)
	responseBody2 := w2.Body.String()
	fmt.Printf("Second upload response: %s\n", responseBody2)
	require.Contains(t, responseBody2, "file_url")
	require.Contains(t, responseBody2, "filename")

	// Test case 3: Upload third image with same name "test.png" but different content again
	fmt.Println("Uploading third image with same name 'test.png' but different content again...")
	img3, err := testutils.CreateDummyImage(150, 150) // Different dimensions = different content
	require.NoError(t, err)

	body3 := &bytes.Buffer{}
	writer3 := multipart.NewWriter(body3)

	part3, err := writer3.CreateFormFile("file", "test.png")
	require.NoError(t, err)

	err = testutils.EncodeImage(part3, img3)
	require.NoError(t, err)

	err = writer3.Close()
	require.NoError(t, err)

	w3 := httptest.NewRecorder()
	c3, _ := gin.CreateTestContext(w3)

	c3.Request = httptest.NewRequest(http.MethodPost, "/api/cdn/upload/media", body3)
	c3.Request.Header.Add("Content-Type", writer3.FormDataContentType())

	mediaHandler.HandleMediaUpload(c3)

	// Assert successful response
	require.Equal(t, http.StatusOK, w3.Result().StatusCode)
	responseBody3 := w3.Body.String()
	fmt.Printf("Third upload response: %s\n", responseBody3)
	require.Contains(t, responseBody3, "file_url")
	require.Contains(t, responseBody3, "filename")

	// Verify all three files exist in the database with the same filename but different checksums
	mediaRepo := database.NewMediaRepo(database.DB)
	allMedia := mediaRepo.GetAllMedia()

	var testPngFiles []models.Media
	for _, media := range allMedia {
		if strings.Contains(media.FileName, "test.png") {
			testPngFiles = append(testPngFiles, media)
		}
	}

	fmt.Printf("Found %d files with 'test.png' in filename\n", len(testPngFiles))
	require.Equal(t, 3, len(testPngFiles), "Should have exactly 3 files with similar names")

	// Verify all have different checksums
	checksum1 := fmt.Sprintf("%x", testPngFiles[0].Checksum)
	checksum2 := fmt.Sprintf("%x", testPngFiles[1].Checksum)
	checksum3 := fmt.Sprintf("%x", testPngFiles[2].Checksum)

	fmt.Printf("Checksum 1: %s\n", checksum1)
	fmt.Printf("Checksum 2: %s\n", checksum2)
	fmt.Printf("Checksum 3: %s\n", checksum3)

	require.NotEqual(t, checksum1, checksum2, "Checksums should be different")
	require.NotEqual(t, checksum2, checksum3, "Checksums should be different")
	require.NotEqual(t, checksum1, checksum3, "Checksums should be different")

	fmt.Println("✓ Test passed: Multiple files with same name but different content uploaded successfully")
}

// TestFilesWithDifferentNamesIdenticalContent tests that files with different names but identical content are properly rejected as duplicates
func TestFilesWithDifferentNamesIdenticalContent(t *testing.T) {
	fmt.Println("\n=== Test: Files with different names but identical content ===")

	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "media-test-duplicate-content")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	util.ExPath = tempDir

	// Initialize database
	database.ConnectToDB()
	defer func() {
		dbPath := filepath.Join(util.ExPath, database.DbFolder, database.DbName)
		if err := os.Remove(dbPath); err != nil && !os.IsNotExist(err) {
			t.Error(err)
		}
	}()

	// Create media handler
	mediaHandler := NewMediaHandler(database.NewMediaRepo(database.DB))

	// Create test image
	img, err := testutils.CreateDummyImage(200, 200)
	require.NoError(t, err)

	// Calculate checksum for verification
	checksum, err := testutils.CalculateImageChecksum(img)
	require.NoError(t, err)
	expectedChecksum := fmt.Sprintf("%x", checksum)
	fmt.Printf("Expected checksum: %s\n", expectedChecksum)

	// Test case 1: Upload first image with name "first.png"
	fmt.Println("Uploading first image with name 'first.png'...")
	body1 := &bytes.Buffer{}
	writer1 := multipart.NewWriter(body1)

	part1, err := writer1.CreateFormFile("file", "first.png")
	require.NoError(t, err)

	err = testutils.EncodeImage(part1, img)
	require.NoError(t, err)

	err = writer1.Close()
	require.NoError(t, err)

	w1 := httptest.NewRecorder()
	c1, _ := gin.CreateTestContext(w1)

	c1.Request = httptest.NewRequest(http.MethodPost, "/api/cdn/upload/media", body1)
	c1.Request.Header.Add("Content-Type", writer1.FormDataContentType())

	mediaHandler.HandleMediaUpload(c1)

	// Assert successful response
	require.Equal(t, http.StatusOK, w1.Result().StatusCode)
	responseBody1 := w1.Body.String()
	fmt.Printf("First upload response: %s\n", responseBody1)
	require.Contains(t, responseBody1, "file_url")
	require.Contains(t, responseBody1, "filename")

	// Test case 2: Upload second image with different name "second.png" but identical content
	fmt.Println("Uploading second image with different name 'second.png' but identical content...")
	body2 := &bytes.Buffer{}
	writer2 := multipart.NewWriter(body2)

	part2, err := writer2.CreateFormFile("file", "second.png")
	require.NoError(t, err)

	err = testutils.EncodeImage(part2, img) // Same image content
	require.NoError(t, err)

	err = writer2.Close()
	require.NoError(t, err)

	w2 := httptest.NewRecorder()
	c2, _ := gin.CreateTestContext(w2)

	c2.Request = httptest.NewRequest(http.MethodPost, "/api/cdn/upload/media", body2)
	c2.Request.Header.Add("Content-Type", writer2.FormDataContentType())

	mediaHandler.HandleMediaUpload(c2)

	// Assert conflict response (should be rejected as duplicate)
	require.Equal(t, http.StatusConflict, w2.Result().StatusCode)
	responseBody2 := w2.Body.String()
	fmt.Printf("Second upload response: %s\n", responseBody2)
	require.Contains(t, responseBody2, "error")
	require.Contains(t, responseBody2, "File with this content already exists")
	require.Contains(t, responseBody2, "existing_file")

	// Test case 3: Upload third image with different name "third.png" but identical content
	fmt.Println("Uploading third image with different name 'third.png' but identical content...")
	body3 := &bytes.Buffer{}
	writer3 := multipart.NewWriter(body3)

	part3, err := writer3.CreateFormFile("file", "third.png")
	require.NoError(t, err)

	err = testutils.EncodeImage(part3, img) // Same image content
	require.NoError(t, err)

	err = writer3.Close()
	require.NoError(t, err)

	w3 := httptest.NewRecorder()
	c3, _ := gin.CreateTestContext(w3)

	c3.Request = httptest.NewRequest(http.MethodPost, "/api/cdn/upload/media", body3)
	c3.Request.Header.Add("Content-Type", writer3.FormDataContentType())

	mediaHandler.HandleMediaUpload(c3)

	// Assert conflict response (should be rejected as duplicate)
	require.Equal(t, http.StatusConflict, w3.Result().StatusCode)
	responseBody3 := w3.Body.String()
	fmt.Printf("Third upload response: %s\n", responseBody3)
	require.Contains(t, responseBody3, "error")
	require.Contains(t, responseBody3, "File with this content already exists")
	require.Contains(t, responseBody3, "existing_file")

	// Verify only one file exists in the database
	mediaRepo := database.NewMediaRepo(database.DB)
	allMedia := mediaRepo.GetAllMedia()

	var pngFiles []models.Media
	for _, media := range allMedia {
		if strings.HasSuffix(media.FileName, ".png") {
			pngFiles = append(pngFiles, media)
		}
	}

	fmt.Printf("Found %d PNG files in database\n", len(pngFiles))
	require.Equal(t, 1, len(pngFiles), "Should have exactly 1 file since others were duplicates")

	// Verify the checksum matches
	actualChecksum := fmt.Sprintf("%x", pngFiles[0].Checksum)
	require.Equal(t, expectedChecksum, actualChecksum, "Checksum should match the uploaded image")

	fmt.Println("✓ Test passed: Files with different names but identical content properly rejected as duplicates")
}

// TestDifferentFileTypes tests that the upload process works for different file types
func TestDifferentFileTypes(t *testing.T) {
	fmt.Println("\n=== Test: Different file types ===")

	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "media-test-file-types")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	util.ExPath = tempDir

	// Initialize database
	database.ConnectToDB()
	defer func() {
		dbPath := filepath.Join(util.ExPath, database.DbFolder, database.DbName)
		if err := os.Remove(dbPath); err != nil && !os.IsNotExist(err) {
			t.Error(err)
		}
	}()

	// Create media handler
	mediaHandler := NewMediaHandler(database.NewMediaRepo(database.DB))

	// Test case 1: Upload PNG image
	fmt.Println("Testing PNG image upload...")
	img, err := testutils.CreateDummyImage(200, 200)
	require.NoError(t, err)

	body1 := &bytes.Buffer{}
	writer1 := multipart.NewWriter(body1)

	part1, err := writer1.CreateFormFile("file", "test-image.png")
	require.NoError(t, err)

	err = testutils.EncodeImage(part1, img)
	require.NoError(t, err)

	err = writer1.Close()
	require.NoError(t, err)

	w1 := httptest.NewRecorder()
	c1, _ := gin.CreateTestContext(w1)

	c1.Request = httptest.NewRequest(http.MethodPost, "/api/cdn/upload/media", body1)
	c1.Request.Header.Add("Content-Type", writer1.FormDataContentType())

	mediaHandler.HandleMediaUpload(c1)

	require.Equal(t, http.StatusOK, w1.Result().StatusCode)
	responseBody1 := w1.Body.String()
	fmt.Printf("PNG upload response: %s\n", responseBody1)
	require.Contains(t, responseBody1, "file_url")
	require.Contains(t, responseBody1, "type")
	require.Contains(t, responseBody1, "image")

	// Test case 2: Upload text document
	fmt.Println("Testing text document upload...")
	docContent := testutils.CreateDummyDocument()

	body2 := &bytes.Buffer{}
	writer2 := multipart.NewWriter(body2)

	part2, err := writer2.CreateFormFile("file", "test-document.txt")
	require.NoError(t, err)

	_, err = part2.Write(docContent)
	require.NoError(t, err)

	err = writer2.Close()
	require.NoError(t, err)

	w2 := httptest.NewRecorder()
	c2, _ := gin.CreateTestContext(w2)

	c2.Request = httptest.NewRequest(http.MethodPost, "/api/cdn/upload/media", body2)
	c2.Request.Header.Add("Content-Type", writer2.FormDataContentType())

	mediaHandler.HandleMediaUpload(c2)

	require.Equal(t, http.StatusOK, w2.Result().StatusCode)
	responseBody2 := w2.Body.String()
	fmt.Printf("Text document upload response: %s\n", responseBody2)
	require.Contains(t, responseBody2, "file_url")
	require.Contains(t, responseBody2, "type")
	require.Contains(t, responseBody2, "document")

	// Test case 3: Upload JSON document
	fmt.Println("Testing JSON document upload...")
	jsonContent := []byte(`{"test": "data", "number": 123}`)

	body3 := &bytes.Buffer{}
	writer3 := multipart.NewWriter(body3)

	part3, err := writer3.CreateFormFile("file", "test-data.json")
	require.NoError(t, err)

	_, err = part3.Write(jsonContent)
	require.NoError(t, err)

	err = writer3.Close()
	require.NoError(t, err)

	w3 := httptest.NewRecorder()
	c3, _ := gin.CreateTestContext(w3)

	c3.Request = httptest.NewRequest(http.MethodPost, "/api/cdn/upload/media", body3)
	c3.Request.Header.Add("Content-Type", writer3.FormDataContentType())

	mediaHandler.HandleMediaUpload(c3)

	require.Equal(t, http.StatusOK, w3.Result().StatusCode)
	responseBody3 := w3.Body.String()
	fmt.Printf("JSON document upload response: %s\n", responseBody3)
	require.Contains(t, responseBody3, "file_url")
	require.Contains(t, responseBody3, "type")
	require.Contains(t, responseBody3, "document")

	// Verify all files exist in the database
	mediaRepo := database.NewMediaRepo(database.DB)
	allMedia := mediaRepo.GetAllMedia()

	fmt.Printf("Total files in database: %d\n", len(allMedia))
	require.Equal(t, 3, len(allMedia), "Should have exactly 3 files")

	// Verify file types
	var imageCount, documentCount int
	for _, media := range allMedia {
		if media.Type == models.MediaTypeImage {
			imageCount++
		} else if media.Type == models.MediaTypeDocument {
			documentCount++
		}
	}

	require.Equal(t, 1, imageCount, "Should have 1 image")
	require.Equal(t, 2, documentCount, "Should have 2 documents")

	fmt.Println("✓ Test passed: Different file types uploaded successfully")
}

// TestErrorMessageQuality tests that error messages are clear and informative
func TestErrorMessageQuality(t *testing.T) {
	fmt.Println("\n=== Test: Error message quality ===")

	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "media-test-error-messages")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	util.ExPath = tempDir

	// Initialize database
	database.ConnectToDB()
	defer func() {
		dbPath := filepath.Join(util.ExPath, database.DbFolder, database.DbName)
		if err := os.Remove(dbPath); err != nil && !os.IsNotExist(err) {
			t.Error(err)
		}
	}()

	// Create media handler
	mediaHandler := NewMediaHandler(database.NewMediaRepo(database.DB))

	// Test case 1: Duplicate file error message
	fmt.Println("Testing duplicate file error message...")
	img, err := testutils.CreateDummyImage(200, 200)
	require.NoError(t, err)

	// First upload
	body1 := &bytes.Buffer{}
	writer1 := multipart.NewWriter(body1)

	part1, err := writer1.CreateFormFile("file", "original.png")
	require.NoError(t, err)

	err = testutils.EncodeImage(part1, img)
	require.NoError(t, err)

	err = writer1.Close()
	require.NoError(t, err)

	w1 := httptest.NewRecorder()
	c1, _ := gin.CreateTestContext(w1)

	c1.Request = httptest.NewRequest(http.MethodPost, "/api/cdn/upload/media", body1)
	c1.Request.Header.Add("Content-Type", writer1.FormDataContentType())

	mediaHandler.HandleMediaUpload(c1)
	require.Equal(t, http.StatusOK, w1.Result().StatusCode)

	// Second upload (duplicate)
	body2 := &bytes.Buffer{}
	writer2 := multipart.NewWriter(body2)

	part2, err := writer2.CreateFormFile("file", "duplicate.png")
	require.NoError(t, err)

	err = testutils.EncodeImage(part2, img) // Same content
	require.NoError(t, err)

	err = writer2.Close()
	require.NoError(t, err)

	w2 := httptest.NewRecorder()
	c2, _ := gin.CreateTestContext(w2)

	c2.Request = httptest.NewRequest(http.MethodPost, "/api/cdn/upload/media", body2)
	c2.Request.Header.Add("Content-Type", writer2.FormDataContentType())

	mediaHandler.HandleMediaUpload(c2)

	require.Equal(t, http.StatusConflict, w2.Result().StatusCode)
	responseBody2 := w2.Body.String()
	fmt.Printf("Duplicate error response: %s\n", responseBody2)

	// Verify error message quality
	require.Contains(t, responseBody2, "error", "Response should contain 'error' field")
	require.Contains(t, responseBody2, "File with this content already exists", "Should explain it's a content-based duplicate")
	require.Contains(t, responseBody2, "existing_file", "Should provide the existing filename")
	require.Contains(t, responseBody2, "original.png", "Should show the actual existing filename")

	fmt.Println("✓ Test passed: Error messages are clear and informative")
}

// TestDatabaseBehavior tests the database behavior to ensure it matches expectations
func TestDatabaseBehavior(t *testing.T) {
	fmt.Println("\n=== Test: Database behavior ===")

	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "media-test-database-behavior")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	util.ExPath = tempDir

	// Initialize database
	database.ConnectToDB()
	defer func() {
		dbPath := filepath.Join(util.ExPath, database.DbFolder, database.DbName)
		if err := os.Remove(dbPath); err != nil && !os.IsNotExist(err) {
			t.Error(err)
		}
	}()

	// Create media handler
	mediaHandler := NewMediaHandler(database.NewMediaRepo(database.DB))
	mediaRepo := database.NewMediaRepo(database.DB)

	// Create test content with known checksums
	testContent1 := []byte("Test content 1")
	testContent2 := []byte("Test content 2")

	hash1 := md5.Sum(testContent1)
	hash2 := md5.Sum(testContent2)

	checksum1 := fmt.Sprintf("%x", hash1)
	checksum2 := fmt.Sprintf("%x", hash2)

	fmt.Printf("Test content 1 checksum: %s\n", checksum1)
	fmt.Printf("Test content 2 checksum: %s\n", checksum2)

	// Upload first file
	fmt.Println("Uploading first file...")
	body1 := &bytes.Buffer{}
	writer1 := multipart.NewWriter(body1)

	part1, err := writer1.CreateFormFile("file", "same-name.txt")
	require.NoError(t, err)

	_, err = part1.Write(testContent1)
	require.NoError(t, err)

	err = writer1.Close()
	require.NoError(t, err)

	w1 := httptest.NewRecorder()
	c1, _ := gin.CreateTestContext(w1)

	c1.Request = httptest.NewRequest(http.MethodPost, "/api/cdn/upload/media", body1)
	c1.Request.Header.Add("Content-Type", writer1.FormDataContentType())

	mediaHandler.HandleMediaUpload(c1)
	require.Equal(t, http.StatusOK, w1.Result().StatusCode)

	// Upload second file with same name but different content
	fmt.Println("Uploading second file with same name but different content...")
	body2 := &bytes.Buffer{}
	writer2 := multipart.NewWriter(body2)

	part2, err := writer2.CreateFormFile("file", "same-name.txt")
	require.NoError(t, err)

	_, err = part2.Write(testContent2)
	require.NoError(t, err)

	err = writer2.Close()
	require.NoError(t, err)

	w2 := httptest.NewRecorder()
	c2, _ := gin.CreateTestContext(w2)

	c2.Request = httptest.NewRequest(http.MethodPost, "/api/cdn/upload/media", body2)
	c2.Request.Header.Add("Content-Type", writer2.FormDataContentType())

	mediaHandler.HandleMediaUpload(c2)
	require.Equal(t, http.StatusOK, w2.Result().StatusCode)

	// Verify database state
	allMedia := mediaRepo.GetAllMedia()
	fmt.Printf("Total files in database: %d\n", len(allMedia))
	require.Equal(t, 2, len(allMedia), "Should have exactly 2 files")

	// Verify checksums are different
	dbChecksum1 := fmt.Sprintf("%x", allMedia[0].Checksum)
	dbChecksum2 := fmt.Sprintf("%x", allMedia[1].Checksum)

	fmt.Printf("Database checksum 1: %s\n", dbChecksum1)
	fmt.Printf("Database checksum 2: %s\n", dbChecksum2)

	require.NotEqual(t, dbChecksum1, dbChecksum2, "Database checksums should be different")

	// Verify both checksums match our expected values
	foundChecksum1 := (dbChecksum1 == checksum1 || dbChecksum2 == checksum1)
	foundChecksum2 := (dbChecksum1 == checksum2 || dbChecksum2 == checksum2)

	require.True(t, foundChecksum1, "Should find checksum1 in database")
	require.True(t, foundChecksum2, "Should find checksum2 in database")

	// Test database lookup by checksum
	foundMedia1 := mediaRepo.GetMediaByCheckSum(hash1[:])
	require.NotEmpty(t, foundMedia1.Checksum, "Should find media by checksum1")

	foundMedia2 := mediaRepo.GetMediaByCheckSum(hash2[:])
	require.NotEmpty(t, foundMedia2.Checksum, "Should find media by checksum2")

	fmt.Printf("Found media 1: %s (checksum: %s)\n", foundMedia1.FileName, fmt.Sprintf("%x", foundMedia1.Checksum))
	fmt.Printf("Found media 2: %s (checksum: %s)\n", foundMedia2.FileName, fmt.Sprintf("%x", foundMedia2.Checksum))

	fmt.Println("✓ Test passed: Database behavior is correct")
}
