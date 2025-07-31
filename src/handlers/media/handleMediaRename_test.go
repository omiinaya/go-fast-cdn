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

// Test backward compatibility for image rename endpoint
func TestHandleImageRename_BackwardCompatibility(t *testing.T) {
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

	// First, upload the image
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest(http.MethodPost, "/api/cdn/upload/image", body)
	c.Request.Header.Add("Content-Type", writer.FormDataContentType())

	mediaHandler := NewMediaHandler(database.NewMediaRepo(database.DB))
	mediaHandler.HandleImageUpload(c)

	// Verify upload was successful
	require.Equal(t, http.StatusOK, w.Result().StatusCode)

	// Now test the rename endpoint
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	renameBody := strings.NewReader(`{"old_filename": "test-image.png", "new_filename": "renamed-image.png"}`)
	c.Request = httptest.NewRequest(http.MethodPut, "/api/cdn/rename/image", renameBody)
	c.Request.Header.Add("Content-Type", "application/json")

	// Test the legacy image rename endpoint
	mediaHandler.HandleImageRename(c)

	// Assert successful response
	require.Equal(t, http.StatusOK, w.Result().StatusCode)

	// Verify response contains success message
	responseBody := w.Body.String()
	require.Contains(t, responseBody, "message")
	require.Contains(t, responseBody, "renamed successfully")

	// Verify the file was actually renamed
	renamedPath := filepath.Join(util.ExPath, "uploads", "images", "renamed-image.png")
	_, err = os.Stat(renamedPath)
	require.NoError(t, err)

	// Verify the old file no longer exists
	oldPath := filepath.Join(util.ExPath, "uploads", "images", "test-image.png")
	_, err = os.Stat(oldPath)
	require.True(t, os.IsNotExist(err))
}

// Test backward compatibility for document rename endpoint
func TestHandleDocRename_BackwardCompatibility(t *testing.T) {
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

	// First, upload the document
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest(http.MethodPost, "/api/cdn/upload/doc", body)
	c.Request.Header.Add("Content-Type", writer.FormDataContentType())

	mediaHandler := NewMediaHandler(database.NewMediaRepo(database.DB))
	mediaHandler.HandleDocUpload(c)

	// Verify upload was successful
	require.Equal(t, http.StatusOK, w.Result().StatusCode)

	// Now test the rename endpoint
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	renameBody := strings.NewReader(`{"old_filename": "test-document.txt", "new_filename": "renamed-document.txt"}`)
	c.Request = httptest.NewRequest(http.MethodPut, "/api/cdn/rename/doc", renameBody)
	c.Request.Header.Add("Content-Type", "application/json")

	// Test the legacy document rename endpoint
	mediaHandler.HandleDocsRename(c)

	// Assert successful response
	require.Equal(t, http.StatusOK, w.Result().StatusCode)

	// Verify response contains success message
	responseBody := w.Body.String()
	require.Contains(t, responseBody, "message")
	require.Contains(t, responseBody, "renamed successfully")

	// Verify the file was actually renamed
	renamedPath := filepath.Join(util.ExPath, "uploads", "docs", "renamed-document.txt")
	_, err = os.Stat(renamedPath)
	require.NoError(t, err)

	// Verify the old file no longer exists
	oldPath := filepath.Join(util.ExPath, "uploads", "docs", "test-document.txt")
	_, err = os.Stat(oldPath)
	require.True(t, os.IsNotExist(err))
}

