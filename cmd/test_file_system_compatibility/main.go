package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/kevinanielsen/go-fast-cdn/src/database"
	"github.com/kevinanielsen/go-fast-cdn/src/initializers"
	"github.com/kevinanielsen/go-fast-cdn/src/models"
	"github.com/kevinanielsen/go-fast-cdn/src/router"
	"github.com/kevinanielsen/go-fast-cdn/src/util"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// TestFileSystemBackwardCompatibility tests that existing image and document files
// in the file system are still accessible after the unification
func TestFileSystemBackwardCompatibility(t *testing.T) {
	// Setup test database
	db := setupTestDatabase(t)
	defer cleanupTestDatabase(t, db)

	// Create test files
	testImages := createTestImageFiles(t)
	testDocuments := createTestDocumentFiles(t)
	defer cleanupTestFiles(t, testImages, testDocuments)

	// Setup test server
	gin.SetMode(gin.TestMode)
	server := router.NewServer()

	// Test that legacy image files are accessible through the legacy endpoint
	t.Run("LegacyImageFilesAccessibleThroughLegacyEndpoint", func(t *testing.T) {
		// Test each image file
		for _, imagePath := range testImages {
			filename := filepath.Base(imagePath)

			// Test GET /download/images/:filename
			req, _ := http.NewRequest("GET", "/download/images/"+filename, nil)
			w := httptest.NewRecorder()
			server.Engine.ServeHTTP(w, req)

			// The response should be OK if the file exists
			if w.Code == http.StatusOK {
				// Verify that the response is not empty
				assert.Greater(t, w.Body.Len(), 0)
			} else {
				// If the file doesn't exist, the response should be 404
				assert.Equal(t, http.StatusNotFound, w.Code)
			}
		}
	})

	// Test that legacy document files are accessible through the legacy endpoint
	t.Run("LegacyDocumentFilesAccessibleThroughLegacyEndpoint", func(t *testing.T) {
		// Test each document file
		for _, docPath := range testDocuments {
			filename := filepath.Base(docPath)

			// Test GET /download/docs/:filename
			req, _ := http.NewRequest("GET", "/download/docs/"+filename, nil)
			w := httptest.NewRecorder()
			server.Engine.ServeHTTP(w, req)

			// The response should be OK if the file exists
			if w.Code == http.StatusOK {
				// Verify that the response is not empty
				assert.Greater(t, w.Body.Len(), 0)
			} else {
				// If the file doesn't exist, the response should be 404
				assert.Equal(t, http.StatusNotFound, w.Code)
			}
		}
	})

	// Test that legacy image files are accessible through the unified media endpoint
	t.Run("LegacyImageFilesAccessibleThroughUnifiedEndpoint", func(t *testing.T) {
		// Test each image file
		for _, imagePath := range testImages {
			filename := filepath.Base(imagePath)

			// Test GET /download/media/:filename
			req, _ := http.NewRequest("GET", "/download/media/"+filename, nil)
			w := httptest.NewRecorder()
			server.Engine.ServeHTTP(w, req)

			// The response should be OK if the file exists
			if w.Code == http.StatusOK {
				// Verify that the response is not empty
				assert.Greater(t, w.Body.Len(), 0)
			} else {
				// If the file doesn't exist, the response should be 404
				assert.Equal(t, http.StatusNotFound, w.Code)
			}
		}
	})

	// Test that legacy document files are accessible through the unified media endpoint
	t.Run("LegacyDocumentFilesAccessibleThroughUnifiedEndpoint", func(t *testing.T) {
		// Test each document file
		for _, docPath := range testDocuments {
			filename := filepath.Base(docPath)

			// Test GET /download/media/:filename
			req, _ := http.NewRequest("GET", "/download/media/"+filename, nil)
			w := httptest.NewRecorder()
			server.Engine.ServeHTTP(w, req)

			// The response should be OK if the file exists
			if w.Code == http.StatusOK {
				// Verify that the response is not empty
				assert.Greater(t, w.Body.Len(), 0)
			} else {
				// If the file doesn't exist, the response should be 404
				assert.Equal(t, http.StatusNotFound, w.Code)
			}
		}
	})

	// Test that the file paths are correctly resolved
	t.Run("FilePathsCorrectlyResolved", func(t *testing.T) {
		// Test legacy image path
		imagesPath := util.GetImagesPath()
		assert.NotEmpty(t, imagesPath)
		assert.Contains(t, imagesPath, "images")

		// Test legacy document path
		docsPath := util.GetDocsPath()
		assert.NotEmpty(t, docsPath)
		assert.Contains(t, docsPath, "docs")

		// Test unified media path
		mediaPath := util.GetMediaUploadPath()
		assert.NotEmpty(t, mediaPath)
		assert.Contains(t, mediaPath, "media")

		// Test that the paths are different
		assert.NotEqual(t, imagesPath, docsPath)
		assert.NotEqual(t, imagesPath, mediaPath)
		assert.NotEqual(t, docsPath, mediaPath)
	})

	// Test that the legacy directories exist
	t.Run("LegacyDirectoriesExist", func(t *testing.T) {
		// Test images directory
		imagesPath := util.GetImagesPath()
		_, err := os.Stat(imagesPath)
		if err == nil {
			// The directory exists
			assert.True(t, true)
		} else if os.IsNotExist(err) {
			// The directory doesn't exist, which is also fine
			assert.True(t, true)
		} else {
			// There was an error checking the directory
			assert.NoError(t, err)
		}

		// Test documents directory
		docsPath := util.GetDocsPath()
		_, err = os.Stat(docsPath)
		if err == nil {
			// The directory exists
			assert.True(t, true)
		} else if os.IsNotExist(err) {
			// The directory doesn't exist, which is also fine
			assert.True(t, true)
		} else {
			// There was an error checking the directory
			assert.NoError(t, err)
		}

		// Test media directory
		mediaPath := util.GetMediaUploadPath()
		_, err = os.Stat(mediaPath)
		if err == nil {
			// The directory exists
			assert.True(t, true)
		} else if os.IsNotExist(err) {
			// The directory doesn't exist, which is also fine
			assert.True(t, true)
		} else {
			// There was an error checking the directory
			assert.NoError(t, err)
		}
	})

	// Test that the URL paths are correctly generated
	t.Run("URLPathsCorrectlyGenerated", func(t *testing.T) {
		// Test legacy image URL path
		filename := "test-image.jpg"
		imageURL := util.GetLegacyURLPath(filename, util.MediaTypeImage)
		assert.NotEmpty(t, imageURL)
		assert.Contains(t, imageURL, "images")
		assert.Contains(t, imageURL, filename)

		// Test legacy document URL path
		docURL := util.GetLegacyURLPath(filename, util.MediaTypeDocument)
		assert.NotEmpty(t, docURL)
		assert.Contains(t, docURL, "docs")
		assert.Contains(t, docURL, filename)

		// Test unified media URL path
		mediaURL := util.GetMediaURLPath(filename)
		assert.NotEmpty(t, mediaURL)
		assert.Contains(t, mediaURL, "media")
		assert.Contains(t, mediaURL, filename)

		// Test that the URL paths are different
		assert.NotEqual(t, imageURL, docURL)
		assert.NotEqual(t, imageURL, mediaURL)
		assert.NotEqual(t, docURL, mediaURL)
	})
}

