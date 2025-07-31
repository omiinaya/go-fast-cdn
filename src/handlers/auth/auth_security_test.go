package auth

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kevinanielsen/go-fast-cdn/src/database"
	"github.com/kevinanielsen/go-fast-cdn/src/middleware"
	"github.com/kevinanielsen/go-fast-cdn/src/models"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

// TestAuthenticationSecurity tests various security aspects of the authentication system
func TestAuthenticationSecurity(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "auth-security-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Set environment variables for testing
	os.Setenv("JWT_SECRET", "test-super-secret-jwt-key-for-testing-only")
	os.Setenv("JWT_EXPIRES_IN", "900") // 15 minutes
	defer func() {
		os.Unsetenv("JWT_SECRET")
		os.Unsetenv("JWT_EXPIRES_IN")
	}()

	// Initialize database
	database.ConnectToDB()
	defer func() {
		dbPath := filepath.Join(tempDir, database.DbFolder, database.DbName)
		if err := os.Remove(dbPath); err != nil && !os.IsNotExist(err) {
			t.Error(err)
		}
	}()

	// Create auth handler
	userRepo := database.NewUserRepo(database.DB)
	authHandler := NewAuthHandler(userRepo)
	authMiddleware := middleware.NewAuthMiddleware()

	t.Run("PasswordSecurity", func(t *testing.T) {
		t.Run("WeakPasswordRejected", func(t *testing.T) {
			// Test with weak password
			weakPasswordReq := RegisterRequest{
				Email:    "test@example.com",
				Password: "123", // Too weak
			}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			jsonData, _ := json.Marshal(weakPasswordReq)
			c.Request = httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBuffer(jsonData))
			c.Request.Header.Add("Content-Type", "application/json")

			authHandler.Register(c)

			// Should reject weak password
			require.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
			responseBody := w.Body.String()
			require.Contains(t, responseBody, "Validation failed")
		})

		t.Run("StrongPasswordAccepted", func(t *testing.T) {
			// Test with strong password
			strongPasswordReq := RegisterRequest{
				Email:    "strong@example.com",
				Password: "Str0ngP@ssw0rd!",
			}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			jsonData, _ := json.Marshal(strongPasswordReq)
			c.Request = httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBuffer(jsonData))
			c.Request.Header.Add("Content-Type", "application/json")

			authHandler.Register(c)

			// Should accept strong password
			require.Equal(t, http.StatusCreated, w.Result().StatusCode)
		})
	})

	t.Run("TokenSecurity", func(t *testing.T) {
		// First, register and login a user
		registerReq := RegisterRequest{
			Email:    "token@example.com",
			Password: "SecureP@ssw0rd123",
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		jsonData, _ := json.Marshal(registerReq)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBuffer(jsonData))
		c.Request.Header.Add("Content-Type", "application/json")

		authHandler.Register(c)
		require.Equal(t, http.StatusCreated, w.Result().StatusCode)

		// Login to get tokens
		loginReq := LoginRequest{
			Email:    "token@example.com",
			Password: "SecureP@ssw0rd123",
		}

		w = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(w)

		jsonData, _ = json.Marshal(loginReq)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(jsonData))
		c.Request.Header.Add("Content-Type", "application/json")

		authHandler.Login(c)
		require.Equal(t, http.StatusOK, w.Result().StatusCode)

		var authResponse AuthResponse
		err := json.Unmarshal(w.Body.Bytes(), &authResponse)
		require.NoError(t, err)

		// Test token security
		t.Run("InvalidTokenRejected", func(t *testing.T) {
			w = httptest.NewRecorder()
			c, _ = gin.CreateTestContext(w)

			c.Request = httptest.NewRequest(http.MethodGet, "/api/auth/profile", nil)
			c.Request.Header.Add("Authorization", "Bearer invalid-token-12345")

			// Apply auth middleware
			handler := authMiddleware.RequireAuth()
			handler(c)

			// Should reject invalid token
			require.Equal(t, http.StatusUnauthorized, w.Result().StatusCode)
			responseBody := w.Body.String()
			require.Contains(t, responseBody, "Invalid token")
		})

		t.Run("MissingAuthorizationHeader", func(t *testing.T) {
			w = httptest.NewRecorder()
			c, _ = gin.CreateTestContext(w)

			c.Request = httptest.NewRequest(http.MethodGet, "/api/auth/profile", nil)
			// No Authorization header

			// Apply auth middleware
			handler := authMiddleware.RequireAuth()
			handler(c)

			// Should reject missing authorization header
			require.Equal(t, http.StatusUnauthorized, w.Result().StatusCode)
			responseBody := w.Body.String()
			require.Contains(t, responseBody, "Authorization header required")
		})

		t.Run("MalformedAuthorizationHeader", func(t *testing.T) {
			w = httptest.NewRecorder()
			c, _ = gin.CreateTestContext(w)

			c.Request = httptest.NewRequest(http.MethodGet, "/api/auth/profile", nil)
			c.Request.Header.Add("Authorization", "InvalidFormat "+authResponse.AccessToken)

			// Apply auth middleware
			handler := authMiddleware.RequireAuth()
			handler(c)

			// Should reject malformed authorization header
			require.Equal(t, http.StatusUnauthorized, w.Result().StatusCode)
			responseBody := w.Body.String()
			require.Contains(t, responseBody, "Invalid authorization header format")
		})

		t.Run("ValidTokenAccepted", func(t *testing.T) {
			w = httptest.NewRecorder()
			c, _ = gin.CreateTestContext(w)

			c.Request = httptest.NewRequest(http.MethodGet, "/api/auth/profile", nil)
			c.Request.Header.Add("Authorization", "Bearer "+authResponse.AccessToken)

			// Apply auth middleware
			handler := authMiddleware.RequireAuth()
			handler(c)

			// Should accept valid token
			require.Equal(t, http.StatusOK, w.Result().StatusCode)
		})
	})

	t.Run("BruteForceProtection", func(t *testing.T) {
		// Create a user for testing
		registerReq := RegisterRequest{
			Email:    "bruteforce@example.com",
			Password: "SecureP@ssw0rd123",
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		jsonData, _ := json.Marshal(registerReq)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBuffer(jsonData))
		c.Request.Header.Add("Content-Type", "application/json")

		authHandler.Register(c)
		require.Equal(t, http.StatusCreated, w.Result().StatusCode)

		// Attempt multiple failed logins
		loginReq := LoginRequest{
			Email:    "bruteforce@example.com",
			Password: "WrongPassword",
		}

		for i := 0; i < 10; i++ {
			w = httptest.NewRecorder()
			c, _ = gin.CreateTestContext(w)

			jsonData, _ := json.Marshal(loginReq)
			c.Request = httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(jsonData))
			c.Request.Header.Add("Content-Type", "application/json")

			authHandler.Login(c)

			// All attempts should fail
			require.Equal(t, http.StatusUnauthorized, w.Result().StatusCode)
		}

		// Verify that the account is not locked out (current implementation doesn't have lockout)
		// This test documents the current behavior and highlights a potential security improvement
	})

	t.Run("SessionManagement", func(t *testing.T) {
		// Register and login a user
		registerReq := RegisterRequest{
			Email:    "session@example.com",
			Password: "SecureP@ssw0rd123",
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		jsonData, _ := json.Marshal(registerReq)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBuffer(jsonData))
		c.Request.Header.Add("Content-Type", "application/json")

		authHandler.Register(c)
		require.Equal(t, http.StatusCreated, w.Result().StatusCode)

		// Login to get tokens
		loginReq := LoginRequest{
			Email:    "session@example.com",
			Password: "SecureP@ssw0rd123",
		}

		w = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(w)

		jsonData, _ = json.Marshal(loginReq)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(jsonData))
		c.Request.Header.Add("Content-Type", "application/json")

		authHandler.Login(c)
		require.Equal(t, http.StatusOK, w.Result().StatusCode)

		var authResponse AuthResponse
		err := json.Unmarshal(w.Body.Bytes(), &authResponse)
		require.NoError(t, err)

		// Test logout functionality
		t.Run("LogoutInvalidatesToken", func(t *testing.T) {
			// Logout the user
			logoutReq := RefreshRequest{
				RefreshToken: authResponse.RefreshToken,
			}

			w = httptest.NewRecorder()
			c, _ = gin.CreateTestContext(w)

			jsonData, _ := json.Marshal(logoutReq)
			c.Request = httptest.NewRequest(http.MethodPost, "/api/auth/logout", bytes.NewBuffer(jsonData))
			c.Request.Header.Add("Content-Type", "application/json")

			authHandler.Logout(c)
			require.Equal(t, http.StatusOK, w.Result().StatusCode)

			// Try to access protected endpoint with the same token
			w = httptest.NewRecorder()
			c, _ = gin.CreateTestContext(w)

			c.Request = httptest.NewRequest(http.MethodGet, "/api/auth/profile", nil)
			c.Request.Header.Add("Authorization", "Bearer "+authResponse.AccessToken)

			// Apply auth middleware
			handler := authMiddleware.RequireAuth()
			handler(c)

			// Token should still be valid as JWT doesn't support server-side invalidation
			// This test documents the current behavior and highlights a potential security improvement
			require.Equal(t, http.StatusOK, w.Result().StatusCode)
		})

		t.Run("RefreshTokenRotation", func(t *testing.T) {
			// Create a new user for this test
			registerReq := RegisterRequest{
				Email:    "refresh@example.com",
				Password: "SecureP@ssw0rd123",
			}

			w = httptest.NewRecorder()
			c, _ = gin.CreateTestContext(w)

			jsonData, _ := json.Marshal(registerReq)
			c.Request = httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBuffer(jsonData))
			c.Request.Header.Add("Content-Type", "application/json")

			authHandler.Register(c)
			require.Equal(t, http.StatusCreated, w.Result().StatusCode)

			// Login to get tokens
			loginReq := LoginRequest{
				Email:    "refresh@example.com",
				Password: "SecureP@ssw0rd123",
			}

			w = httptest.NewRecorder()
			c, _ = gin.CreateTestContext(w)

			jsonData, _ = json.Marshal(loginReq)
			c.Request = httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(jsonData))
			c.Request.Header.Add("Content-Type", "application/json")

			authHandler.Login(c)
			require.Equal(t, http.StatusOK, w.Result().StatusCode)

			var firstAuthResponse AuthResponse
			err := json.Unmarshal(w.Body.Bytes(), &firstAuthResponse)
			require.NoError(t, err)

			// Store the original refresh token
			originalRefreshToken := firstAuthResponse.RefreshToken

			// Refresh the tokens
			refreshReq := RefreshRequest{
				RefreshToken: originalRefreshToken,
			}

			w = httptest.NewRecorder()
			c, _ = gin.CreateTestContext(w)

			jsonData, _ = json.Marshal(refreshReq)
			c.Request = httptest.NewRequest(http.MethodPost, "/api/auth/refresh", bytes.NewBuffer(jsonData))
			c.Request.Header.Add("Content-Type", "application/json")

			authHandler.RefreshToken(c)
			require.Equal(t, http.StatusOK, w.Result().StatusCode)

			var secondAuthResponse AuthResponse
			err = json.Unmarshal(w.Body.Bytes(), &secondAuthResponse)
			require.NoError(t, err)

			// Verify that a new refresh token was issued (token rotation)
			require.NotEqual(t, originalRefreshToken, secondAuthResponse.RefreshToken)

			// Try to use the old refresh token
			w = httptest.NewRecorder()
			c, _ = gin.CreateTestContext(w)

			jsonData, _ = json.Marshal(refreshReq)
			c.Request = httptest.NewRequest(http.MethodPost, "/api/auth/refresh", bytes.NewBuffer(jsonData))
			c.Request.Header.Add("Content-Type", "application/json")

			authHandler.RefreshToken(c)

			// Should reject the old refresh token
			require.Equal(t, http.StatusUnauthorized, w.Result().StatusCode)
			responseBody := w.Body.String()
			require.Contains(t, responseBody, "Invalid refresh token")
		})
	})

	t.Run("TwoFactorAuthenticationSecurity", func(t *testing.T) {
		// Register and login a user
		registerReq := RegisterRequest{
			Email:    "2fa@example.com",
			Password: "SecureP@ssw0rd123",
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		jsonData, _ := json.Marshal(registerReq)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBuffer(jsonData))
		c.Request.Header.Add("Content-Type", "application/json")

		authHandler.Register(c)
		require.Equal(t, http.StatusCreated, w.Result().StatusCode)

		// Login to get tokens
		loginReq := LoginRequest{
			Email:    "2fa@example.com",
			Password: "SecureP@ssw0rd123",
		}

		w = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(w)

		jsonData, _ = json.Marshal(loginReq)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(jsonData))
		c.Request.Header.Add("Content-Type", "application/json")

		authHandler.Login(c)
		require.Equal(t, http.StatusOK, w.Result().StatusCode)

		var authResponse AuthResponse
		err := json.Unmarshal(w.Body.Bytes(), &authResponse)
		require.NoError(t, err)

		// Set up 2FA
		t.Run("Setup2FA", func(t *testing.T) {
			w = httptest.NewRecorder()
			c, _ = gin.CreateTestContext(w)

			c.Request = httptest.NewRequest(http.MethodPost, "/api/auth/2fa", bytes.NewBuffer([]byte(`{"enable": true}`)))
			c.Request.Header.Add("Content-Type", "application/json")
			c.Request.Header.Add("Authorization", "Bearer "+authResponse.AccessToken)

			authHandler.Setup2FA(c)
			require.Equal(t, http.StatusOK, w.Result().StatusCode)

			var setupResponse struct {
				Secret     string `json:"secret"`
				OtpauthURL string `json:"otpauth_url"`
			}
			err = json.Unmarshal(w.Body.Bytes(), &setupResponse)
			require.NoError(t, err)
			require.NotEmpty(t, setupResponse.Secret)
			require.NotEmpty(t, setupResponse.OtpauthURL)
		})

		t.Run("LoginWith2FAEnabled", func(t *testing.T) {
			// First, enable 2FA completely (this would normally require a valid TOTP token)
			user, err := userRepo.GetUserByEmail("2fa@example.com")
			require.NoError(t, err)

			// Simulate 2FA being enabled
			is2FAEnabled := true
			user.Is2FAEnabled = &is2FAEnabled
			err = userRepo.UpdateUser(user)
			require.NoError(t, err)

			// Try to login without 2FA token
			loginReq := LoginRequest{
				Email:    "2fa@example.com",
				Password: "SecureP@ssw0rd123",
				// No TwoFAToken provided
			}

			w = httptest.NewRecorder()
			c, _ = gin.CreateTestContext(w)

			jsonData, _ := json.Marshal(loginReq)
			c.Request = httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(jsonData))
			c.Request.Header.Add("Content-Type", "application/json")

			authHandler.Login(c)

			// Should require 2FA token
			require.Equal(t, http.StatusUnauthorized, w.Result().StatusCode)
			responseBody := w.Body.String()
			require.Contains(t, responseBody, "2FA token required")
			require.Contains(t, responseBody, "requires_2fa")
		})
	})
}

