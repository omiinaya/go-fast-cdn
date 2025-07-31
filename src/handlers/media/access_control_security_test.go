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
	"github.com/kevinanielsen/go-fast-cdn/src/testutils"
	"github.com/kevinanielsen/go-fast-cdn/src/util"
	"github.com/stretchr/testify/require"
)

// TestAccessControlSecurity tests access control for different user roles
func TestAccessControlSecurity(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "access-control-test")
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

	// Create and authenticate users with different roles
	userRepo := database.NewUserRepo(database.DB)

	// Create admin user
	adminUser := &models.User{
		Email: "admin@example.com",
		Role:  "admin",
	}
	err = adminUser.HashPassword("SecureP@ssw0rd123")
	require.NoError(t, err)
	err = userRepo.CreateUser(adminUser)
	require.NoError(t, err)

	// Create regular user
	regularUser := &models.User{
		Email: "user@example.com",
		Role:  "user",
	}
	err = regularUser.HashPassword("SecureP@ssw0rd123")
	require.NoError(t, err)
	err = userRepo.CreateUser(regularUser)
	require.NoError(t, err)

	// Create auth handler to get tokens
	authHandler := authHandlers.NewAuthHandler(userRepo)

	// Get admin token
	adminLoginReq := authHandlers.LoginRequest{
		Email:    "admin@example.com",
		Password: "SecureP@ssw0rd123",
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	jsonData, _ := json.Marshal(adminLoginReq)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(jsonData))
	c.Request.Header.Add("Content-Type", "application/json")

	authHandler.Login(c)
	require.Equal(t, http.StatusOK, w.Result().StatusCode)

	var adminAuthResponse authHandlers.AuthResponse
	err = json.Unmarshal(w.Body.Bytes(), &adminAuthResponse)
	require.NoError(t, err)

	// Get regular user token
	userLoginReq := authHandlers.LoginRequest{
		Email:    "user@example.com",
		Password: "SecureP@ssw0rd123",
	}

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	jsonData, _ = json.Marshal(userLoginReq)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(jsonData))
	c.Request.Header.Add("Content-Type", "application/json")

	authHandler.Login(c)
	require.Equal(t, http.StatusOK, w.Result().StatusCode)

	var userAuthResponse authHandlers.AuthResponse
	err = json.Unmarshal(w.Body.Bytes(), &userAuthResponse)
	require.NoError(t, err)

	t.Run("MediaUploadAccessControl", func(t *testing.T) {
		// Create test image
		img, err := testutils.CreateDummyImage(200, 200)
		require.NoError(t, err)

		// Create multipart form with image
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		part, err := writer.CreateFormFile("file", "test.png")
		require.NoError(t, err)

		err = testutils.EncodeImage(part, img)
		require.NoError(t, err)

		err = writer.Close()
		require.NoError(t, err)

		testCases := []struct {
			name           string
			token          string
			expectedStatus int
		}{
			{"AdminUser", adminAuthResponse.AccessToken, http.StatusOK},
			{"RegularUser", userAuthResponse.AccessToken, http.StatusOK},
			{"NoToken", "", http.StatusUnauthorized},
			{"InvalidToken", "invalid-token", http.StatusUnauthorized},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Create test request
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				c.Request = httptest.NewRequest(http.MethodPost, "/api/cdn/upload/media", body)
				c.Request.Header.Add("Content-Type", writer.FormDataContentType())
				if tc.token != "" {
					c.Request.Header.Add("Authorization", "Bearer "+tc.token)
				}

				// Test the unified media upload endpoint
				mediaHandler.HandleMediaUpload(c)

				// Should enforce access control
				require.Equal(t, tc.expectedStatus, w.Result().StatusCode)
			})
		}
	})

	t.Run("MediaDeleteAccessControl", func(t *testing.T) {
		// First upload a file with admin user
		img, err := testutils.CreateDummyImage(200, 200)
		require.NoError(t, err)

		// Create multipart form with image
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		part, err := writer.CreateFormFile("file", "test-delete.png")
		require.NoError(t, err)

		err = testutils.EncodeImage(part, img)
		require.NoError(t, err)

		err = writer.Close()
		require.NoError(t, err)

		w = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(w)

		c.Request = httptest.NewRequest(http.MethodPost, "/api/cdn/upload/media", body)
		c.Request.Header.Add("Content-Type", writer.FormDataContentType())
		c.Request.Header.Add("Authorization", "Bearer "+adminAuthResponse.AccessToken)

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

		testCases := []struct {
			name           string
			token          string
			expectedStatus int
		}{
			{"AdminUser", adminAuthResponse.AccessToken, http.StatusOK},
			{"RegularUser", userAuthResponse.AccessToken, http.StatusOK},
			{"NoToken", "", http.StatusUnauthorized},
			{"InvalidToken", "invalid-token", http.StatusUnauthorized},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Create test request
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				c.Request = httptest.NewRequest(http.MethodDelete, "/api/cdn/media/"+filename+"?type=image", nil)
				c.Params = gin.Params{{Key: "filename", Value: filename}}
				if tc.token != "" {
					c.Request.Header.Add("Authorization", "Bearer "+tc.token)
				}

				// Test the unified media delete endpoint
				mediaHandler.HandleMediaDelete(c)

				// Should enforce access control
				require.Equal(t, tc.expectedStatus, w.Result().StatusCode)
			})
		}
	})

	t.Run("MediaMetadataAccessControl", func(t *testing.T) {
		// First upload a file with admin user
		img, err := testutils.CreateDummyImage(200, 200)
		require.NoError(t, err)

		// Create multipart form with image
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		part, err := writer.CreateFormFile("file", "test-metadata.png")
		require.NoError(t, err)

		err = testutils.EncodeImage(part, img)
		require.NoError(t, err)

		err = writer.Close()
		require.NoError(t, err)

		w = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(w)

		c.Request = httptest.NewRequest(http.MethodPost, "/api/cdn/upload/media", body)
		c.Request.Header.Add("Content-Type", writer.FormDataContentType())
		c.Request.Header.Add("Authorization", "Bearer "+adminAuthResponse.AccessToken)

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

		testCases := []struct {
			name           string
			token          string
			expectedStatus int
		}{
			{"AdminUser", adminAuthResponse.AccessToken, http.StatusOK},
			{"RegularUser", userAuthResponse.AccessToken, http.StatusOK},
			{"NoToken", "", http.StatusUnauthorized},
			{"InvalidToken", "invalid-token", http.StatusUnauthorized},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Create test request
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				c.Request = httptest.NewRequest(http.MethodGet, "/api/cdn/media/"+filename+"?type=image", nil)
				c.Params = gin.Params{{Key: "filename", Value: filename}}
				if tc.token != "" {
					c.Request.Header.Add("Authorization", "Bearer "+tc.token)
				}

				// Test the unified media metadata endpoint
				mediaHandler.HandleMediaMetadata(c)

				// Should enforce access control
				require.Equal(t, tc.expectedStatus, w.Result().StatusCode)
			})
		}
	})

	t.Run("MediaRenameAccessControl", func(t *testing.T) {
		// First upload a file with admin user
		img, err := testutils.CreateDummyImage(200, 200)
		require.NoError(t, err)

		// Create multipart form with image
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		part, err := writer.CreateFormFile("file", "test-rename.png")
		require.NoError(t, err)

		err = testutils.EncodeImage(part, img)
		require.NoError(t, err)

		err = writer.Close()
		require.NoError(t, err)

		w = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(w)

		c.Request = httptest.NewRequest(http.MethodPost, "/api/cdn/upload/media", body)
		c.Request.Header.Add("Content-Type", writer.FormDataContentType())
		c.Request.Header.Add("Authorization", "Bearer "+adminAuthResponse.AccessToken)

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
		_ = parts[len(parts)-1] // filename is not used in this test

		// Create rename request
		renameReq := map[string]string{
			"new_name": "renamed.png",
		}

		jsonData, _ := json.Marshal(renameReq)

		testCases := []struct {
			name           string
			token          string
			expectedStatus int
		}{
			{"AdminUser", adminAuthResponse.AccessToken, http.StatusOK},
			{"RegularUser", userAuthResponse.AccessToken, http.StatusOK},
			{"NoToken", "", http.StatusUnauthorized},
			{"InvalidToken", "invalid-token", http.StatusUnauthorized},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Create test request
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				c.Request = httptest.NewRequest(http.MethodPut, "/api/cdn/rename/media", bytes.NewBuffer(jsonData))
				c.Request.Header.Add("Content-Type", "application/json")
				if tc.token != "" {
					c.Request.Header.Add("Authorization", "Bearer "+tc.token)
				}

				// Test the unified media rename endpoint
				mediaHandler.HandleMediaRename(c)

				// Should enforce access control
				require.Equal(t, tc.expectedStatus, w.Result().StatusCode)
			})
		}
	})

	t.Run("MediaResizeAccessControl", func(t *testing.T) {
		// First upload a file with admin user
		img, err := testutils.CreateDummyImage(200, 200)
		require.NoError(t, err)

		// Create multipart form with image
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		part, err := writer.CreateFormFile("file", "test-resize.png")
		require.NoError(t, err)

		err = testutils.EncodeImage(part, img)
		require.NoError(t, err)

		err = writer.Close()
		require.NoError(t, err)

		w = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(w)

		c.Request = httptest.NewRequest(http.MethodPost, "/api/cdn/upload/media", body)
		c.Request.Header.Add("Content-Type", writer.FormDataContentType())
		c.Request.Header.Add("Authorization", "Bearer "+adminAuthResponse.AccessToken)

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

		testCases := []struct {
			name           string
			token          string
			expectedStatus int
		}{
			{"AdminUser", adminAuthResponse.AccessToken, http.StatusOK},
			{"RegularUser", userAuthResponse.AccessToken, http.StatusOK},
			{"NoToken", "", http.StatusUnauthorized},
			{"InvalidToken", "invalid-token", http.StatusUnauthorized},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Create test request
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				c.Request = httptest.NewRequest(http.MethodPost, "/api/cdn/resize/media/"+filename+"?type=image&width=100&height=100", nil)
				c.Params = gin.Params{{Key: "filename", Value: filename}}
				if tc.token != "" {
					c.Request.Header.Add("Authorization", "Bearer "+tc.token)
				}

				// Test the unified media resize endpoint
				mediaHandler.HandleMediaResize(c)

				// Should enforce access control
				require.Equal(t, tc.expectedStatus, w.Result().StatusCode)
			})
		}
	})

	t.Run("MediaListAccessControl", func(t *testing.T) {
		testCases := []struct {
			name           string
			token          string
			expectedStatus int
		}{
			{"AdminUser", adminAuthResponse.AccessToken, http.StatusOK},
			{"RegularUser", userAuthResponse.AccessToken, http.StatusOK},
			{"NoToken", "", http.StatusUnauthorized},
			{"InvalidToken", "invalid-token", http.StatusUnauthorized},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Create test request
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				c.Request = httptest.NewRequest(http.MethodGet, "/api/cdn/media?type=image", nil)
				if tc.token != "" {
					c.Request.Header.Add("Authorization", "Bearer "+tc.token)
				}

				// Test the unified media list endpoint
				mediaHandler.HandleAllMedia(c)

				// Should enforce access control
				require.Equal(t, tc.expectedStatus, w.Result().StatusCode)
			})
		}
	})

	t.Run("CrossUserAccessControl", func(t *testing.T) {
		// Create a second regular user
		secondUser := &models.User{
			Email: "user2@example.com",
			Role:  "user",
		}
		err = secondUser.HashPassword("SecureP@ssw0rd123")
		require.NoError(t, err)
		err = userRepo.CreateUser(secondUser)
		require.NoError(t, err)

		// Get second user token
		secondUserLoginReq := authHandlers.LoginRequest{
			Email:    "user2@example.com",
			Password: "SecureP@ssw0rd123",
		}

		w = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(w)

		jsonData, _ := json.Marshal(secondUserLoginReq)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(jsonData))
		c.Request.Header.Add("Content-Type", "application/json")

		authHandler.Login(c)
		require.Equal(t, http.StatusOK, w.Result().StatusCode)

		var secondUserAuthResponse authHandlers.AuthResponse
		err = json.Unmarshal(w.Body.Bytes(), &secondUserAuthResponse)
		require.NoError(t, err)

		// First user uploads a file
		img, err := testutils.CreateDummyImage(200, 200)
		require.NoError(t, err)

		// Create multipart form with image
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		part, err := writer.CreateFormFile("file", "user1-file.png")
		require.NoError(t, err)

		err = testutils.EncodeImage(part, img)
		require.NoError(t, err)

		err = writer.Close()
		require.NoError(t, err)

		w = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(w)

		c.Request = httptest.NewRequest(http.MethodPost, "/api/cdn/upload/media", body)
		c.Request.Header.Add("Content-Type", writer.FormDataContentType())
		c.Request.Header.Add("Authorization", "Bearer "+userAuthResponse.AccessToken)

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

		// Second user tries to delete first user's file
		w = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(w)

		c.Request = httptest.NewRequest(http.MethodDelete, "/api/cdn/media/"+filename+"?type=image", nil)
		c.Params = gin.Params{{Key: "filename", Value: filename}}
		c.Request.Header.Add("Authorization", "Bearer "+secondUserAuthResponse.AccessToken)

		// Test the unified media delete endpoint
		mediaHandler.HandleMediaDelete(c)

		// Should allow deletion (currently no user-specific access control)
		require.Equal(t, http.StatusOK, w.Result().StatusCode)
	})

	t.Run("RoleBasedAccessControl", func(t *testing.T) {
		// Test if admin-only endpoints are properly protected
		// Note: Currently, there are no admin-only endpoints in the media handlers
		// This test documents the current behavior and highlights a potential security improvement

		// Create test image
		img, err := testutils.CreateDummyImage(200, 200)
		require.NoError(t, err)

		// Create multipart form with image
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		part, err := writer.CreateFormFile("file", "test-rbac.png")
		require.NoError(t, err)

		err = testutils.EncodeImage(part, img)
		require.NoError(t, err)

		err = writer.Close()
		require.NoError(t, err)

		// Both admin and regular users can upload files
		testCases := []struct {
			name           string
			token          string
			expectedStatus int
		}{
			{"AdminUser", adminAuthResponse.AccessToken, http.StatusOK},
			{"RegularUser", userAuthResponse.AccessToken, http.StatusOK},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Create test request
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				c.Request = httptest.NewRequest(http.MethodPost, "/api/cdn/upload/media", body)
				c.Request.Header.Add("Content-Type", writer.FormDataContentType())
				c.Request.Header.Add("Authorization", "Bearer "+tc.token)

				// Test the unified media upload endpoint
				mediaHandler.HandleMediaUpload(c)

				// Should allow both admin and regular users to upload files
				require.Equal(t, tc.expectedStatus, w.Result().StatusCode)
			})
		}
	})

	t.Run("ExpiredTokenAccessControl", func(t *testing.T) {
		// Create test image
		img, err := testutils.CreateDummyImage(200, 200)
		require.NoError(t, err)

		// Create multipart form with image
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		part, err := writer.CreateFormFile("file", "test-expired.png")
		require.NoError(t, err)

		err = testutils.EncodeImage(part, img)
		require.NoError(t, err)

		err = writer.Close()
		require.NoError(t, err)

		// Create test request with expired token
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest(http.MethodPost, "/api/cdn/upload/media", body)
		c.Request.Header.Add("Content-Type", writer.FormDataContentType())
		c.Request.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJlbWFpbCI6InRlc3RAZXhhbXBsZS5jb20iLCJyb2xlIjoidXNlciIsImV4cCI6MTAwMDAwMDAwMH0.invalid")

		// Test the unified media upload endpoint
		mediaHandler.HandleMediaUpload(c)

		// Should reject requests with expired tokens
		require.Equal(t, http.StatusUnauthorized, w.Result().StatusCode)
	})

	t.Run("TamperedTokenAccessControl", func(t *testing.T) {
		// Create test image
		img, err := testutils.CreateDummyImage(200, 200)
		require.NoError(t, err)

		// Create multipart form with image
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		part, err := writer.CreateFormFile("file", "test-tampered.png")
		require.NoError(t, err)

		err = testutils.EncodeImage(part, img)
		require.NoError(t, err)

		err = writer.Close()
		require.NoError(t, err)

		// Create test request with tampered token
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest(http.MethodPost, "/api/cdn/upload/media", body)
		c.Request.Header.Add("Content-Type", writer.FormDataContentType())
		c.Request.Header.Add("Authorization", "Bearer "+userAuthResponse.AccessToken+"tampered")

		// Test the unified media upload endpoint
		mediaHandler.HandleMediaUpload(c)

		// Should reject requests with tampered tokens
		require.Equal(t, http.StatusUnauthorized, w.Result().StatusCode)
	})
}
