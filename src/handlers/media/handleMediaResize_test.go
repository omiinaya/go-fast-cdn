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
	testutils "github.com/kevinanielsen/go-fast-cdn/src/testutils"
	"github.com/kevinanielsen/go-fast-cdn/src/util"
	"github.com/stretchr/testify/require"
)

// Test backward compatibility for image resize endpoint
func TestHandleImageResize_BackwardCompatibility(t *testing.T) {
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
	img, err := testutils.CreateDummyImage(400, 400)
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

	// First, upload the image
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest(http.MethodPost, "/api/cdn/upload/image", body)
	c.Request.Header.Add("Content-Type", writer.FormDataContentType())

	mediaHandler := NewMediaHandler(database.NewMediaRepo(database.DB))
	mediaHandler.HandleImageUpload(c)

	// Verify upload was successful
	require.Equal(t, http.StatusOK, w.Result().StatusCode)

	// Now test the resize endpoint
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	resizeBody := strings.NewReader(`{"filename": "test-image.png", "width": 200, "height": 200}`)
	c.Request = httptest.NewRequest(http.MethodPut, "/api/cdn/resize/image", resizeBody)
	c.Request.Header.Add("Content-Type", "application/json")

	// Test the legacy image resize endpoint
	HandleImageResize(c)

	// Assert successful response
	require.Equal(t, http.StatusOK, w.Result().StatusCode)

	// Verify response contains success message
	responseBody := w.Body.String()
	require.Contains(t, responseBody, "message")
	require.Contains(t, responseBody, "resized successfully")

	// Verify the resized file was created
	resizedPath := filepath.Join(util.ExPath, "uploads", "images", "resized_test-image.png")
	_, err = os.Stat(resizedPath)
	require.NoError(t, err)
}

// Test unified media resize endpoint with image
func TestHandleMediaResize_Image(t *testing.T) {
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
	img, err := testutils.CreateDummyImage(400, 400)
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

	// First, upload the image using unified endpoint
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest(http.MethodPost, "/api/cdn/upload/media", body)
	c.Request.Header.Add("Content-Type", writer.FormDataContentType())

	mediaHandler := NewMediaHandler(database.NewMediaRepo(database.DB))
	mediaHandler.HandleMediaUpload(c)

	// Verify upload was successful
	require.Equal(t, http.StatusOK, w.Result().StatusCode)

	// Now test the unified resize endpoint
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	resizeBody := strings.NewReader(`{"filename": "test-image.png", "width": 200, "height": 200, "type": "image"}`)
	c.Request = httptest.NewRequest(http.MethodPut, "/api/cdn/resize/media", resizeBody)
	c.Request.Header.Add("Content-Type", "application/json")

	// Test the unified media resize endpoint
	mediaHandler.HandleMediaResize(c)

	// Assert successful response
	require.Equal(t, http.StatusOK, w.Result().StatusCode)

	// Verify response contains success message
	responseBody := w.Body.String()
	require.Contains(t, responseBody, "message")
	require.Contains(t, responseBody, "resized successfully")

	// Verify the resized file was created
	resizedPath := filepath.Join(util.ExPath, "uploads", "media", "resized_test-image.png")
	_, err = os.Stat(resizedPath)
	require.NoError(t, err)
}

// Test resize non-existent image
func TestHandleImageResize_NotFound(t *testing.T) {
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
	_ = NewMediaHandler(database.NewMediaRepo(database.DB))

	// Create test request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	resizeBody := strings.NewReader(`{"filename": "non-existent.png", "width": 200, "height": 200}`)
	c.Request = httptest.NewRequest(http.MethodPut, "/api/cdn/resize/image", resizeBody)
	c.Request.Header.Add("Content-Type", "application/json")

	// Test the legacy image resize endpoint
	HandleImageResize(c)

	// Assert not found response
	require.Equal(t, http.StatusNotFound, w.Result().StatusCode)

	// Verify response contains error message
	responseBody := w.Body.String()
	require.Contains(t, responseBody, "error")
	require.Contains(t, responseBody, "not found")
}

// Test resize with invalid JSON
func TestHandleImageResize_InvalidJSON(t *testing.T) {
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
	_ = NewMediaHandler(database.NewMediaRepo(database.DB))

	// Create test request with invalid JSON
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	resizeBody := strings.NewReader(`{"invalid": "json"}`)
	c.Request = httptest.NewRequest(http.MethodPut, "/api/cdn/resize/image", resizeBody)
	c.Request.Header.Add("Content-Type", "application/json")

	// Test the legacy image resize endpoint
	HandleImageResize(c)

	// Assert bad request response
	require.Equal(t, http.StatusBadRequest, w.Result().StatusCode)

	// Verify response contains error message
	responseBody := w.Body.String()
	require.Contains(t, responseBody, "error")
}

// Test resize with invalid dimensions
func TestHandleImageResize_InvalidDimensions(t *testing.T) {
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
	_ = NewMediaHandler(database.NewMediaRepo(database.DB))

	// Create test request with invalid dimensions
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	resizeBody := strings.NewReader(`{"filename": "test-image.png", "width": -1, "height": 200}`)
	c.Request = httptest.NewRequest(http.MethodPut, "/api/cdn/resize/image", resizeBody)
	c.Request.Header.Add("Content-Type", "application/json")

	// Test the legacy image resize endpoint
	HandleImageResize(c)

	// Assert bad request response
	require.Equal(t, http.StatusBadRequest, w.Result().StatusCode)

	// Verify response contains error message
	responseBody := w.Body.String()
	require.Contains(t, responseBody, "error")
}

// Test resize document (should fail)
func TestHandleMediaResize_Document(t *testing.T) {
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

	// Create test request for document resize
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	resizeBody := strings.NewReader(`{"filename": "test-document.txt", "width": 200, "height": 200, "type": "document"}`)
	c.Request = httptest.NewRequest(http.MethodPut, "/api/cdn/resize/media", resizeBody)
	c.Request.Header.Add("Content-Type", "application/json")

	// Test the unified media resize endpoint
	mediaHandler.HandleMediaResize(c)

	// Assert bad request response
	require.Equal(t, http.StatusBadRequest, w.Result().StatusCode)

	// Verify response contains error message
	responseBody := w.Body.String()
	require.Contains(t, responseBody, "error")
	require.Contains(t, responseBody, "cannot be resized")
}
