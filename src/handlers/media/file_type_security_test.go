package handlers

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/kevinanielsen/go-fast-cdn/src/database"
	authHandlers "github.com/kevinanielsen/go-fast-cdn/src/handlers/auth"
	"github.com/kevinanielsen/go-fast-cdn/src/models"
	"github.com/kevinanielsen/go-fast-cdn/src/testutils"
	"github.com/kevinanielsen/go-fast-cdn/src/util"
	"github.com/stretchr/testify/require"
)

// TestFileTypeSecurity tests various security aspects of file type validation for media operations
func TestFileTypeSecurity(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "file-type-test")
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

	t.Run("ValidImageTypes", func(t *testing.T) {
		testCases := []struct {
			name         string
			imageType    string
			filename     string
			expectedMIME string
		}{
			{"JPEGImage", "jpeg", "test.jpg", "image/jpeg"},
			{"PNGImage", "png", "test.png", "image/png"},
			{"GIFImage", "gif", "test.gif", "image/gif"},
			{"WebPImage", "webp", "test.webp", "image/webp"},
			{"BMPImage", "bmp", "test.bmp", "image/bmp"},
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

				// Should accept valid image types
				require.Equal(t, http.StatusOK, w.Result().StatusCode)
			})
		}
	})

	t.Run("ValidDocumentTypes", func(t *testing.T) {
		testCases := []struct {
			name         string
			documentType string
			filename     string
			content      string
			expectedMIME string
		}{
			{"PDFDocument", "pdf", "test.pdf", "%PDF-1.4\n1 0 obj\n<<\n/Type /Catalog\n/Pages 2 0 R\n>>\nendobj\n", "application/pdf"},
			{"TextDocument", "txt", "test.txt", "This is a text document.", "text/plain; charset=utf-8"},
			{"JSONDocument", "json", "test.json", `{"key": "value"}`, "application/json"},
			{"XMLDocument", "xml", "test.xml", "<?xml version=\"1.0\" encoding=\"UTF-8\"?><root></root>", "text/xml; charset=utf-8"},
			{"CSVDocument", "csv", "test.csv", "name,age\nJohn,30\nJane,25\n", "text/csv"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Create multipart form with document
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

				// Should accept valid document types
				require.Equal(t, http.StatusOK, w.Result().StatusCode)
			})
		}
	})

	t.Run("InvalidFileTypes", func(t *testing.T) {
		testCases := []struct {
			name          string
			filename      string
			content       string
			expectedError string
		}{
			{"ExecutableFile", "test.exe", "MZ\x90\x00\x03\x00\x00\x00", "Invalid file type"},
			{"ScriptFile", "test.js", "console.log('Hello, world!');", "Invalid file type"},
			{"PHPFile", "test.php", "<?php echo 'Hello, world!'; ?>", "Invalid file type"},
			{"ShellScript", "test.sh", "#!/bin/bash\necho 'Hello, world!'", "Invalid file type"},
			{"PythonScript", "test.py", "print('Hello, world!')", "Invalid file type"},
			{"RubyScript", "test.rb", "puts 'Hello, world!'", "Invalid file type"},
			{"PerlScript", "test.pl", "print 'Hello, world!';", "Invalid file type"},
			{"JavaClass", "test.class", "\xCA\xFE\xBA\xBE", "Invalid file type"},
			{"DLLFile", "test.dll", "MZ\x90\x00\x03\x00\x00\x00", "Invalid file type"},
			{"ShortcutFile", "test.lnk", "\x4C\x00\x00\x00\x01\x14\x02\x00", "Invalid file type"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Create multipart form with invalid file
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

				// Should reject invalid file types
				require.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
				responseBody := w.Body.String()
				require.Contains(t, responseBody, tc.expectedError)
			})
		}
	})

	t.Run("MIMETypeSpoofing", func(t *testing.T) {
		testCases := []struct {
			name          string
			filename      string
			content       string
			contentType   string
			expectedError string
		}{
			{"EXEWithImageExtension", "test.jpg", "MZ\x90\x00\x03\x00\x00\x00", "image/jpeg", "Invalid file type"},
			{"ScriptWithDocumentExtension", "test.pdf", "<?php echo 'Hello, world!'; ?>", "application/pdf", "Invalid file type"},
			{"ShellScriptWithTextExtension", "test.txt", "#!/bin/bash\necho 'Hello, world!'", "text/plain", "Invalid file type"},
			{"JavaClassWithPNGExtension", "test.png", "\xCA\xFE\xBA\xBE", "image/png", "Invalid file type"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Create multipart form with spoofed file
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

				// Should reject files with spoofed MIME types
				require.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
				responseBody := w.Body.String()
				require.Contains(t, responseBody, tc.expectedError)
			})
		}
	})

	t.Run("FileHeaderManipulation", func(t *testing.T) {
		testCases := []struct {
			name          string
			filename      string
			content       string
			expectedError string
		}{
			{"EXEWithPNGHeader", "test.png", "\x89PNG\r\n\x1A\nMZ\x90\x00\x03\x00\x00\x00", "Invalid file type"},
			{"ScriptWithPDFHeader", "test.pdf", "%PDF-1.4\n<?php echo 'Hello, world!'; ?>", "Invalid file type"},
			{"ShellScriptWithJPEGHeader", "test.jpg", "\xFF\xD8\xFF\xE0#!/bin/bash\necho 'Hello, world!'", "Invalid file type"},
			{"JavaClassWithGIFHeader", "test.gif", "GIF87a\xCA\xFE\xBA\xBE", "Invalid file type"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Create multipart form with manipulated file header
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

				// Should reject files with manipulated headers
				require.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
				responseBody := w.Body.String()
				require.Contains(t, responseBody, tc.expectedError)
			})
		}
	})

	t.Run("PolyglotFiles", func(t *testing.T) {
		testCases := []struct {
			name          string
			filename      string
			content       string
			expectedError string
		}{
			{"GIFARFile", "test.gif", "GIF87a\xFF\xD8\xFF\xE0", "Invalid file type"},                // GIF + JPEG
			{"PDFEXEFile", "test.pdf", "%PDF-1.4\nMZ\x90\x00\x03\x00\x00\x00", "Invalid file type"}, // PDF + EXE
			{"ZIPGIFFile", "test.gif", "GIF87aPK\x03\x04", "Invalid file type"},                     // GIF + ZIP
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Create multipart form with polyglot file
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

				// Should reject polyglot files
				require.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
				responseBody := w.Body.String()
				require.Contains(t, responseBody, tc.expectedError)
			})
		}
	})

	t.Run("EmptyFiles", func(t *testing.T) {
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
		responseBody := w.Body.String()
		require.Contains(t, responseBody, "Invalid file type")
	})

	t.Run("FilesWithNullBytes", func(t *testing.T) {
		// Create multipart form with file containing null bytes
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		part, err := writer.CreateFormFile("file", "null.txt")
		require.NoError(t, err)

		// Write content with null bytes
		_, err = part.Write([]byte("This\x00is\x00a\x00test\x00"))
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

		// Should handle files with null bytes gracefully
		require.Equal(t, http.StatusOK, w.Result().StatusCode)
	})

	t.Run("FilesWithControlCharacters", func(t *testing.T) {
		// Create multipart form with file containing control characters
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		part, err := writer.CreateFormFile("file", "control.txt")
		require.NoError(t, err)

		// Write content with control characters
		_, err = part.Write([]byte("This\x01is\x02a\x03test\x04"))
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

		// Should handle files with control characters gracefully
		require.Equal(t, http.StatusOK, w.Result().StatusCode)
	})

	t.Run("FilesWithUnicode", func(t *testing.T) {
		// Create multipart form with file containing unicode characters
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		part, err := writer.CreateFormFile("file", "unicode.txt")
		require.NoError(t, err)

		// Write content with unicode characters
		_, err = part.Write([]byte("这是一个测试文件"))
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

		// Should handle files with unicode characters gracefully
		require.Equal(t, http.StatusOK, w.Result().StatusCode)
	})
}