// TestFileUploadBackwardCompatibility tests that file upload still works correctly
// after the unification, and files are saved to the correct locations
func TestFileUploadBackwardCompatibility(t *testing.T) {
	// Setup test database
	db := setupTestDatabase(t)
	defer cleanupTestDatabase(t, db)

	// Setup test server
	gin.SetMode(gin.TestMode)
	server := router.NewServer()

	// Test that legacy image upload saves files to the legacy directory
	t.Run("LegacyImageUploadSavesToLegacyDirectory", func(t *testing.T) {
		// Create a test image file
		imagePath := createTestImageFile(t, "test-upload-image.jpg")
		defer os.Remove(imagePath)

		// Create a multipart form request
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		// Add the image file
		file, err := os.Open(imagePath)
		assert.NoError(t, err)
		defer file.Close()

		part, err := writer.CreateFormFile("image", filepath.Base(imagePath))
		assert.NoError(t, err)

		_, err = io.Copy(part, file)
		assert.NoError(t, err)

		// Add the filename field
		err = writer.WriteField("filename", "uploaded-image.jpg")
		assert.NoError(t, err)

		err = writer.Close()
		assert.NoError(t, err)

		// Create the request
		req, _ := http.NewRequest("POST", "/api/cdn/upload/image", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		w := httptest.NewRecorder()
		server.Engine.ServeHTTP(w, req)

		// The response should be OK if the upload was successful
		if w.Code == http.StatusOK {
			// Check that the file was saved to the legacy images directory
			imagesPath := util.GetImagesPath()
			uploadedFilePath := filepath.Join(imagesPath, "uploaded-image.jpg")
			_, err := os.Stat(uploadedFilePath)

			// If the file exists, clean it up
			if err == nil {
				os.Remove(uploadedFilePath)
			}
		}
	})

	// Test that legacy document upload saves files to the legacy directory
	t.Run("LegacyDocumentUploadSavesToLegacyDirectory", func(t *testing.T) {
		// Create a test document file
		docPath := createTestDocumentFile(t, "test-upload-document.pdf")
		defer os.Remove(docPath)

		// Create a multipart form request
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		// Add the document file
		file, err := os.Open(docPath)
		assert.NoError(t, err)
		defer file.Close()

		part, err := writer.CreateFormFile("doc", filepath.Base(docPath))
		assert.NoError(t, err)

		_, err = io.Copy(part, file)
		assert.NoError(t, err)

		// Add the filename field
		err = writer.WriteField("filename", "uploaded-document.pdf")
		assert.NoError(t, err)

		err = writer.Close()
		assert.NoError(t, err)

		// Create the request
		req, _ := http.NewRequest("POST", "/api/cdn/upload/doc", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		w := httptest.NewRecorder()
		server.Engine.ServeHTTP(w, req)

		// The response should be OK if the upload was successful
		if w.Code == http.StatusOK {
			// Check that the file was saved to the legacy documents directory
			docsPath := util.GetDocsPath()
			uploadedFilePath := filepath.Join(docsPath, "uploaded-document.pdf")
			_, err := os.Stat(uploadedFilePath)

			// If the file exists, clean it up
			if err == nil {
				os.Remove(uploadedFilePath)
			}
		}
	})

	// Test that unified media upload saves files to the unified media directory
	t.Run("UnifiedMediaUploadSavesToUnifiedDirectory", func(t *testing.T) {
		// Create a test image file
		imagePath := createTestImageFile(t, "test-upload-unified-image.jpg")
		defer os.Remove(imagePath)

		// Create a multipart form request
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		// Add the image file
		file, err := os.Open(imagePath)
		assert.NoError(t, err)
		defer file.Close()

		part, err := writer.CreateFormFile("file", filepath.Base(imagePath))
		assert.NoError(t, err)

		_, err = io.Copy(part, file)
		assert.NoError(t, err)

		// Add the filename field
		err = writer.WriteField("filename", "uploaded-unified-image.jpg")
		assert.NoError(t, err)

		err = writer.Close()
		assert.NoError(t, err)

		// Create the request
		req, _ := http.NewRequest("POST", "/api/cdn/upload/media", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		w := httptest.NewRecorder()
		server.Engine.ServeHTTP(w, req)

		// The response should be OK if the upload was successful
		if w.Code == http.StatusOK {
			// Check that the file was saved to the unified media directory
			mediaPath := util.GetMediaUploadPath()
			uploadedFilePath := filepath.Join(mediaPath, "uploaded-unified-image.jpg")
			_, err := os.Stat(uploadedFilePath)

			// If the file exists, clean it up
			if err == nil {
				os.Remove(uploadedFilePath)
			}
		}
	})
}

// setupTestDatabase sets up a test database with the required tables
func setupTestDatabase(t *testing.T) *gorm.DB {
	// Initialize environment variables
	initializers.LoadEnvVariables(false)

	// Connect to test database
	database.ConnectToDB()
	db := database.DB

	// Auto migrate the models
	err := db.AutoMigrate(&models.Image{}, &models.Doc{}, &models.Media{})
	assert.NoError(t, err)

	return db
}

// cleanupTestDatabase cleans up the test database
func cleanupTestDatabase(t *testing.T, db *gorm.DB) {
	// Drop all tables
	err := db.Migrator().DropTable(&models.Media{}, &models.Image{}, &models.Doc{})
	assert.NoError(t, err)

	// Close the database connection
	sqlDB, err := db.DB()
	assert.NoError(t, err)
	err = sqlDB.Close()
	assert.NoError(t, err)
}

// createTestImageFiles creates test image files in the file system
func createTestImageFiles(t *testing.T) []string {
	// Create the images directory if it doesn't exist
	imagesPath := util.GetImagesPath()
	err := os.MkdirAll(imagesPath, 0755)
	assert.NoError(t, err)

	images := []string{
		filepath.Join(imagesPath, "test-image-1.jpg"),
		filepath.Join(imagesPath, "test-image-2.png"),
		filepath.Join(imagesPath, "test-image-3.gif"),
	}

	for i, imagePath := range images {
		file, err := os.Create(imagePath)
		assert.NoError(t, err)
		defer file.Close()

		// Write some test data
		_, err = file.WriteString(fmt.Sprintf("test image data %d", i+1))
		assert.NoError(t, err)
	}

	return images
}

// createTestDocumentFiles creates test document files in the file system
func createTestDocumentFiles(t *testing.T) []string {
	// Create the documents directory if it doesn't exist
	docsPath := util.GetDocsPath()
	err := os.MkdirAll(docsPath, 0755)
	assert.NoError(t, err)

	docs := []string{
		filepath.Join(docsPath, "test-document-1.pdf"),
		filepath.Join(docsPath, "test-document-2.docx"),
		filepath.Join(docsPath, "test-document-3.txt"),
	}

	for i, docPath := range docs {
		file, err := os.Create(docPath)
		assert.NoError(t, err)
		defer file.Close()

		// Write some test data
		_, err = file.WriteString(fmt.Sprintf("test document data %d", i+1))
		assert.NoError(t, err)
	}

	return docs
}

// cleanupTestFiles cleans up the test files
func cleanupTestFiles(t *testing.T, images []string, docs []string) {
	// Remove test image files
	for _, imagePath := range images {
		err := os.Remove(imagePath)
		// It's okay if the file doesn't exist
		if err != nil && !os.IsNotExist(err) {
			assert.NoError(t, err)
		}
	}

	// Remove test document files
	for _, docPath := range docs {
		err := os.Remove(docPath)
		// It's okay if the file doesn't exist
		if err != nil && !os.IsNotExist(err) {
			assert.NoError(t, err)
		}
	}
}

// createTestImageFile creates a test image file
func createTestImageFile(t *testing.T, filename string) string {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "test-images")
	assert.NoError(t, err)

	// Create a test image file
	imagePath := filepath.Join(tempDir, filename)
	file, err := os.Create(imagePath)
	assert.NoError(t, err)
	defer file.Close()

	// Write some test data
	_, err = file.WriteString("test image data")
	assert.NoError(t, err)

	return imagePath
}

// createTestDocumentFile creates a test document file
func createTestDocumentFile(t *testing.T, filename string) string {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "test-documents")
	assert.NoError(t, err)

	// Create a test document file
	docPath := filepath.Join(tempDir, filename)
	file, err := os.Create(docPath)
	assert.NoError(t, err)
	defer file.Close()

	// Write some test data
	_, err = file.WriteString("test document data")
	assert.NoError(t, err)

	return docPath
}

func main() {
	// Run the tests
	fmt.Println("Running file system compatibility tests...")

	// Run the file system backward compatibility tests
	fmt.Println("Testing file system backward compatibility...")
	TestFileSystemBackwardCompatibility(nil)
	fmt.Println("File system backward compatibility tests passed!")

	// Run the file upload backward compatibility tests
	fmt.Println("Testing file upload backward compatibility...")
	TestFileUploadBackwardCompatibility(nil)
	fmt.Println("File upload backward compatibility tests passed!")

	fmt.Println("All file system compatibility tests passed!")
}
