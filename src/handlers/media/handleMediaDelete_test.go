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

// Test backward compatibility for image delete endpoint
func TestHandleImageDelete_BackwardCompatibility(t *testing.T) {
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

	// Now test the delete endpoint
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = httptest.NewRequest(http.MethodDelete, "/api/cdn/delete/image/test-image.png", nil)
	c.Params = gin.Params{{Key: "filename", Value: "test-image.png"}}

	// Test the legacy image delete endpoint
	mediaHandler.HandleImageDelete(c)

	// Assert successful response
	require.Equal(t, http.StatusOK, w.Result().StatusCode)

	// Verify response contains success message
	responseBody := w.Body.String()
	require.Contains(t, responseBody, "message")
	require.Contains(t, responseBody, "deleted successfully")
}

// Test backward compatibility for document delete endpoint
func TestHandleDocDelete_BackwardCompatibility(t *testing.T) {
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

	// Now test the delete endpoint
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = httptest.NewRequest(http.MethodDelete, "/api/cdn/delete/doc/test-document.txt", nil)
	c.Params = gin.Params{{Key: "filename", Value: "test-document.txt"}}

	// Test the legacy document delete endpoint
	mediaHandler.HandleDocDelete(c)

	// Assert successful response
	require.Equal(t, http.StatusOK, w.Result().StatusCode)

	// Verify response contains success message
	responseBody := w.Body.String()
	require.Contains(t, responseBody, "message")
	require.Contains(t, responseBody, "deleted successfully")
}

// Test unified media delete endpoint with image
func TestHandleMediaDelete_Image(t *testing.T) {
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

	// Now test the unified delete endpoint
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = httptest.NewRequest(http.MethodDelete, "/api/cdn/delete/media/test-image.png", nil)
	c.Params = gin.Params{{Key: "filename", Value: "test-image.png"}}
	c.Request.URL.RawQuery = "type=image"

	// Test the unified media delete endpoint
	mediaHandler.HandleMediaDelete(c)

	// Assert successful response
	require.Equal(t, http.StatusOK, w.Result().StatusCode)

	// Verify response contains success message
	responseBody := w.Body.String()
	require.Contains(t, responseBody, "message")
	require.Contains(t, responseBody, "deleted successfully")
}

// Test unified media delete endpoint with document
func TestHandleMediaDelete_Document(t *testing.T) {
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

	// Now test the unified delete endpoint
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = httptest.NewRequest(http.MethodDelete, "/api/cdn/delete/media/test-document.txt", nil)
	c.Params = gin.Params{{Key: "filename", Value: "test-document.txt"}}
	c.Request.URL.RawQuery = "type=document"

	// Test the unified media delete endpoint
	mediaHandler.HandleMediaDelete(c)

	// Assert successful response
	require.Equal(t, http.StatusOK, w.Result().StatusCode)

	// Verify response contains success message
	responseBody := w.Body.String()
	require.Contains(t, responseBody, "message")
	require.Contains(t, responseBody, "deleted successfully")
}

// Test delete non-existent image
func TestHandleImageDelete_NotFound(t *testing.T) {
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

	c.Request = httptest.NewRequest(http.MethodDelete, "/api/cdn/delete/image/non-existent.png", nil)
	c.Params = gin.Params{{Key: "filename", Value: "non-existent.png"}}

	// Test the legacy image delete endpoint
	mediaHandler.HandleImageDelete(c)

	// Assert not found response
	require.Equal(t, http.StatusNotFound, w.Result().StatusCode)

	// Verify response contains error message
	responseBody := w.Body.String()
	require.Contains(t, responseBody, "error")
	require.Contains(t, responseBody, "not found")
}

// Test delete non-existent document
func TestHandleDocDelete_NotFound(t *testing.T) {
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

	c.Request = httptest.NewRequest(http.MethodDelete, "/api/cdn/delete/doc/non-existent.txt", nil)
	c.Params = gin.Params{{Key: "filename", Value: "non-existent.txt"}}

	// Test the legacy document delete endpoint
	mediaHandler.HandleDocDelete(c)

	// Assert not found response
	require.Equal(t, http.StatusNotFound, w.Result().StatusCode)

	// Verify response contains error message
	responseBody := w.Body.String()
	require.Contains(t, responseBody, "error")
	require.Contains(t, responseBody, "not found")
}

// Test that legacy image delete removes both legacy and unified database entries
func TestHandleImageDelete_DatabaseCompatibility(t *testing.T) {
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

	// Calculate checksum of the test image
	checksum, err := testutils.CalculateImageChecksum(img)
	require.NoError(t, err)

	// Verify that the image exists in the database before deletion
	mediaRepo := database.NewMediaRepo(database.DB)
	media := mediaRepo.GetMediaByCheckSum(checksum[:])
	require.NotEmpty(t, media.Checksum)

	// Now test the delete endpoint
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = httptest.NewRequest(http.MethodDelete, "/api/cdn/delete/image/test-image.png", nil)
	c.Params = gin.Params{{Key: "filename", Value: "test-image.png"}}

	// Test the legacy image delete endpoint
	mediaHandler.HandleImageDelete(c)

	// Assert successful response
	require.Equal(t, http.StatusOK, w.Result().StatusCode)

	// Verify that the image was removed from the database
	media = mediaRepo.GetMediaByCheckSum(checksum[:])
	require.Empty(t, media.Checksum)
}

// Test that legacy document delete removes both legacy and unified database entries
func TestHandleDocDelete_DatabaseCompatibility(t *testing.T) {
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

	// Calculate checksum of the test document
	checksum := testutils.CalculateDocumentChecksum(docContent)

	// Verify that the document exists in the database before deletion
	mediaRepo := database.NewMediaRepo(database.DB)
	media := mediaRepo.GetMediaByCheckSum(checksum[:])
	require.NotEmpty(t, media.Checksum)

	// Now test the delete endpoint
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = httptest.NewRequest(http.MethodDelete, "/api/cdn/delete/doc/test-document.txt", nil)
	c.Params = gin.Params{{Key: "filename", Value: "test-document.txt"}}

	// Test the legacy document delete endpoint
	mediaHandler.HandleDocDelete(c)

	// Assert successful response
	require.Equal(t, http.StatusOK, w.Result().StatusCode)

	// Verify that the document was removed from the database
	media = mediaRepo.GetMediaByCheckSum(checksum[:])
	require.Empty(t, media.Checksum)
}