// TestFileExtensionSecurity tests security aspects related to file extensions
func TestFileExtensionSecurity(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "file-extension-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	util.ExPath = tempDir

	// Set environment variables for testing
	os.Setenv("JWT_SECRET", "test-super-secret-jwt-key-for-testing-only")
	defer os.Unsetenv("JWT_SECRET")

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

	t.Run("DoubleExtensions", func(t *testing.T) {
		testCases := []struct {
			name          string
			filename      string
			content       string
			expectedError string
		}{
			{"EXEWithJPGExtension", "test.jpg.exe", "MZ\x90\x00\x03\x00\x00\x00", "Invalid filename"},
			{"PHPWithTXTExtension", "test.txt.php", "<?php echo 'Hello, world!'; ?>", "Invalid filename"},
			{"ShellWithPDFExtension", "test.pdf.sh", "#!/bin/bash\necho 'Hello, world!'", "Invalid filename"},
			{"JSWithPNGExtension", "test.png.js", "console.log('Hello, world!');", "Invalid filename"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Create multipart form with file having double extension
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

				// Should reject files with double extensions
				require.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
				responseBody := w.Body.String()
				require.Contains(t, responseBody, tc.expectedError)
			})
		}
	})

	t.Run("HiddenFiles", func(t *testing.T) {
		testCases := []struct {
			name          string
			filename      string
			content       string
			shouldSucceed bool
		}{
			{"HiddenTextFile", ".hidden.txt", "This is a hidden text file.", true},
			{"HiddenImageFile", ".hidden.png", "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNkYPhfDwAChwGA60e6kgAAAABJRU5ErkJggg==", true},
			{"HiddenExecutableFile", ".hidden.exe", "MZ\x90\x00\x03\x00\x00\x00", false},
			{"HiddenScriptFile", ".hidden.php", "<?php echo 'Hello, world!'; ?>", false},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Create multipart form with hidden file
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

				if tc.shouldSucceed {
					require.Equal(t, http.StatusOK, w.Result().StatusCode)
				} else {
					require.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
					responseBody := w.Body.String()
					require.Contains(t, responseBody, "Invalid file type")
				}
			})
		}
	})

	t.Run("CaseSensitiveExtensions", func(t *testing.T) {
		testCases := []struct {
			name          string
			filename      string
			content       string
			shouldSucceed bool
		}{
			{"UppercaseJPG", "test.JPG", "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNkYPhfDwAChwGA60e6kgAAAABJRU5ErkJggg==", true},
			{"UppercasePNG", "test.PNG", "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNkYPhfDwAChwGA60e6kgAAAABJRU5ErkJggg==", true},
			{"UppercaseEXE", "test.EXE", "MZ\x90\x00\x03\x00\x00\x00", false},
			{"UppercasePHP", "test.PHP", "<?php echo 'Hello, world!'; ?>", false},
			{"MixedCaseJpg", "test.JpG", "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNkYPhfDwAChwGA60e6kgAAAABJRU5ErkJggg==", true},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Create multipart form with file having case-sensitive extension
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

				if tc.shouldSucceed {
					require.Equal(t, http.StatusOK, w.Result().StatusCode)
				} else {
					require.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
					responseBody := w.Body.String()
					require.Contains(t, responseBody, "Invalid file type")
				}
			})
		}
	})

	t.Run("ExtensionlessFiles", func(t *testing.T) {
		// Create multipart form with extensionless file
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		part, err := writer.CreateFormFile("file", "extensionless")
		require.NoError(t, err)

		_, err = part.Write([]byte("This file has no extension."))
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

		// Should handle extensionless files gracefully
		require.Equal(t, http.StatusOK, w.Result().StatusCode)
	})

	t.Run("FilesWithTrailingSpaces", func(t *testing.T) {
		testCases := []struct {
			name          string
			filename      string
			content       string
			shouldSucceed bool
		}{
			{"JPGWithTrailingSpace", "test.jpg ", "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNkYPhfDwAChwGA60e6kgAAAABJRU5ErkJggg==", true},
			{"PNGWithTrailingSpace", "test.png ", "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNkYPhfDwAChwGA60e6kgAAAABJRU5ErkJggg==", true},
			{"EXEWithTrailingSpace", "test.exe ", "MZ\x90\x00\x03\x00\x00\x00", false},
			{"PHPWithTrailingSpace", "test.php ", "<?php echo 'Hello, world!'; ?>", false},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Create multipart form with file having trailing space in filename
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

				if tc.shouldSucceed {
					require.Equal(t, http.StatusOK, w.Result().StatusCode)
				} else {
					require.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
					responseBody := w.Body.String()
					require.Contains(t, responseBody, "Invalid file type")
				}
			})
		}
	})

	t.Run("FilesWithLeadingSpaces", func(t *testing.T) {
		testCases := []struct {
			name          string
			filename      string
			content       string
			shouldSucceed bool
		}{
			{"JPGWithLeadingSpace", " test.jpg", "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNkYPhfDwAChwGA60e6kgAAAABJRU5ErkJggg==", true},
			{"PNGWithLeadingSpace", " test.png", "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNkYPhfDwAChwGA60e6kgAAAABJRU5ErkJggg==", true},
			{"EXEWithLeadingSpace", " test.exe", "MZ\x90\x00\x03\x00\x00\x00", false},
			{"PHPWithLeadingSpace", " test.php", "<?php echo 'Hello, world!'; ?>", false},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Create multipart form with file having leading space in filename
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

				if tc.shouldSucceed {
					require.Equal(t, http.StatusOK, w.Result().StatusCode)
				} else {
					require.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
					responseBody := w.Body.String()
					require.Contains(t, responseBody, "Invalid file type")
				}
			})
		}
	})
}