// TestAuthorizationSecurity tests various security aspects of the authorization system
func TestAuthorizationSecurity(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "authz-security-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Set environment variables for testing
	os.Setenv("JWT_SECRET", "test-super-secret-jwt-key-for-testing-only")
	defer os.Unsetenv("JWT_SECRET")

	// Initialize database
	database.ConnectToDB()
	defer func() {
		dbPath := filepath.Join(tempDir, database.DbFolder, database.DbName)
		if err := os.Remove(dbPath); err != nil && !os.IsNotExist(err) {
			t.Error(err)
		}
	}()

	// Create auth handler and middleware
	userRepo := database.NewUserRepo(database.DB)
	authHandler := NewAuthHandler(userRepo)
	authMiddleware := middleware.NewAuthMiddleware()

	// Create admin user
	adminReq := RegisterRequest{
		Email:    "admin@example.com",
		Password: "AdminP@ssw0rd123",
		Role:     "admin",
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	jsonData, _ := json.Marshal(adminReq)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBuffer(jsonData))
	c.Request.Header.Add("Content-Type", "application/json")

	authHandler.Register(c)
	require.Equal(t, http.StatusCreated, w.Result().StatusCode)

	// Login as admin
	adminLoginReq := LoginRequest{
		Email:    "admin@example.com",
		Password: "AdminP@ssw0rd123",
	}

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	jsonData, _ = json.Marshal(adminLoginReq)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(jsonData))
	c.Request.Header.Add("Content-Type", "application/json")

	authHandler.Login(c)
	require.Equal(t, http.StatusOK, w.Result().StatusCode)

	var adminAuthResponse AuthResponse
	err = json.Unmarshal(w.Body.Bytes(), &adminAuthResponse)
	require.NoError(t, err)

	// Create regular user
	userReq := RegisterRequest{
		Email:    "user@example.com",
		Password: "UserP@ssw0rd123",
		Role:     "user",
	}

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	jsonData, _ = json.Marshal(userReq)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBuffer(jsonData))
	c.Request.Header.Add("Content-Type", "application/json")

	authHandler.Register(c)
	require.Equal(t, http.StatusCreated, w.Result().StatusCode)

	// Login as regular user
	userLoginReq := LoginRequest{
		Email:    "user@example.com",
		Password: "UserP@ssw0rd123",
	}

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	jsonData, _ = json.Marshal(userLoginReq)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(jsonData))
	c.Request.Header.Add("Content-Type", "application/json")

	authHandler.Login(c)
	require.Equal(t, http.StatusOK, w.Result().StatusCode)

	var userAuthResponse AuthResponse
	err = json.Unmarshal(w.Body.Bytes(), &userAuthResponse)
	require.NoError(t, err)

	t.Run("RoleBasedAccessControl", func(t *testing.T) {
		t.Run("AdminAccessToAdminEndpoints", func(t *testing.T) {
			w = httptest.NewRecorder()
			c, _ = gin.CreateTestContext(w)

			c.Request = httptest.NewRequest(http.MethodGet, "/api/admin/users", nil)
			c.Request.Header.Add("Authorization", "Bearer "+adminAuthResponse.AccessToken)

			// Apply admin middleware
			handler := authMiddleware.RequireAdmin()
			handler(c)

			// Admin should be able to access admin endpoints
			require.Equal(t, http.StatusOK, w.Result().StatusCode)
		})

		t.Run("UserDeniedFromAdminEndpoints", func(t *testing.T) {
			w = httptest.NewRecorder()
			c, _ = gin.CreateTestContext(w)

			c.Request = httptest.NewRequest(http.MethodGet, "/api/admin/users", nil)
			c.Request.Header.Add("Authorization", "Bearer "+userAuthResponse.AccessToken)

			// Apply admin middleware
			handler := authMiddleware.RequireAdmin()
			handler(c)

			// Regular user should be denied access to admin endpoints
			require.Equal(t, http.StatusForbidden, w.Result().StatusCode)
			responseBody := w.Body.String()
			require.Contains(t, responseBody, "Insufficient permissions")
		})

		t.Run("SpecificRoleAccess", func(t *testing.T) {
			w = httptest.NewRecorder()
			c, _ = gin.CreateTestContext(w)

			c.Request = httptest.NewRequest(http.MethodGet, "/api/admin/users", nil)
			c.Request.Header.Add("Authorization", "Bearer "+userAuthResponse.AccessToken)

			// Apply role middleware requiring admin role
			handler := authMiddleware.RequireRole("admin")
			handler(c)

			// Regular user should be denied access when admin role is required
			require.Equal(t, http.StatusForbidden, w.Result().StatusCode)
			responseBody := w.Body.String()
			require.Contains(t, responseBody, "Insufficient permissions")
		})
	})

	t.Run("ContextSecurity", func(t *testing.T) {
		t.Run("UserContextSetCorrectly", func(t *testing.T) {
			w = httptest.NewRecorder()
			c, _ = gin.CreateTestContext(w)

			c.Request = httptest.NewRequest(http.MethodGet, "/api/auth/profile", nil)
			c.Request.Header.Add("Authorization", "Bearer "+userAuthResponse.AccessToken)

			// Apply auth middleware
			handler := authMiddleware.RequireAuth()
			handler(c)

			// Verify that user context is set correctly
			userID, exists := c.Get("user_id")
			require.True(t, exists)
			require.NotNil(t, userID)

			userEmail, exists := c.Get("user_email")
			require.True(t, exists)
			require.Equal(t, "user@example.com", userEmail)

			userRole, exists := c.Get("user_role")
			require.True(t, exists)
			require.Equal(t, "user", userRole)

			user, exists := c.Get("user")
			require.True(t, exists)
			require.NotNil(t, user)
			require.Equal(t, "user@example.com", user.(*models.User).Email)
		})

		t.Run("MissingUserContext", func(t *testing.T) {
			w = httptest.NewRecorder()
			c, _ = gin.CreateTestContext(w)

			c.Request = httptest.NewRequest(http.MethodGet, "/api/auth/profile", nil)
			// No authorization header

			// Apply auth middleware
			handler := authMiddleware.RequireAuth()
			handler(c)

			// Verify that user context is not set
			_, exists := c.Get("user_id")
			require.False(t, exists)

			_, exists = c.Get("user_email")
			require.False(t, exists)

			_, exists = c.Get("user_role")
			require.False(t, exists)

			_, exists = c.Get("user")
			require.False(t, exists)
		})
	})

	t.Run("OptionalAuthentication", func(t *testing.T) {
		t.Run("OptionalAuthWithToken", func(t *testing.T) {
			w = httptest.NewRecorder()
			c, _ = gin.CreateTestContext(w)

			c.Request = httptest.NewRequest(http.MethodGet, "/api/public/endpoint", nil)
			c.Request.Header.Add("Authorization", "Bearer "+userAuthResponse.AccessToken)

			// Apply optional auth middleware
			handler := authMiddleware.OptionalAuth()
			handler(c)

			// Should succeed and set user context
			require.Equal(t, http.StatusOK, w.Result().StatusCode)

			userID, exists := c.Get("user_id")
			require.True(t, exists)
			require.NotNil(t, userID)
		})

		t.Run("OptionalAuthWithoutToken", func(t *testing.T) {
			w = httptest.NewRecorder()
			c, _ = gin.CreateTestContext(w)

			c.Request = httptest.NewRequest(http.MethodGet, "/api/public/endpoint", nil)
			// No authorization header

			// Apply optional auth middleware
			handler := authMiddleware.OptionalAuth()
			handler(c)

			// Should succeed without setting user context
			require.Equal(t, http.StatusOK, w.Result().StatusCode)

			_, exists := c.Get("user_id")
			require.False(t, exists)
		})
	})
}

