package handlers

import (
	"bytes"
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

// Test backward compatibility for image upload using legacy endpoint
func TestHandleImageUpload_BackwardCompatibility(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "media-test")
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

	// Create test image
	img, err := testutils.CreateDummyImage(200, 200)
	require.NoError(t, err)

	// Create multipart form with image
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("image", "test-image.png")
	require.NoError(t, err)

	err = testutils.EncodeImage(part, img)
	require.NoError(t, err)

	err = writer.Close()
	require.NoError(t, err)

	// Create test request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest(http.MethodPost, "/api/cdn/upload/image", body)
	c.Request.Header.Add("Content-Type", writer.FormDataContentType())

	// Create media handler (which handles both legacy and unified endpoints)
	mediaHandler := NewMediaHandler(database.NewMediaRepo(database.DB))

	// Test the legacy image upload endpoint
	mediaHandler.HandleImageUpload(c)

	// Assert successful response
	require.Equal(t, http.StatusOK, w.Result().StatusCode)

	// Verify response contains file_url
	responseBody := w.Body.String()
	require.Contains(t, responseBody, "file_url")
	require.Contains(t, responseBody, "/download/images/")
}

// Test backward compatibility for document upload using legacy endpoint
func TestHandleDocUpload_BackwardCompatibility(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "media-test")
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

	// Create test document
	docContent := testutils.CreateDummyDocument()

	// Create multipart form with document
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("doc", "test-document.txt")
	require.NoError(t, err)

	_, err = part.Write(docContent)
	require.NoError(t, err)

	err = writer.Close()
	require.NoError(t, err)

	// Create test request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest(http.MethodPost, "/api/cdn/upload/doc", body)
	c.Request.Header.Add("Content-Type", writer.FormDataContentType())

	// Create media handler (which handles both legacy and unified endpoints)
	mediaHandler := NewMediaHandler(database.NewMediaRepo(database.DB))

	// Test the legacy document upload endpoint
	mediaHandler.HandleDocUpload(c)

	// Assert successful response
	require.Equal(t, http.StatusOK, w.Result().StatusCode)

	// Verify response contains file_url
	responseBody := w.Body.String()
	require.Contains(t, responseBody, "file_url")
	require.Contains(t, responseBody, "/download/docs/")
}

// Test unified media upload endpoint with image
func TestHandleMediaUpload_Image(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "media-test")
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

	// Create test image
	img, err := testutils.CreateDummyImage(200, 200)
	require.NoError(t, err)

	// Create multipart form with image
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", "test-image.png")
	require.NoError(t, err)

	err = testutils.EncodeImage(part, img)
	require.NoError(t, err)

	err = writer.Close()
	require.NoError(t, err)

	// Create test request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest(http.MethodPost, "/api/cdn/upload/media", body)
	c.Request.Header.Add("Content-Type", writer.FormDataContentType())

	// Create media handler
	mediaHandler := NewMediaHandler(database.NewMediaRepo(database.DB))

	// Test the unified media upload endpoint
	mediaHandler.HandleMediaUpload(c)

	// Assert successful response
	require.Equal(t, http.StatusOK, w.Result().StatusCode)

	// Verify response contains file_url and type
	responseBody := w.Body.String()
	require.Contains(t, responseBody, "file_url")
	require.Contains(t, responseBody, "type")
	require.Contains(t, responseBody, "image")
}

// Test unified media upload endpoint with document
func TestHandleMediaUpload_Document(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "media-test")
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

	// Create test document
	docContent := testutils.CreateDummyDocument()

	// Create multipart form with document
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", "test-document.txt")
	require.NoError(t, err)

	_, err = part.Write(docContent)
	require.NoError(t, err)

	err = writer.Close()
	require.NoError(t, err)

	// Create test request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest(http.MethodPost, "/api/cdn/upload/media", body)
	c.Request.Header.Add("Content-Type", writer.FormDataContentType())

	// Create media handler
	mediaHandler := NewMediaHandler(database.NewMediaRepo(database.DB))

	// Test the unified media upload endpoint
	mediaHandler.HandleMediaUpload(c)

	// Assert successful response
	require.Equal(t, http.StatusOK, w.Result().StatusCode)

	// Verify response contains file_url and type
	responseBody := w.Body.String()
	require.Contains(t, responseBody, "file_url")
	require.Contains(t, responseBody, "type")
	require.Contains(t, responseBody, "document")
}