// Test unified media rename endpoint with image
func TestHandleMediaRename_Image(t *testing.T) {
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

	// Now test the unified rename endpoint
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	renameBody := strings.NewReader(`{"old_filename": "test-image.png", "new_filename": "renamed-image.png", "type": "image"}`)
	c.Request = httptest.NewRequest(http.MethodPut, "/api/cdn/rename/media", renameBody)
	c.Request.Header.Add("Content-Type", "application/json")

	// Test the unified media rename endpoint
	mediaHandler.HandleMediaRename(c)

	// Assert successful response
	require.Equal(t, http.StatusOK, w.Result().StatusCode)

	// Verify response contains success message
	responseBody := w.Body.String()
	require.Contains(t, responseBody, "message")
	require.Contains(t, responseBody, "renamed successfully")

	// Verify the file was actually renamed
	renamedPath := filepath.Join(util.ExPath, "uploads", "media", "renamed-image.png")
	_, err = os.Stat(renamedPath)
	require.NoError(t, err)

	// Verify the old file no longer exists
	oldPath := filepath.Join(util.ExPath, "uploads", "media", "test-image.png")
	_, err = os.Stat(oldPath)
	require.True(t, os.IsNotExist(err))
}

// Test unified media rename endpoint with document
func TestHandleMediaRename_Document(t *testing.T) {
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

	// Now test the unified rename endpoint
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	renameBody := strings.NewReader(`{"old_filename": "test-document.txt", "new_filename": "renamed-document.txt", "type": "document"}`)
	c.Request = httptest.NewRequest(http.MethodPut, "/api/cdn/rename/media", renameBody)
	c.Request.Header.Add("Content-Type", "application/json")

	// Test the unified media rename endpoint
	mediaHandler.HandleMediaRename(c)

	// Assert successful response
	require.Equal(t, http.StatusOK, w.Result().StatusCode)

	// Verify response contains success message
	responseBody := w.Body.String()
	require.Contains(t, responseBody, "message")
	require.Contains(t, responseBody, "renamed successfully")

	// Verify the file was actually renamed
	renamedPath := filepath.Join(util.ExPath, "uploads", "media", "renamed-document.txt")
	_, err = os.Stat(renamedPath)
	require.NoError(t, err)

	// Verify the old file no longer exists
	oldPath := filepath.Join(util.ExPath, "uploads", "media", "test-document.txt")
	_, err = os.Stat(oldPath)
	require.True(t, os.IsNotExist(err))
}

// Test rename non-existent image
func TestHandleImageRename_NotFound(t *testing.T) {
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

	renameBody := strings.NewReader(`{"old_filename": "non-existent.png", "new_filename": "new-name.png"}`)
	c.Request = httptest.NewRequest(http.MethodPut, "/api/cdn/rename/image", renameBody)
	c.Request.Header.Add("Content-Type", "application/json")

	// Test the legacy image rename endpoint
	mediaHandler.HandleImageRename(c)

	// Assert not found response
	require.Equal(t, http.StatusNotFound, w.Result().StatusCode)

	// Verify response contains error message
	responseBody := w.Body.String()
	require.Contains(t, responseBody, "error")
	require.Contains(t, responseBody, "not found")
}

// Test rename non-existent document
func TestHandleDocRename_NotFound(t *testing.T) {
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

	renameBody := strings.NewReader(`{"old_filename": "non-existent.txt", "new_filename": "new-name.txt"}`)
	c.Request = httptest.NewRequest(http.MethodPut, "/api/cdn/rename/doc", renameBody)
	c.Request.Header.Add("Content-Type", "application/json")

	// Test the legacy document rename endpoint
	mediaHandler.HandleDocsRename(c)

	// Assert not found response
	require.Equal(t, http.StatusNotFound, w.Result().StatusCode)

	// Verify response contains error message
	responseBody := w.Body.String()
	require.Contains(t, responseBody, "error")
	require.Contains(t, responseBody, "not found")
}

// Test rename with invalid JSON
func TestHandleImageRename_InvalidJSON(t *testing.T) {
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

	// Create test request with invalid JSON
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	renameBody := strings.NewReader(`{"invalid": "json"}`)
	c.Request = httptest.NewRequest(http.MethodPut, "/api/cdn/rename/image", renameBody)
	c.Request.Header.Add("Content-Type", "application/json")

	// Test the legacy image rename endpoint
	mediaHandler.HandleImageRename(c)

	// Assert bad request response
	require.Equal(t, http.StatusBadRequest, w.Result().StatusCode)

	// Verify response contains error message
	responseBody := w.Body.String()
	require.Contains(t, responseBody, "error")
}