// TestInputValidationSecurity tests input validation security for authentication endpoints
func TestInputValidationSecurity(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "input-validation-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Initialize database
	database.ConnectToDB()
	defer func() {
		dbPath := filepath.Join(tempDir, database.DbFolder, database.DbName)
		if err := os.Remove(dbPath); err != nil && !os.IsNotExist(err) {
			t.Error(err)
		}
	}()

	// Create auth handler
	userRepo := database.NewUserRepo(database.DB)
	authHandler := NewAuthHandler(userRepo)

	t.Run("EmailValidation", func(t *testing.T) {
		testCases := []struct {
			name     string
			email    string
			expected int
		}{
			{"ValidEmail", "valid@example.com", http.StatusCreated},
			{"InvalidEmailFormat", "invalid-email", http.StatusBadRequest},
			{"EmptyEmail", "", http.StatusBadRequest},
			{"EmailWithSQLInjection", "test@example.com'; DROP TABLE users; --", http.StatusBadRequest},
			{"EmailWithXSS", "test@example.com<script>alert('xss')</script>", http.StatusBadRequest},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				registerReq := RegisterRequest{
					Email:    tc.email,
					Password: "SecureP@ssw0rd123",
				}

				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				jsonData, _ := json.Marshal(registerReq)
				c.Request = httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBuffer(jsonData))
				c.Request.Header.Add("Content-Type", "application/json")

				authHandler.Register(c)

				require.Equal(t, tc.expected, w.Result().StatusCode)
			})
		}
	})

	t.Run("PasswordValidation", func(t *testing.T) {
		testCases := []struct {
			name     string
			password string
			expected int
		}{
			{"ValidPassword", "SecureP@ssw0rd123", http.StatusCreated},
			{"TooShortPassword", "123", http.StatusBadRequest},
			{"EmptyPassword", "", http.StatusBadRequest},
			{"PasswordWithSQLInjection", "password'; DROP TABLE users; --", http.StatusBadRequest},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				registerReq := RegisterRequest{
					Email:    tc.name + "@example.com",
					Password: tc.password,
				}

				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				jsonData, _ := json.Marshal(registerReq)
				c.Request = httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBuffer(jsonData))
				c.Request.Header.Add("Content-Type", "application/json")

				authHandler.Register(c)

				require.Equal(t, tc.expected, w.Result().StatusCode)
			})
		}
	})

	t.Run("RoleValidation", func(t *testing.T) {
		testCases := []struct {
			name     string
			role     string
			expected int
		}{
			{"ValidRoleAdmin", "admin", http.StatusCreated},
			{"ValidRoleUser", "user", http.StatusCreated},
			{"InvalidRole", "hacker", http.StatusBadRequest},
			{"RoleWithSQLInjection", "admin'; DROP TABLE users; --", http.StatusBadRequest},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				registerReq := RegisterRequest{
					Email:    tc.name + "@example.com",
					Password: "SecureP@ssw0rd123",
					Role:     tc.role,
				}

				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				jsonData, _ := json.Marshal(registerReq)
				c.Request = httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBuffer(jsonData))
				c.Request.Header.Add("Content-Type", "application/json")

				authHandler.Register(c)

				require.Equal(t, tc.expected, w.Result().StatusCode)
			})
		}
	})

	t.Run("JSONInjection", func(t *testing.T) {
		// Test for JSON injection attacks
		maliciousJSON := `{
			"email": "test@example.com",
			"password": "SecureP@ssw0rd123",
			"role": "user",
			"malicious": {"$ref": "https://attacker.com/malicious.json"}
		}`

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBuffer([]byte(maliciousJSON)))
		c.Request.Header.Add("Content-Type", "application/json")

		authHandler.Register(c)

		// Should reject malformed JSON or ignore extra fields
		require.Equal(t, http.StatusCreated, w.Result().StatusCode)
	})
}