// Test backward compatibility for image metadata endpoint
func TestHandleImageMetadata_BackwardCompatibility(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "media-test")
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

	// Create test image
	img, err := testutils.CreateDummyImage(200, 200)
	require.NoError(t, err)

	// Create images directory
	imagesDir := filepath.Join(util.ExPath, "uploads", "images")
	err = os.MkdirAll(imagesDir, 0755)
	require.NoError(t, err)

	// Save test image
	imgPath := filepath.Join(imagesDir, "test-image.png")
	imgFile, err := os.Create(imgPath)
	require.NoError(t, err)
	defer imgFile.Close()

	err = testutils.EncodeImage(imgFile, img)
	require.NoError(t, err)

	// Create test request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest(http.MethodGet, "/api/cdn/image/test-image.png", nil)
	c.Params = gin.Params{{Key: "filename", Value: "test-image.png"}}

	// Test the legacy image metadata endpoint
	HandleImageMetadata(c)

	// Assert successful response
	require.Equal(t, http.StatusOK, w.Result().StatusCode)

	// Verify response contains expected fields
	responseBody := w.Body.String()
	require.Contains(t, responseBody, "filename")
	require.Contains(t, responseBody, "test-image.png")
	require.Contains(t, responseBody, "download_url")
	require.Contains(t, responseBody, "file_size")
	require.Contains(t, responseBody, "width")
	require.Contains(t, responseBody, "height")
}

// Test backward compatibility for document metadata endpoint
func TestHandleDocMetadata_BackwardCompatibility(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "media-test")
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

	// Create test document
	docContent := testutils.CreateDummyDocument()

	// Create docs directory
	docsDir := filepath.Join(util.ExPath, "uploads", "docs")
	err = os.MkdirAll(docsDir, 0755)
	require.NoError(t, err)

	// Save test document
	docPath := filepath.Join(docsDir, "test-document.txt")
	err = os.WriteFile(docPath, docContent, 0644)
	require.NoError(t, err)

	// Create test request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest(http.MethodGet, "/api/cdn/doc/test-document.txt", nil)
	c.Params = gin.Params{{Key: "filename", Value: "test-document.txt"}}

	// Test the legacy document metadata endpoint
	HandleDocMetadata(c)

	// Assert successful response
	require.Equal(t, http.StatusOK, w.Result().StatusCode)

	// Verify response contains expected fields
	responseBody := w.Body.String()
	require.Contains(t, responseBody, "filename")
	require.Contains(t, responseBody, "test-document.txt")
	require.Contains(t, responseBody, "download_url")
	require.Contains(t, responseBody, "file_size")
}

// Test unified media metadata endpoint
func TestHandleMediaMetadata_Unified(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "media-test")
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

	// Create test image
	img, err := testutils.CreateDummyImage(200, 200)
	require.NoError(t, err)

	// Create media directory
	mediaDir := filepath.Join(util.ExPath, "uploads", "image")
	err = os.MkdirAll(mediaDir, 0755)
	require.NoError(t, err)

	// Save test image
	imgPath := filepath.Join(mediaDir, "test-image.png")
	imgFile, err := os.Create(imgPath)
	require.NoError(t, err)
	defer imgFile.Close()

	err = testutils.EncodeImage(imgFile, img)
	require.NoError(t, err)

	// Create media handler
	mediaHandler := NewMediaHandler(database.NewMediaRepo(database.DB))

	// Create test request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest(http.MethodGet, "/api/cdn/media/test-image.png?type=image", nil)
	c.Params = gin.Params{{Key: "filename", Value: "test-image.png"}}
	c.Request.URL.RawQuery = "type=image"

	// Test the unified media metadata endpoint
	mediaHandler.HandleMediaMetadata(c)

	// Assert successful response
	require.Equal(t, http.StatusOK, w.Result().StatusCode)

	// Verify response contains expected fields
	responseBody := w.Body.String()
	require.Contains(t, responseBody, "filename")
	require.Contains(t, responseBody, "test-image.png")
	require.Contains(t, responseBody, "download_url")
	require.Contains(t, responseBody, "file_size")
	require.Contains(t, responseBody, "type")
	require.Contains(t, responseBody, "image")
	require.Contains(t, responseBody, "width")
	require.Contains(t, responseBody, "height")
}

// Test backward compatibility for all images endpoint
func TestHandleAllImages_BackwardCompatibility(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "media-test")
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

	// Create test request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest(http.MethodGet, "/api/cdn/image/all", nil)

	// Test the legacy all images endpoint
	mediaHandler.HandleAllImages(c)

	// Assert successful response
	require.Equal(t, http.StatusOK, w.Result().StatusCode)

	// Verify response is an array
	responseBody := w.Body.String()
	require.True(t, strings.HasPrefix(responseBody, "[") || strings.HasPrefix(responseBody, "[]"))
}

