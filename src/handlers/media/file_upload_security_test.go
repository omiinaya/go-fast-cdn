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

// TestFileUploadSecurity tests various security aspects of file upload operations
func TestFileUploadSecurity(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "file-upload-test")
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

	t.Run("FileSizeLimits", func(t *testing.T) {
		testCases := []struct {
			name           string
			fileSize       int
			expectedStatus int
		}{
			{"SmallFile", 1024, http.StatusOK},
			{"MediumFile", 1024 * 1024, http.StatusOK},          // 1MB
			{"LargeFile", 10 * 1024 * 1024, http.StatusOK},      // 10MB
			{"VeryLargeFile", 100 * 1024 * 1024, http.StatusOK}, // 100MB
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Create a file with the specified size
				content := make([]byte, tc.fileSize)
				for i := range content {
					content[i] = byte(i % 256)
				}

				// Create multipart form with file
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)

				part, err := writer.CreateFormFile("file", "test.txt")
				require.NoError(t, err)

				_, err = part.Write(content)
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

				// Should handle files of various sizes
				require.Equal(t, tc.expectedStatus, w.Result().StatusCode)
			})
		}
	})

	t.Run("MaliciousFileContent", func(t *testing.T) {
		testCases := []struct {
			name           string
			filename       string
			content        string
			expectedStatus int
		}{
			{"FileWithHTMLScript", "test.html", "<script>alert('XSS')</script>", http.StatusOK},
			{"FileWithJavaScript", "test.js", "alert('XSS')", http.StatusOK},
			{"FileWithPHPCode", "test.php", "<?php echo 'Hello, world!'; ?>", http.StatusOK},
			{"FileWithShellScript", "test.sh", "#!/bin/bash\necho 'Hello, world!'", http.StatusOK},
			{"FileWithPythonScript", "test.py", "print('Hello, world!')", http.StatusOK},
			{"FileWithRubyScript", "test.rb", "puts 'Hello, world!'", http.StatusOK},
			{"FileWithPerlScript", "test.pl", "print 'Hello, world!';", http.StatusOK},
			{"FileWithJavaClass", "test.class", "\xCA\xFE\xBA\xBE", http.StatusOK},
			{"FileWithDLL", "test.dll", "MZ\x90\x00\x03\x00\x00\x00", http.StatusOK},
			{"FileWithEXE", "test.exe", "MZ\x90\x00\x03\x00\x00\x00", http.StatusOK},
			{"FileWithShortcut", "test.lnk", "\x4C\x00\x00\x00\x01\x14\x02\x00", http.StatusOK},
			{"FileWithBatchScript", "test.bat", "@echo off\necho Hello, world!", http.StatusOK},
			{"FileWithPowerShellScript", "test.ps1", "Write-Host 'Hello, world!'", http.StatusOK},
			{"FileWithVBScript", "test.vbs", "MsgBox \"Hello, world!\"", http.StatusOK},
			{"FileWithHTA", "test.hta", "<html><head><hta:application><script>alert('XSS')</script></head></html>", http.StatusOK},
			{"FileWithSCRF", "test.scrf", "[InternetShortcut]\nURL=javascript:alert('XSS')", http.StatusOK},
			{"FileWithURL", "test.url", "[InternetShortcut]\nURL=http://example.com", http.StatusOK},
			{"FileWithWSF", "test.wsf", "<package><job><script>alert('XSS')</script></job></package>", http.StatusOK},
			{"FileWithPif", "test.pif", "MZ\x90\x00\x03\x00\x00\x00", http.StatusOK},
			{"FileWithCom", "test.com", "MZ\x90\x00\x03\x00\x00\x00", http.StatusOK},
			{"FileWithCpl", "test.cpl", "MZ\x90\x00\x03\x00\x00\x00", http.StatusOK},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Create multipart form with malicious file
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)

				part, err := writer.CreateFormFile("file", tc.filename)
				require.NoError(t, err)

				_, err = part.Write([]byte(tc.content))
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

				// Should handle malicious files based on their content type
				require.Equal(t, tc.expectedStatus, w.Result().StatusCode)
			})
		}
	})

	t.Run("FileUploadWithoutAuthentication", func(t *testing.T) {
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

		// Create test request without authentication
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest(http.MethodPost, "/api/cdn/upload/media", body)
		c.Request.Header.Add("Content-Type", writer.FormDataContentType())

		// Test the unified media upload endpoint
		mediaHandler.HandleMediaUpload(c)

		// Should reject requests without authentication
		require.Equal(t, http.StatusUnauthorized, w.Result().StatusCode)
	})

	t.Run("FileUploadWithInvalidToken", func(t *testing.T) {
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

		// Create test request with invalid token
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest(http.MethodPost, "/api/cdn/upload/media", body)
		c.Request.Header.Add("Content-Type", writer.FormDataContentType())
		c.Request.Header.Add("Authorization", "Bearer invalid-token")

		// Test the unified media upload endpoint
		mediaHandler.HandleMediaUpload(c)

		// Should reject requests with invalid tokens
		require.Equal(t, http.StatusUnauthorized, w.Result().StatusCode)
	})

	t.Run("FileUploadWithExpiredToken", func(t *testing.T) {
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

	t.Run("FileUploadWithTamperedToken", func(t *testing.T) {
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

		// Create test request with tampered token
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest(http.MethodPost, "/api/cdn/upload/media", body)
		c.Request.Header.Add("Content-Type", writer.FormDataContentType())
		c.Request.Header.Add("Authorization", "Bearer "+authResponse.AccessToken+"tampered")

		// Test the unified media upload endpoint
		mediaHandler.HandleMediaUpload(c)

		// Should reject requests with tampered tokens
		require.Equal(t, http.StatusUnauthorized, w.Result().StatusCode)
	})

	t.Run("FileUploadWithMissingFile", func(t *testing.T) {
		// Create multipart form without file
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		err := writer.Close()
		require.NoError(t, err)

		// Create test request
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest(http.MethodPost, "/api/cdn/upload/media", body)
		c.Request.Header.Add("Content-Type", writer.FormDataContentType())
		c.Request.Header.Add("Authorization", "Bearer "+authResponse.AccessToken)

		// Test the unified media upload endpoint
		mediaHandler.HandleMediaUpload(c)

		// Should reject requests without files
		require.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	})

	t.Run("FileUploadWithEmptyFile", func(t *testing.T) {
		// Create multipart form with empty file
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		part, err := writer.CreateFormFile("file", "empty.txt")
		require.NoError(t, err)

		// Write empty content
		_, err = part.Write([]byte(""))
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

		// Should reject empty files
		require.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	})

	t.Run("FileUploadWithMultipleFiles", func(t *testing.T) {
		// Create test images
		img1, err := testutils.CreateDummyImage(200, 200)
		require.NoError(t, err)

		img2, err := testutils.CreateDummyImage(300, 300)
		require.NoError(t, err)

		// Create multipart form with multiple files
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		part1, err := writer.CreateFormFile("file", "test1.png")
		require.NoError(t, err)

		err = testutils.EncodeImage(part1, img1)
		require.NoError(t, err)

		part2, err := writer.CreateFormFile("file", "test2.png")
		require.NoError(t, err)

		err = testutils.EncodeImage(part2, img2)
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

		// Should handle multiple files (depending on implementation)
		// Currently, the API only processes the first file
		require.Equal(t, http.StatusOK, w.Result().StatusCode)
	})

	t.Run("FileUploadWithInvalidContentType", func(t *testing.T) {
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

		// Create test request with invalid content type
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest(http.MethodPost, "/api/cdn/upload/media", body)
		c.Request.Header.Add("Content-Type", "invalid-content-type")
		c.Request.Header.Add("Authorization", "Bearer "+authResponse.AccessToken)

		// Test the unified media upload endpoint
		mediaHandler.HandleMediaUpload(c)

		// Should handle invalid content types gracefully
		require.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	})

	t.Run("FileUploadWithMissingContentType", func(t *testing.T) {
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

		// Create test request without content type
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest(http.MethodPost, "/api/cdn/upload/media", body)
		c.Request.Header.Add("Authorization", "Bearer "+authResponse.AccessToken)

		// Test the unified media upload endpoint
		mediaHandler.HandleMediaUpload(c)

		// Should handle missing content types gracefully
		require.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	})

	t.Run("FileUploadWithMaliciousBoundary", func(t *testing.T) {
		// Create test image
		img, err := testutils.CreateDummyImage(200, 200)
		require.NoError(t, err)

		// Create multipart form with malicious boundary
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		// Manually create multipart form with malicious boundary
		boundary := "----WebKitFormBoundary" + strings.Repeat("A", 1000)
		writer.SetBoundary(boundary)

		part, err := writer.CreateFormFile("file", "test.png")
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

		// Should handle malicious boundaries gracefully
		require.Equal(t, http.StatusOK, w.Result().StatusCode)
	})

	t.Run("FileUploadWithMaliciousFormData", func(t *testing.T) {
		// Create test image
		img, err := testutils.CreateDummyImage(200, 200)
		require.NoError(t, err)

		// Create multipart form with malicious form data
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		// Add malicious form data
		err = writer.WriteField("filename", "<script>alert('XSS')</script>")
		require.NoError(t, err)

		part, err := writer.CreateFormFile("file", "test.png")
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

		// Should handle malicious form data gracefully
		require.Equal(t, http.StatusOK, w.Result().StatusCode)
	})

	t.Run("FileUploadWithMaliciousHeader", func(t *testing.T) {
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

		// Create test request with malicious headers
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest(http.MethodPost, "/api/cdn/upload/media", body)
		c.Request.Header.Add("Content-Type", writer.FormDataContentType())
		c.Request.Header.Add("Authorization", "Bearer "+authResponse.AccessToken)
		c.Request.Header.Add("X-Forwarded-For", "<script>alert('XSS')</script>")
		c.Request.Header.Add("User-Agent", "<script>alert('XSS')</script>")
		c.Request.Header.Add("Referer", "javascript:alert('XSS')")

		// Test the unified media upload endpoint
		mediaHandler.HandleMediaUpload(c)

		// Should handle malicious headers gracefully
		require.Equal(t, http.StatusOK, w.Result().StatusCode)
	})

	t.Run("FileUploadWithMaliciousFilename", func(t *testing.T) {
		testCases := []struct {
			name           string
			filename       string
			expectedStatus int
		}{
			{"FilenameWithNullByte", "test\x00.png", http.StatusBadRequest},
			{"FilenameWithControlChar", "test\x01.png", http.StatusBadRequest},
			{"FilenameWithBackspace", "test\x08.png", http.StatusBadRequest},
			{"FilenameWithVerticalTab", "test\x0b.png", http.StatusBadRequest},
			{"FilenameWithFormFeed", "test\x0c.png", http.StatusBadRequest},
			{"FilenameWithCarriageReturn", "test\x0d.png", http.StatusBadRequest},
			{"FilenameWithDelete", "test\x7f.png", http.StatusBadRequest},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Create test image
				img, err := testutils.CreateDummyImage(200, 200)
				require.NoError(t, err)

				// Create multipart form with image
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)

				part, err := writer.CreateFormFile("file", tc.filename)
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

				// Should reject filenames with control characters
				require.Equal(t, tc.expectedStatus, w.Result().StatusCode)
			})
		}
	})

	t.Run("FileUploadWithDuplicateFiles", func(t *testing.T) {
		// Create test image
		img, err := testutils.CreateDummyImage(200, 200)
		require.NoError(t, err)

		// Upload the same file twice
		for i := 0; i < 2; i++ {
			// Create multipart form with image
			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)

			part, err := writer.CreateFormFile("file", "test.png")
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

			if i == 0 {
				// First upload should succeed
				require.Equal(t, http.StatusOK, w.Result().StatusCode)
			} else {
				// Second upload with the same content should fail
				require.Equal(t, http.StatusConflict, w.Result().StatusCode)
			}
		}
	})

	t.Run("FileUploadWithSpecialCharactersInFilename", func(t *testing.T) {
		testCases := []struct {
			name           string
			filename       string
			expectedStatus int
		}{
			{"FilenameWithSpaces", "test image.png", http.StatusOK},
			{"FilenameWithUnderscore", "test_image.png", http.StatusOK},
			{"FilenameWithDash", "test-image.png", http.StatusOK},
			{"FilenameWithDots", "test.image.png", http.StatusOK},
			{"FilenameWithUnicode", "测试图像.png", http.StatusOK},
			{"FilenameWithEmoji", "test🖼️.png", http.StatusOK},
			{"FilenameWithAccents", "testéàüö.png", http.StatusOK},
			{"FilenameWithCyrillic", "тест.png", http.StatusOK},
			{"FilenameWithChinese", "测试.png", http.StatusOK},
			{"FilenameWithJapanese", "テスト.png", http.StatusOK},
			{"FilenameWithKorean", "테스트.png", http.StatusOK},
			{"FilenameWithArabic", "اختبار.png", http.StatusOK},
			{"FilenameWithHebrew", "בדיקה.png", http.StatusOK},
			{"FilenameWithThai", "ทดสอบ.png", http.StatusOK},
			{"FilenameWithHindi", "परीक्षण.png", http.StatusOK},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Create test image
				img, err := testutils.CreateDummyImage(200, 200)
				require.NoError(t, err)

				// Create multipart form with image
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)

				part, err := writer.CreateFormFile("file", tc.filename)
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

				// Should handle special characters in filenames
				require.Equal(t, tc.expectedStatus, w.Result().StatusCode)
			})
		}
	})
}