// TestTokenExpirationSecurity tests token expiration and refresh security
func TestTokenExpirationSecurity(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "token-expiry-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Set environment variables for testing
	os.Setenv("JWT_SECRET", "test-super-secret-jwt-key-for-testing-only")
	os.Setenv("JWT_EXPIRES_IN", "1") // 1 second for testing
	defer func() {
		os.Unsetenv("JWT_SECRET")
		os.Unsetenv("JWT_EXPIRES_IN")
	}()

	// Initialize database
	database.ConnectToDB()
	defer func() {
		dbPath := filepath.Join(tempDir, database.DbFolder, database.DbName)
		if err := os.Remove(dbPath); err != nil && !os.IsNotExist(err) {
			t.Error(err)
		}
	}()

	// Create auth handler and middleware
	userRepo := database.NewUserRepo(database.DB)
	authHandler := NewAuthHandler(userRepo)
	authMiddleware := middleware.NewAuthMiddleware()

	// Register and login a user
	registerReq := RegisterRequest{
		Email:    "expiry@example.com",
		Password: "SecureP@ssw0rd123",
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	jsonData, _ := json.Marshal(registerReq)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBuffer(jsonData))
	c.Request.Header.Add("Content-Type", "application/json")

	authHandler.Register(c)
	require.Equal(t, http.StatusCreated, w.Result().StatusCode)

	// Login to get tokens
	loginReq := LoginRequest{
		Email:    "expiry@example.com",
		Password: "SecureP@ssw0rd123",
	}

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	jsonData, _ = json.Marshal(loginReq)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(jsonData))
	c.Request.Header.Add("Content-Type", "application/json")

	authHandler.Login(c)
	require.Equal(t, http.StatusOK, w.Result().StatusCode)

	var authResponse AuthResponse
	err = json.Unmarshal(w.Body.Bytes(), &authResponse)
	require.NoError(t, err)

	t.Run("ExpiredTokenRejected", func(t *testing.T) {
		// Wait for token to expire
		time.Sleep(2 * time.Second)

		w = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(w)

		c.Request = httptest.NewRequest(http.MethodGet, "/api/auth/profile", nil)
		c.Request.Header.Add("Authorization", "Bearer "+authResponse.AccessToken)

		// Apply auth middleware
		handler := authMiddleware.RequireAuth()
		handler(c)

		// Should reject expired token
		require.Equal(t, http.StatusUnauthorized, w.Result().StatusCode)
		responseBody := w.Body.String()
		require.Contains(t, responseBody, "Invalid token")
	})

	t.Run("RefreshTokenExpiration", func(t *testing.T) {
		// Set refresh token expiration to a short time for testing
		os.Setenv("REFRESH_TOKEN_EXPIRES_IN", "1") // 1 second
		defer os.Unsetenv("REFRESH_TOKEN_EXPIRES_IN")

		// Create a new user for this test
		registerReq := RegisterRequest{
			Email:    "refresh-expiry@example.com",
			Password: "SecureP@ssw0rd123",
		}

		w = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(w)

		jsonData, _ := json.Marshal(registerReq)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBuffer(jsonData))
		c.Request.Header.Add("Content-Type", "application/json")

		authHandler.Register(c)
		require.Equal(t, http.StatusCreated, w.Result().StatusCode)

		// Login to get tokens
		loginReq := LoginRequest{
			Email:    "refresh-expiry@example.com",
			Password: "SecureP@ssw0rd123",
		}

		w = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(w)

		jsonData, _ = json.Marshal(loginReq)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(jsonData))
		c.Request.Header.Add("Content-Type", "application/json")

		authHandler.Login(c)
		require.Equal(t, http.StatusOK, w.Result().StatusCode)

		var refreshAuthResponse AuthResponse
		err = json.Unmarshal(w.Body.Bytes(), &refreshAuthResponse)
		require.NoError(t, err)

		// Wait for refresh token to expire
		time.Sleep(2 * time.Second)

		// Try to refresh with expired refresh token
		refreshReq := RefreshRequest{
			RefreshToken: refreshAuthResponse.RefreshToken,
		}

		w = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(w)

		jsonData, _ = json.Marshal(refreshReq)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/auth/refresh", bytes.NewBuffer(jsonData))
		c.Request.Header.Add("Content-Type", "application/json")

		authHandler.RefreshToken(c)

		// Should reject expired refresh token
		require.Equal(t, http.StatusUnauthorized, w.Result().StatusCode)
		responseBody := w.Body.String()
		require.Contains(t, responseBody, "Invalid refresh token")
	})
}