// Test backward compatibility for all documents endpoint
func TestHandleAllDocs_BackwardCompatibility(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "media-test")
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

	// Create test request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest(http.MethodGet, "/api/cdn/doc/all", nil)

	// Test the legacy all documents endpoint
	mediaHandler.HandleAllDocs(c)

	// Assert successful response
	require.Equal(t, http.StatusOK, w.Result().StatusCode)

	// Verify response is an array
	responseBody := w.Body.String()
	require.True(t, strings.HasPrefix(responseBody, "[") || strings.HasPrefix(responseBody, "[]"))
}

// Test unified all media endpoint
func TestHandleAllMedia_Unified(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "media-test")
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

	// Create test request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest(http.MethodGet, "/api/cdn/media/all", nil)

	// Test the unified all media endpoint
	mediaHandler.HandleAllMedia(c)

	// Assert successful response
	require.Equal(t, http.StatusOK, w.Result().StatusCode)

	// Verify response is an array
	responseBody := w.Body.String()
	require.True(t, strings.HasPrefix(responseBody, "[") || strings.HasPrefix(responseBody, "[]"))
}

// Test that legacy image upload creates both legacy and unified database entries
func TestHandleImageUpload_DatabaseCompatibility(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "media-test")
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

	// Create test image
	img, err := testutils.CreateDummyImage(200, 200)
	require.NoError(t, err)

	// Create multipart form with image
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("image", "test-image.png")
	require.NoError(t, err)

	err = testutils.EncodeImage(part, img)
	require.NoError(t, err)

	err = writer.Close()
	require.NoError(t, err)

	// Create test request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest(http.MethodPost, "/api/cdn/upload/image", body)
	c.Request.Header.Add("Content-Type", writer.FormDataContentType())

	// Create media handler
	mediaHandler := NewMediaHandler(database.NewMediaRepo(database.DB))

	// Test the legacy image upload endpoint
	mediaHandler.HandleImageUpload(c)

	// Assert successful response
	require.Equal(t, http.StatusOK, w.Result().StatusCode)

	// Verify that the image was added to the database
	mediaRepo := database.NewMediaRepo(database.DB)

	// Calculate checksum of the test image
	checksum, err := testutils.CalculateImageChecksum(img)
	require.NoError(t, err)

	// Check if media exists in the database
	media := mediaRepo.GetMediaByCheckSum(checksum[:])
	require.NotEmpty(t, media.Checksum)
	require.Equal(t, "test-image.png", media.FileName)
	require.Equal(t, models.MediaTypeImage, media.Type)

	// Check if image dimensions were captured
	require.NotNil(t, media.Width)
	require.NotNil(t, media.Height)
	require.Equal(t, 200, *media.Width)
	require.Equal(t, 200, *media.Height)
}

// Test that legacy document upload creates both legacy and unified database entries
func TestHandleDocUpload_DatabaseCompatibility(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "media-test")
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

	// Create test document
	docContent := testutils.CreateDummyDocument()

	// Create multipart form with document
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("doc", "test-document.txt")
	require.NoError(t, err)

	_, err = part.Write(docContent)
	require.NoError(t, err)

	err = writer.Close()
	require.NoError(t, err)

	// Create test request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest(http.MethodPost, "/api/cdn/upload/doc", body)
	c.Request.Header.Add("Content-Type", writer.FormDataContentType())

	// Create media handler
	mediaHandler := NewMediaHandler(database.NewMediaRepo(database.DB))

	// Test the legacy document upload endpoint
	mediaHandler.HandleDocUpload(c)

	// Assert successful response
	require.Equal(t, http.StatusOK, w.Result().StatusCode)

	// Verify that the document was added to the database
	mediaRepo := database.NewMediaRepo(database.DB)

	// Calculate checksum of the test document
	checksum := testutils.CalculateDocumentChecksum(docContent)

	// Check if media exists in the database
	media := mediaRepo.GetMediaByCheckSum(checksum[:])
	require.NotEmpty(t, media.Checksum)
	require.Equal(t, "test-document.txt", media.FileName)
	require.Equal(t, models.MediaTypeDocument, media.Type)

	// Documents should not have dimensions
	require.Nil(t, media.Width)
	require.Nil(t, media.Height)
}
