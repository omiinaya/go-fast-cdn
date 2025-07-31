package main

import (
	"bytes"
	"encoding/json"
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
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// TestDatabaseBackwardCompatibility tests that existing image and document data
// in the database is still accessible after the unification
func TestDatabaseBackwardCompatibility(t *testing.T) {
	// Setup test database
	db := setupTestDatabase(t)
	defer cleanupTestDatabase(t, db)

	// Create test data
	testImages := createTestImages(t, db)
	testDocuments := createTestDocuments(t, db)

	// Test that legacy image data is accessible through the unified media repository
	t.Run("LegacyImageDataAccessibleThroughUnifiedRepository", func(t *testing.T) {
		mediaRepo := database.NewMediaRepo(db)

		// Get all images using the unified repository
		images := mediaRepo.GetAllImages()

		// Verify that all test images are returned
		assert.Equal(t, len(testImages), len(images))

		// Verify that each image has the correct data
		for _, image := range images {
			found := false
			for _, testImage := range testImages {
				if image.FileName == testImage.FileName {
					found = true
					assert.Equal(t, testImage.Checksum, image.Checksum)
					assert.Equal(t, models.MediaTypeImage, image.Type)
					break
				}
			}
			assert.True(t, found, "Image %s not found in test data", image.FileName)
		}
	})

	// Test that legacy document data is accessible through the unified media repository
	t.Run("LegacyDocumentDataAccessibleThroughUnifiedRepository", func(t *testing.T) {
		mediaRepo := database.NewMediaRepo(db)

		// Get all documents using the unified repository
		documents := mediaRepo.GetAllDocs()

		// Verify that all test documents are returned
		assert.Equal(t, len(testDocuments), len(documents))

		// Verify that each document has the correct data
		for _, doc := range documents {
			found := false
			for _, testDoc := range testDocuments {
				if doc.FileName == testDoc.FileName {
					found = true
					assert.Equal(t, testDoc.Checksum, doc.Checksum)
					assert.Equal(t, models.MediaTypeDocument, doc.Type)
					break
				}
			}
			assert.True(t, found, "Document %s not found in test data", doc.FileName)
		}
	})

	// Test that legacy image data can be accessed by checksum
	t.Run("LegacyImageDataAccessibleByChecksum", func(t *testing.T) {
		mediaRepo := database.NewMediaRepo(db)

		// Test each image
		for _, testImage := range testImages {
			image := mediaRepo.GetImageByCheckSum(testImage.Checksum)

			// Verify that the image was found
			assert.NotEmpty(t, image.Checksum)
			assert.Equal(t, testImage.FileName, image.FileName)
			assert.Equal(t, testImage.Checksum, image.Checksum)
			assert.Equal(t, models.MediaTypeImage, image.Type)
		}
	})

	// Test that legacy document data can be accessed by checksum
	t.Run("LegacyDocumentDataAccessibleByChecksum", func(t *testing.T) {
		mediaRepo := database.NewMediaRepo(db)

		// Test each document
		for _, testDoc := range testDocuments {
			doc := mediaRepo.GetDocByCheckSum(testDoc.Checksum)

			// Verify that the document was found
			assert.NotEmpty(t, doc.Checksum)
			assert.Equal(t, testDoc.FileName, doc.FileName)
			assert.Equal(t, testDoc.Checksum, doc.Checksum)
			assert.Equal(t, models.MediaTypeDocument, doc.Type)
		}
	})

	// Test that legacy image data can be accessed by filename
	t.Run("LegacyImageDataAccessibleByFilename", func(t *testing.T) {
		mediaRepo := database.NewMediaRepo(db)

		// Test each image
		for _, testImage := range testImages {
			media := mediaRepo.GetMediaByFileName(testImage.FileName)

			// Verify that the media was found
			assert.NotEmpty(t, media.Checksum)
			assert.Equal(t, testImage.FileName, media.FileName)
			assert.Equal(t, testImage.Checksum, media.Checksum)
			assert.Equal(t, models.MediaTypeImage, media.Type)
		}
	})

	// Test that legacy document data can be accessed by filename
	t.Run("LegacyDocumentDataAccessibleByFilename", func(t *testing.T) {
		mediaRepo := database.NewMediaRepo(db)

		// Test each document
		for _, testDoc := range testDocuments {
			media := mediaRepo.GetMediaByFileName(testDoc.FileName)

			// Verify that the media was found
			assert.NotEmpty(t, media.Checksum)
			assert.Equal(t, testDoc.FileName, media.FileName)
			assert.Equal(t, testDoc.Checksum, media.Checksum)
			assert.Equal(t, models.MediaTypeDocument, media.Type)
		}
	})

	// Test that legacy image data can be accessed by type
	t.Run("LegacyImageDataAccessibleByType", func(t *testing.T) {
		mediaRepo := database.NewMediaRepo(db)

		// Get all media of type image
		images := mediaRepo.GetMediaByType(models.MediaTypeImage)

		// Verify that all test images are returned
		assert.Equal(t, len(testImages), len(images))

		// Verify that each media has the correct type
		for _, image := range images {
			assert.Equal(t, models.MediaTypeImage, image.Type)
		}
	})

	// Test that legacy document data can be accessed by type
	t.Run("LegacyDocumentDataAccessibleByType", func(t *testing.T) {
		mediaRepo := database.NewMediaRepo(db)

		// Get all media of type document
		documents := mediaRepo.GetMediaByType(models.MediaTypeDocument)

		// Verify that all test documents are returned
		assert.Equal(t, len(testDocuments), len(documents))

		// Verify that each media has the correct type
		for _, doc := range documents {
			assert.Equal(t, models.MediaTypeDocument, doc.Type)
		}
	})

	// Test that legacy image data can be converted to Image model
	t.Run("LegacyImageDataCanBeConvertedToImageModel", func(t *testing.T) {
		mediaRepo := database.NewMediaRepo(db)

		// Get all images using the unified repository
		images := mediaRepo.GetAllImages()

		// Convert each media to Image model
		for _, media := range images {
			image := media.ToImage()

			// Verify that the conversion is correct
			assert.Equal(t, media.ID, image.ID)
			assert.Equal(t, media.FileName, image.FileName)
			assert.Equal(t, media.Checksum, image.Checksum)
			assert.Equal(t, media.CreatedAt, image.CreatedAt)
			assert.Equal(t, media.UpdatedAt, image.UpdatedAt)
			assert.Equal(t, media.DeletedAt, image.DeletedAt)
		}
	})

	// Test that legacy document data can be converted to Doc model
	t.Run("LegacyDocumentDataCanBeConvertedToDocModel", func(t *testing.T) {
		mediaRepo := database.NewMediaRepo(db)

		// Get all documents using the unified repository
		documents := mediaRepo.GetAllDocs()

		// Convert each media to Doc model
		for _, media := range documents {
			doc := media.ToDoc()

			// Verify that the conversion is correct
			assert.Equal(t, media.ID, doc.ID)
			assert.Equal(t, media.FileName, doc.FileName)
			assert.Equal(t, media.Checksum, doc.Checksum)
			assert.Equal(t, media.CreatedAt, doc.CreatedAt)
			assert.Equal(t, media.UpdatedAt, doc.UpdatedAt)
			assert.Equal(t, media.DeletedAt, doc.DeletedAt)
		}
	})
}

// TestAPIBackwardCompatibility tests that existing image and document API endpoints
// still work correctly after the unification
func TestAPIBackwardCompatibility(t *testing.T) {
	// Setup test database
	db := setupTestDatabase(t)
	defer cleanupTestDatabase(t, db)

	// Create test data
	createTestImages(t, db)
	createTestDocuments(t, db)

	// Setup test server
	gin.SetMode(gin.TestMode)
	server := router.NewServer()

	// Test that legacy image endpoints still work
	t.Run("LegacyImageEndpoints", func(t *testing.T) {
		// Test GET /api/cdn/image/all
		req, _ := http.NewRequest("GET", "/api/cdn/image/all", nil)
		w := httptest.NewRecorder()
		server.Engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var images []models.Image
		err := json.Unmarshal(w.Body.Bytes(), &images)
		assert.NoError(t, err)
		assert.Greater(t, len(images), 0)

		// Test GET /api/cdn/image/:filename
		if len(images) > 0 {
			req, _ := http.NewRequest("GET", "/api/cdn/image/"+images[0].FileName, nil)
			w := httptest.NewRecorder()
			server.Engine.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)

			var imageMetadata map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &imageMetadata)
			assert.NoError(t, err)
			assert.Equal(t, images[0].FileName, imageMetadata["filename"])
		}
	})

	// Test that legacy document endpoints still work
	t.Run("LegacyDocumentEndpoints", func(t *testing.T) {
		// Test GET /api/cdn/doc/all
		req, _ := http.NewRequest("GET", "/api/cdn/doc/all", nil)
		w := httptest.NewRecorder()
		server.Engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var docs []models.Doc
		err := json.Unmarshal(w.Body.Bytes(), &docs)
		assert.NoError(t, err)
		assert.Greater(t, len(docs), 0)

		// Test GET /api/cdn/doc/:filename
		if len(docs) > 0 {
			req, _ := http.NewRequest("GET", "/api/cdn/doc/"+docs[0].FileName, nil)
			w := httptest.NewRecorder()
			server.Engine.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)

			var docMetadata map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &docMetadata)
			assert.NoError(t, err)
			assert.Equal(t, docs[0].FileName, docMetadata["filename"])
		}
	})

	// Test that unified media endpoints work
	t.Run("UnifiedMediaEndpoints", func(t *testing.T) {
		// Test GET /api/cdn/media/all
		req, _ := http.NewRequest("GET", "/api/cdn/media/all", nil)
		w := httptest.NewRecorder()
		server.Engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var media []models.Media
		err := json.Unmarshal(w.Body.Bytes(), &media)
		assert.NoError(t, err)
		assert.Greater(t, len(media), 0)

		// Test GET /api/cdn/media/:filename
		if len(media) > 0 {
			req, _ := http.NewRequest("GET", "/api/cdn/media/"+media[0].FileName+"?type="+string(media[0].Type), nil)
			w := httptest.NewRecorder()
			server.Engine.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)

			var mediaMetadata map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &mediaMetadata)
			assert.NoError(t, err)
			assert.Equal(t, media[0].FileName, mediaMetadata["filename"])
			assert.Equal(t, string(media[0].Type), mediaMetadata["type"])
		}
	})

	// Test that legacy image upload still works
	t.Run("LegacyImageUpload", func(t *testing.T) {
		// Create a test image file
		imagePath := createTestImageFile(t, "test-image.jpg")
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

		// The response should be a conflict since the image already exists
		assert.Equal(t, http.StatusConflict, w.Code)
	})

	// Test that legacy document upload still works
	t.Run("LegacyDocumentUpload", func(t *testing.T) {
		// Create a test document file
		docPath := createTestDocumentFile(t, "test-document.pdf")
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

		// The response should be a conflict since the document already exists
		assert.Equal(t, http.StatusConflict, w.Code)
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

// createTestImages creates test image data in the database
func createTestImages(t *testing.T, db *gorm.DB) []models.Image {
	images := []models.Image{
		{
			FileName: "test-image-1.jpg",
			Checksum: []byte{1, 2, 3, 4, 5},
		},
		{
			FileName: "test-image-2.png",
			Checksum: []byte{6, 7, 8, 9, 10},
		},
		{
			FileName: "test-image-3.gif",
			Checksum: []byte{11, 12, 13, 14, 15},
		},
	}

	for i := range images {
		err := db.Create(&images[i]).Error
		assert.NoError(t, err)
	}

	return images
}

// createTestDocuments creates test document data in the database
func createTestDocuments(t *testing.T, db *gorm.DB) []models.Doc {
	docs := []models.Doc{
		{
			FileName: "test-document-1.pdf",
			Checksum: []byte{16, 17, 18, 19, 20},
		},
		{
			FileName: "test-document-2.docx",
			Checksum: []byte{21, 22, 23, 24, 25},
		},
		{
			FileName: "test-document-3.txt",
			Checksum: []byte{26, 27, 28, 29, 30},
		},
	}

	for i := range docs {
		err := db.Create(&docs[i]).Error
		assert.NoError(t, err)
	}

	return docs
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
	fmt.Println("Running backward compatibility tests...")

	// Run the database backward compatibility tests
	fmt.Println("Testing database backward compatibility...")
	TestDatabaseBackwardCompatibility(nil)
	fmt.Println("Database backward compatibility tests passed!")

	// Run the API backward compatibility tests
	fmt.Println("Testing API backward compatibility...")
	TestAPIBackwardCompatibility(nil)
	fmt.Println("API backward compatibility tests passed!")

	fmt.Println("All backward compatibility tests passed!")
}