// TestPasswordSecurity tests password hashing and verification security
func TestPasswordSecurity(t *testing.T) {
	t.Run("PasswordHashing", func(t *testing.T) {
		testCases := []struct {
			name     string
			password string
		}{
			{"SimplePassword", "password123"},
			{"ComplexPassword", "Str0ngP@ssw0rd!"},
			{"VeryLongPassword", "ThisIsAVeryLongPasswordThatShouldBeSecureAndHardToCrack123!@#"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Create a user with the test password
				user := &models.User{
					Email: tc.name + "@example.com",
					Role:  "user",
				}

				// Hash the password
				err := user.HashPassword(tc.password)
				require.NoError(t, err)

				// Verify that the password is not stored in plaintext
				require.NotEqual(t, tc.password, user.PasswordHash)

				// Verify that the password hash is bcrypt format
				require.True(t, isBcryptHash(user.PasswordHash))

				// Verify that the password can be checked correctly
				require.True(t, user.CheckPassword(tc.password))

				// Verify that incorrect passwords are rejected
				require.False(t, user.CheckPassword("wrongpassword"))
			})
		}
	})

	t.Run("PasswordHashUniqueness", func(t *testing.T) {
		// Create two users with the same password
		user1 := &models.User{
			Email: "user1@example.com",
			Role:  "user",
		}

		user2 := &models.User{
			Email: "user2@example.com",
			Role:  "user",
		}

		password := "SamePassword123!"

		// Hash the same password for both users
		err := user1.HashPassword(password)
		require.NoError(t, err)

		err = user2.HashPassword(password)
		require.NoError(t, err)

		// Verify that the hashes are different (due to salt)
		require.NotEqual(t, user1.PasswordHash, user2.PasswordHash)

		// Verify that both passwords can be checked correctly
		require.True(t, user1.CheckPassword(password))
		require.True(t, user2.CheckPassword(password))
	})
}

// Helper function to check if a string is a bcrypt hash
func isBcryptHash(hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte("test"))
	return err == bcrypt.ErrMismatchedHashAndPassword || err == nil
}
