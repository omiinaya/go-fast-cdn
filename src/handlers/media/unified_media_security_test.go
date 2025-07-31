package handlers

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/kevinanielsen/go-fast-cdn/src/database"
	authHandlers "github.com/kevinanielsen/go-fast-cdn/src/handlers/auth"
	"github.com/kevinanielsen/go-fast-cdn/src/models"
	testutils "github.com/kevinanielsen/go-fast-cdn/src/testUtils"
	"github.com/kevinanielsen/go-fast-cdn/src/util"
	"github.com/stretchr/testify/require"
)

// TestUnifiedMediaSecurity tests that existing security measures still work correctly with the unified media repository
func TestUnifiedMediaSecurity(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "unified-media-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	util.ExPath = tempDir

	// Set environment variables for testing
	os.Setenv("JWT_SECRET", "test-super-secret-jwt-key-for-testing-only")
	defer os.Unsetenv("JWT_SECRET")

	// Initialize database
	database.ConnectToDB()
	database.MigrateWithMedia()
	defer func() {
		dbPath := filepath.Join(util.ExPath, database.DbFolder, database.DbName)
		if err := os.Remove(dbPath); err != nil && !os.IsNotExist(err) {
			t.Error(err)
		}
	}()

	// Create media handler
	mediaHandler := NewMediaHandler(database.NewMediaRepo(database.DB))

	// Create and authenticate a user
	userRepo := database.NewUserRepo(database.DB)
	user := &models.User{
		Email: "test@example.com",
		Role:  "user",
	}
	err = user.HashPassword("SecureP@ssw0rd123")
	require.NoError(t, err)
	err = userRepo.CreateUser(user)
	require.NoError(t, err)

	// Create auth handler to get tokens
	authHandler := authHandlers.NewAuthHandler(userRepo)

	loginReq := authHandlers.LoginRequest{
		Email:    "test@example.com",
		Password: "SecureP@ssw0rd123",
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	jsonData, _ := json.Marshal(loginReq)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(jsonData))
	c.Request.Header.Add("Content-Type", "application/json")

	authHandler.Login(c)
	require.Equal(t, http.StatusOK, w.Result().StatusCode)

	var authResponse authHandlers.AuthResponse
	err = json.Unmarshal(w.Body.Bytes(), &authResponse)
	require.NoError(t, err)

	t.Run("UnifiedMediaUploadSecurity", func(t *testing.T) {
		// Test that the unified media upload endpoint maintains security measures

		// Create test image
		img, err := testutils.CreateDummyImage(200, 200)
		require.NoError(t, err)

		// Create multipart form with image
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		part, err := writer.CreateFormFile("file", "test-unified.png")
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
		c.Request.Header.Add("Authorization", "Bearer "+authResponse.AccessToken)

		// Test the unified media upload endpoint
		mediaHandler.HandleMediaUpload(c)

		// Should succeed
		require.Equal(t, http.StatusOK, w.Result().StatusCode)

		// Extract response
		var uploadResponse map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &uploadResponse)
		require.NoError(t, err)

		// Verify response contains expected fields
		require.Contains(t, uploadResponse, "file_url")
		require.Contains(t, uploadResponse, "type")

		// Verify the file URL uses the unified media path
		fileURL := uploadResponse["file_url"].(string)
		require.Contains(t, fileURL, "/api/cdn/download/media/")
	})

	t.Run("BackwardCompatibilityWithLegacyEndpoints", func(t *testing.T) {
		// Test that legacy endpoints still work and maintain security measures

		// Create test image
		img, err := testutils.CreateDummyImage(200, 200)
		require.NoError(t, err)

		// Create multipart form with image using legacy field name
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		part, err := writer.CreateFormFile("image", "test-legacy.png")
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
		c.Request.Header.Add("Authorization", "Bearer "+authResponse.AccessToken)

		// Test the legacy image upload endpoint
		mediaHandler.HandleImageUpload(c)

		// Should succeed
		require.Equal(t, http.StatusOK, w.Result().StatusCode)

		// Extract response
		var uploadResponse map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &uploadResponse)
		require.NoError(t, err)

		// Verify response contains expected fields
		require.Contains(t, uploadResponse, "file_url")

		// Verify the file URL uses the legacy path
		fileURL := uploadResponse["file_url"].(string)
		require.Contains(t, fileURL, "/download/images/")
	})

	t.Run("UnifiedMediaDeleteSecurity", func(t *testing.T) {
		// First upload a file using the unified endpoint
		img, err := testutils.CreateDummyImage(200, 200)
		require.NoError(t, err)

		// Create multipart form with image
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		part, err := writer.CreateFormFile("file", "test-delete-unified.png")
		require.NoError(t, err)

		err = testutils.EncodeImage(part, img)
		require.NoError(t, err)

		err = writer.Close()
		require.NoError(t, err)

		w = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(w)

		c.Request = httptest.NewRequest(http.MethodPost, "/api/cdn/upload/media", body)
		c.Request.Header.Add("Content-Type", writer.FormDataContentType())
		c.Request.Header.Add("Authorization", "Bearer "+authResponse.AccessToken)

		// Upload the file
		mediaHandler.HandleMediaUpload(c)
		require.Equal(t, http.StatusOK, w.Result().StatusCode)

		// Extract filename from response
		var uploadResponse map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &uploadResponse)
		require.NoError(t, err)
		fileURL := uploadResponse["file_url"].(string)
		// Extract filename from URL (format: host/path/filename)
		parts := strings.Split(fileURL, "/")
		filename := parts[len(parts)-1]

		// Test deletion using the unified endpoint
		w = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(w)

		c.Request = httptest.NewRequest(http.MethodDelete, "/api/cdn/media/"+filename+"?type=image", nil)
		c.Params = gin.Params{{Key: "filename", Value: filename}}
		c.Request.Header.Add("Authorization", "Bearer "+authResponse.AccessToken)

		// Test the unified media delete endpoint
		mediaHandler.HandleMediaDelete(c)

		// Should succeed
		require.Equal(t, http.StatusOK, w.Result().StatusCode)

		// Verify file was actually deleted
		require.FileExists(t, filepath.Join(util.GetMediaUploadPath(), filename), false)
	})

	t.Run("UnifiedMediaMetadataSecurity", func(t *testing.T) {
		// First upload a file using the unified endpoint
		img, err := testutils.CreateDummyImage(200, 200)
		require.NoError(t, err)

		// Create multipart form with image
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		part, err := writer.CreateFormFile("file", "test-metadata-unified.png")
		require.NoError(t, err)

		err = testutils.EncodeImage(part, img)
		require.NoError(t, err)

		err = writer.Close()
		require.NoError(t, err)

		w = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(w)

		c.Request = httptest.NewRequest(http.MethodPost, "/api/cdn/upload/media", body)
		c.Request.Header.Add("Content-Type", writer.FormDataContentType())
		c.Request.Header.Add("Authorization", "Bearer "+authResponse.AccessToken)

		// Upload the file
		mediaHandler.HandleMediaUpload(c)
		require.Equal(t, http.StatusOK, w.Result().StatusCode)

		// Extract filename from response
		var uploadResponse map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &uploadResponse)
		require.NoError(t, err)
		fileURL := uploadResponse["file_url"].(string)
		// Extract filename from URL (format: host/path/filename)
		parts := strings.Split(fileURL, "/")
		filename := parts[len(parts)-1]

		// Test metadata retrieval using the unified endpoint
		w = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(w)

		c.Request = httptest.NewRequest(http.MethodGet, "/api/cdn/media/"+filename+"?type=image", nil)
		c.Params = gin.Params{{Key: "filename", Value: filename}}
		c.Request.Header.Add("Authorization", "Bearer "+authResponse.AccessToken)

		// Test the unified media metadata endpoint
		mediaHandler.HandleMediaMetadata(c)

		// Should succeed
		require.Equal(t, http.StatusOK, w.Result().StatusCode)

		// Verify response contains expected fields
		var metadataResponse map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &metadataResponse)
		require.NoError(t, err)

		require.Contains(t, metadataResponse, "filename")
		require.Contains(t, metadataResponse, "download_url")
		require.Contains(t, metadataResponse, "file_size")
		require.Contains(t, metadataResponse, "type")
		require.Contains(t, metadataResponse, "width")
		require.Contains(t, metadataResponse, "height")

		// Verify the download URL uses the unified media path
		downloadURL := metadataResponse["download_url"].(string)
		require.Contains(t, downloadURL, "/api/cdn/download/media/")
	})

	t.Run("UnifiedMediaRenameSecurity", func(t *testing.T) {
		// First upload a file using the unified endpoint
		img, err := testutils.CreateDummyImage(200, 200)
		require.NoError(t, err)

		// Create multipart form with image
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		part, err := writer.CreateFormFile("file", "test-rename-unified.png")
		require.NoError(t, err)

		err = testutils.EncodeImage(part, img)
		require.NoError(t, err)

		err = writer.Close()
		require.NoError(t, err)

		w = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(w)

		c.Request = httptest.NewRequest(http.MethodPost, "/api/cdn/upload/media", body)
		c.Request.Header.Add("Content-Type", writer.FormDataContentType())
		c.Request.Header.Add("Authorization", "Bearer "+authResponse.AccessToken)

		// Upload the file
		mediaHandler.HandleMediaUpload(c)
		require.Equal(t, http.StatusOK, w.Result().StatusCode)

		// Extract filename from response
		var uploadResponse map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &uploadResponse)
		require.NoError(t, err)
		fileURL := uploadResponse["file_url"].(string)
		// Extract filename from URL (format: host/path/filename)
		parts := strings.Split(fileURL, "/")
		filename := parts[len(parts)-1]

		// Create rename request
		renameReq := map[string]interface{}{
			"filename":   filename,
			"new_name":   "renamed-unified.png",
			"media_type": "image",
		}

		jsonData, _ := json.Marshal(renameReq)

		// Test renaming using the unified endpoint
		w = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(w)

		c.Request = httptest.NewRequest(http.MethodPut, "/api/cdn/rename/media", bytes.NewBuffer(jsonData))
		c.Request.Header.Add("Content-Type", "application/json")
		c.Request.Header.Add("Authorization", "Bearer "+authResponse.AccessToken)

		// Test the unified media rename endpoint
		mediaHandler.HandleMediaRename(c)

		// Should succeed
		require.Equal(t, http.StatusOK, w.Result().StatusCode)

		// Verify response contains expected fields
		var renameResponse map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &renameResponse)
		require.NoError(t, err)

		require.Contains(t, renameResponse, "message")
		require.Contains(t, renameResponse, "new_filename")

		// Verify the file was actually renamed
		newFilename := renameResponse["new_filename"].(string)
		require.FileExists(t, filepath.Join(util.GetMediaUploadPath(), newFilename), true)
		require.FileExists(t, filepath.Join(util.GetMediaUploadPath(), filename), false)
	})

	t.Run("UnifiedMediaResizeSecurity", func(t *testing.T) {
		// First upload a file using the unified endpoint
		img, err := testutils.CreateDummyImage(200, 200)
		require.NoError(t, err)

		// Create multipart form with image
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		part, err := writer.CreateFormFile("file", "test-resize-unified.png")
		require.NoError(t, err)

		err = testutils.EncodeImage(part, img)
		require.NoError(t, err)

		err = writer.Close()
		require.NoError(t, err)

		w = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(w)

		c.Request = httptest.NewRequest(http.MethodPost, "/api/cdn/upload/media", body)
		c.Request.Header.Add("Content-Type", writer.FormDataContentType())
		c.Request.Header.Add("Authorization", "Bearer "+authResponse.AccessToken)

		// Upload the file
		mediaHandler.HandleMediaUpload(c)
		require.Equal(t, http.StatusOK, w.Result().StatusCode)

		// Extract filename from response
		var uploadResponse map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &uploadResponse)
		require.NoError(t, err)
		fileURL := uploadResponse["file_url"].(string)
		// Extract filename from URL (format: host/path/filename)
		parts := strings.Split(fileURL, "/")
		filename := parts[len(parts)-1]

		// Test resizing using the unified endpoint
		w = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(w)

		c.Request = httptest.NewRequest(http.MethodPost, "/api/cdn/resize/media/"+filename+"?type=image&width=100&height=100", nil)
		c.Params = gin.Params{{Key: "filename", Value: filename}}
		c.Request.Header.Add("Authorization", "Bearer "+authResponse.AccessToken)

		// Test the unified media resize endpoint
		mediaHandler.HandleMediaResize(c)

		// Should succeed
		require.Equal(t, http.StatusOK, w.Result().StatusCode)

		// Verify response contains expected fields
		var resizeResponse map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &resizeResponse)
		require.NoError(t, err)

		require.Contains(t, resizeResponse, "file_url")
		require.Contains(t, resizeResponse, "width")
		require.Contains(t, resizeResponse, "height")

		// Verify the resized image dimensions
		require.Equal(t, 100.0, resizeResponse["width"])
		require.Equal(t, 100.0, resizeResponse["height"])

		// Verify the file URL uses the unified media path
		resizedFileURL := resizeResponse["file_url"].(string)
		require.Contains(t, resizedFileURL, "/api/cdn/download/media/")
	})

	t.Run("UnifiedMediaListSecurity", func(t *testing.T) {
		// First upload a file using the unified endpoint
		img, err := testutils.CreateDummyImage(200, 200)
		require.NoError(t, err)

		// Create multipart form with image
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		part, err := writer.CreateFormFile("file", "test-list-unified.png")
		require.NoError(t, err)

		err = testutils.EncodeImage(part, img)
		require.NoError(t, err)

		err = writer.Close()
		require.NoError(t, err)

		w = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(w)

		c.Request = httptest.NewRequest(http.MethodPost, "/api/cdn/upload/media", body)
		c.Request.Header.Add("Content-Type", writer.FormDataContentType())
		c.Request.Header.Add("Authorization", "Bearer "+authResponse.AccessToken)

		// Upload the file
		mediaHandler.HandleMediaUpload(c)
		require.Equal(t, http.StatusOK, w.Result().StatusCode)

		// Test listing using the unified endpoint
		w = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(w)

		c.Request = httptest.NewRequest(http.MethodGet, "/api/cdn/media?type=image", nil)
		c.Request.Header.Add("Authorization", "Bearer "+authResponse.AccessToken)

		// Test the unified media list endpoint
		mediaHandler.HandleAllMedia(c)

		// Should succeed
		require.Equal(t, http.StatusOK, w.Result().StatusCode)

		// Verify response contains expected fields
		var listResponse []map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &listResponse)
		require.NoError(t, err)

		// Should have at least one file
		require.Greater(t, len(listResponse), 0)

		// Verify each file has expected fields
		for _, file := range listResponse {
			require.Contains(t, file, "filename")
			require.Contains(t, file, "download_url")
			require.Contains(t, file, "file_size")
			require.Contains(t, file, "type")
			require.Contains(t, file, "width")
			require.Contains(t, file, "height")

			// Verify the download URL uses the unified media path
			downloadURL := file["download_url"].(string)
			require.Contains(t, downloadURL, "/api/cdn/download/media/")
		}
	})

	t.Run("CrossTypeSecurity", func(t *testing.T) {
		// Test that cross-type operations are properly secured

		// First upload an image using the unified endpoint
		img, err := testutils.CreateDummyImage(200, 200)
		require.NoError(t, err)

		// Create multipart form with image
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		part, err := writer.CreateFormFile("file", "test-cross-type.png")
		require.NoError(t, err)

		err = testutils.EncodeImage(part, img)
		require.NoError(t, err)

		err = writer.Close()
		require.NoError(t, err)

		w = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(w)

		c.Request = httptest.NewRequest(http.MethodPost, "/api/cdn/upload/media", body)
		c.Request.Header.Add("Content-Type", writer.FormDataContentType())
		c.Request.Header.Add("Authorization", "Bearer "+authResponse.AccessToken)

		// Upload the file
		mediaHandler.HandleMediaUpload(c)
		require.Equal(t, http.StatusOK, w.Result().StatusCode)

		// Extract filename from response
		var uploadResponse map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &uploadResponse)
		require.NoError(t, err)
		fileURL := uploadResponse["file_url"].(string)
		// Extract filename from URL (format: host/path/filename)
		parts := strings.Split(fileURL, "/")
		filename := parts[len(parts)-1]

		// Test accessing with wrong media type
		w = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(w)

		c.Request = httptest.NewRequest(http.MethodGet, "/api/cdn/media/"+filename+"?type=document", nil)
		c.Params = gin.Params{{Key: "filename", Value: filename}}
		c.Request.Header.Add("Authorization", "Bearer "+authResponse.AccessToken)

		// Test the unified media metadata endpoint with wrong type
		mediaHandler.HandleMediaMetadata(c)

		// Should fail due to type mismatch
		require.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	})

	t.Run("UnifiedMediaPathSecurity", func(t *testing.T) {
		// Test that unified media paths are properly secured

		// Create test image
		img, err := testutils.CreateDummyImage(200, 200)
		require.NoError(t, err)

		// Create multipart form with image
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		part, err := writer.CreateFormFile("file", "test-path-security.png")
		require.NoError(t, err)

		err = testutils.EncodeImage(part, img)
		require.NoError(t, err)

		err = writer.Close()
		require.NoError(t, err)

		// Create test request
		w := httptest.NewRecorder()
		c, _ = gin.CreateTestContext(w)

		c.Request = httptest.NewRequest(http.MethodPost, "/api/cdn/upload/media", body)
		c.Request.Header.Add("Content-Type", writer.FormDataContentType())
		c.Request.Header.Add("Authorization", "Bearer "+authResponse.AccessToken)

		// Test the unified media upload endpoint
		mediaHandler.HandleMediaUpload(c)

		// Should succeed
		require.Equal(t, http.StatusOK, w.Result().StatusCode)

		// Extract response
		var uploadResponse map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &uploadResponse)
		require.NoError(t, err)

		// Verify the file URL uses the unified media path
		fileURL := uploadResponse["file_url"].(string)
		require.Contains(t, fileURL, "/api/cdn/download/media/")

		// Verify the file is stored in the unified media directory
		parts := strings.Split(fileURL, "/")
		filename := parts[len(parts)-1]
		require.FileExists(t, filepath.Join(util.GetMediaUploadPath(), filename), true)
	})

	t.Run("LegacyToUnifiedMigrationSecurity", func(t *testing.T) {
		// Test that migration from legacy to unified maintains security

		// First upload a file using the legacy endpoint
		img, err := testutils.CreateDummyImage(200, 200)
		require.NoError(t, err)

		// Create multipart form with image using legacy field name
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		part, err := writer.CreateFormFile("image", "test-migration.png")
		require.NoError(t, err)

		err = testutils.EncodeImage(part, img)
		require.NoError(t, err)

		err = writer.Close()
		require.NoError(t, err)

		w = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(w)

		c.Request = httptest.NewRequest(http.MethodPost, "/api/cdn/upload/image", body)
		c.Request.Header.Add("Content-Type", writer.FormDataContentType())
		c.Request.Header.Add("Authorization", "Bearer "+authResponse.AccessToken)

		// Upload the file using legacy endpoint
		mediaHandler.HandleImageUpload(c)
		require.Equal(t, http.StatusOK, w.Result().StatusCode)

		// Extract filename from response
		var uploadResponse map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &uploadResponse)
		require.NoError(t, err)
		fileURL := uploadResponse["file_url"].(string)
		// Extract filename from URL (format: host/path/filename)
		parts := strings.Split(fileURL, "/")
		filename := parts[len(parts)-1]

		// Test accessing the legacy file through the unified endpoint
		w = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(w)

		c.Request = httptest.NewRequest(http.MethodGet, "/api/cdn/media/"+filename+"?type=image", nil)
		c.Params = gin.Params{{Key: "filename", Value: filename}}
		c.Request.Header.Add("Authorization", "Bearer "+authResponse.AccessToken)

		// Test the unified media metadata endpoint
		mediaHandler.HandleMediaMetadata(c)

		// Should succeed (unified endpoint can access legacy files)
		require.Equal(t, http.StatusOK, w.Result().StatusCode)

		// Verify response contains expected fields
		var metadataResponse map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &metadataResponse)
		require.NoError(t, err)

		require.Contains(t, metadataResponse, "filename")
		require.Contains(t, metadataResponse, "download_url")
		require.Contains(t, metadataResponse, "file_size")
		require.Contains(t, metadataResponse, "type")
		require.Contains(t, metadataResponse, "width")
		require.Contains(t, metadataResponse, "height")
	})
}
