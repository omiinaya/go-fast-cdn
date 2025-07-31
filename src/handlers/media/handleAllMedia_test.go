package handlers

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/kevinanielsen/go-fast-cdn/src/database"
	testutils "github.com/kevinanielsen/go-fast-cdn/src/testutils"
	"github.com/kevinanielsen/go-fast-cdn/src/util"
	"github.com/stretchr/testify/require"
)

// Test unified handleAllMedia endpoint with images
func TestHandleAllMedia_Images(t *testing.T) {
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

	// First, upload the image using unified endpoint
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest(http.MethodPost, "/api/cdn/upload/media", body)
	c.Request.Header.Add("Content-Type", writer.FormDataContentType())

	mediaHandler := NewMediaHandler(database.NewMediaRepo(database.DB))
	mediaHandler.HandleMediaUpload(c)

	// Verify upload was successful
	require.Equal(t, http.StatusOK, w.Result().StatusCode)

	// Now test the unified get all media endpoint
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = httptest.NewRequest(http.MethodGet, "/api/cdn/media/all", nil)

	// Test the unified get all media endpoint
	mediaHandler.HandleAllMedia(c)

	// Assert successful response
	require.Equal(t, http.StatusOK, w.Result().StatusCode)

	// Verify response contains the uploaded image
	responseBody := w.Body.String()
	require.Contains(t, responseBody, "test-image.png")
	require.Contains(t, responseBody, "type")
	require.Contains(t, responseBody, "image")
}

// Test unified handleAllMedia endpoint with documents
func TestHandleAllMedia_Documents(t *testing.T) {
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

	// First, upload the document using unified endpoint
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest(http.MethodPost, "/api/cdn/upload/media", body)
	c.Request.Header.Add("Content-Type", writer.FormDataContentType())

	mediaHandler := NewMediaHandler(database.NewMediaRepo(database.DB))
	mediaHandler.HandleMediaUpload(c)

	// Verify upload was successful
	require.Equal(t, http.StatusOK, w.Result().StatusCode)

	// Now test the unified get all media endpoint
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = httptest.NewRequest(http.MethodGet, "/api/cdn/media/all", nil)

	// Test the unified get all media endpoint
	mediaHandler.HandleAllMedia(c)

	// Assert successful response
	require.Equal(t, http.StatusOK, w.Result().StatusCode)

	// Verify response contains the uploaded document
	responseBody := w.Body.String()
	require.Contains(t, responseBody, "test-document.txt")
	require.Contains(t, responseBody, "type")
	require.Contains(t, responseBody, "document")
}

// Test handleAllImages with empty database
func TestHandleAllImages_Empty(t *testing.T) {
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

	// Test the legacy get all images endpoint
	mediaHandler.HandleAllImages(c)

	// Assert successful response
	require.Equal(t, http.StatusOK, w.Result().StatusCode)

	// Verify response is empty array
	responseBody := w.Body.String()
	require.Equal(t, "[]", responseBody)
}

// Test handleAllDocs with empty database
func TestHandleAllDocs_Empty(t *testing.T) {
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

	// Test the legacy get all docs endpoint
	mediaHandler.HandleAllDocs(c)

	// Assert successful response
	require.Equal(t, http.StatusOK, w.Result().StatusCode)

	// Verify response is empty array
	responseBody := w.Body.String()
	require.Equal(t, "[]", responseBody)
}

// Test handleAllMedia with empty database
func TestHandleAllMedia_Empty(t *testing.T) {
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

	// Test the unified get all media endpoint
	mediaHandler.HandleAllMedia(c)

	// Assert successful response
	require.Equal(t, http.StatusOK, w.Result().StatusCode)

	// Verify response is empty array
	responseBody := w.Body.String()
	require.Equal(t, "[]", responseBody)
}

// Test handleAllMedia with mixed content
func TestHandleAllMedia_Mixed(t *testing.T) {
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

	mediaHandler := NewMediaHandler(database.NewMediaRepo(database.DB))

	// Create and upload test image
	img, err := testutils.CreateDummyImage(200, 200)
	require.NoError(t, err)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", "test-image.png")
	require.NoError(t, err)

	err = testutils.EncodeImage(part, img)
	require.NoError(t, err)

	err = writer.Close()
	require.NoError(t, err)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest(http.MethodPost, "/api/cdn/upload/media", body)
	c.Request.Header.Add("Content-Type", writer.FormDataContentType())

	mediaHandler.HandleMediaUpload(c)
	require.Equal(t, http.StatusOK, w.Result().StatusCode)

	// Create and upload test document
	docContent := testutils.CreateDummyDocument()

	body = &bytes.Buffer{}
	writer = multipart.NewWriter(body)

	part, err = writer.CreateFormFile("file", "test-document.txt")
	require.NoError(t, err)

	_, err = part.Write(docContent)
	require.NoError(t, err)

	err = writer.Close()
	require.NoError(t, err)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = httptest.NewRequest(http.MethodPost, "/api/cdn/upload/media", body)
	c.Request.Header.Add("Content-Type", writer.FormDataContentType())

	mediaHandler.HandleMediaUpload(c)
	require.Equal(t, http.StatusOK, w.Result().StatusCode)

	// Now test the unified get all media endpoint
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = httptest.NewRequest(http.MethodGet, "/api/cdn/media/all", nil)

	// Test the unified get all media endpoint
	mediaHandler.HandleAllMedia(c)

	// Assert successful response
	require.Equal(t, http.StatusOK, w.Result().StatusCode)

	// Verify response contains both the image and document
	responseBody := w.Body.String()
	require.Contains(t, responseBody, "test-image.png")
	require.Contains(t, responseBody, "test-document.txt")
	require.Contains(t, responseBody, "type")
	require.Contains(t, responseBody, "image")
	require.Contains(t, responseBody, "document")
}
