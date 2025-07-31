package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kevinanielsen/go-fast-cdn/src/database"
	"github.com/kevinanielsen/go-fast-cdn/src/handlers"
	authHandlers "github.com/kevinanielsen/go-fast-cdn/src/handlers/auth"
	dbHandlers "github.com/kevinanielsen/go-fast-cdn/src/handlers/db"
	dHandlers "github.com/kevinanielsen/go-fast-cdn/src/handlers/docs"
	iHandlers "github.com/kevinanielsen/go-fast-cdn/src/handlers/image"
	mHandlers "github.com/kevinanielsen/go-fast-cdn/src/handlers/media"
	"github.com/kevinanielsen/go-fast-cdn/src/middleware"
	"github.com/kevinanielsen/go-fast-cdn/src/util"
)

func (s *Server) AddApiRoutes() {
	api := s.Engine.Group("/api")
	api.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, "pong")
	})

	// Authentication routes (public)
	authHandler := authHandlers.NewAuthHandler(database.NewUserRepo(database.DB))
	auth := api.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh", authHandler.RefreshToken)
		auth.POST("/logout", authHandler.Logout)
	}

	// Initialize auth middleware
	authMiddleware := middleware.NewAuthMiddleware()

	// Protected auth routes
	authProtected := api.Group("/auth")
	authProtected.Use(authMiddleware.RequireAuth())
	{
		authProtected.GET("/profile", authHandler.GetProfile)
		authProtected.PUT("/change-password", authHandler.ChangePassword)
		authProtected.PUT("/change-email", authHandler.ChangeEmail)
		authProtected.POST("/2fa", authHandler.Setup2FA)
		authProtected.POST("/2fa/verify", authHandler.Verify2FA)
	}

	cdn := api.Group("/cdn")
	docHandler := dHandlers.NewDocHandler(database.NewDocRepo(database.DB))
	imageHandler := iHandlers.NewImageHandler(database.NewImageRepo(database.DB))
	mediaHandler := mHandlers.NewMediaHandler(database.NewMediaRepo(database.DB))

	// Public CDN routes (read-only)
	{
		cdn.GET("/size", handlers.GetSizeHandler)

		// Unified media endpoints
		cdn.GET("/media/all", mediaHandler.HandleAllMedia)
		cdn.GET("/media/:filename", mediaHandler.HandleMediaMetadata)
		cdn.Static("/download/media", util.ExPath+"/uploads/media")
		cdn.Static("/download/images", util.ExPath+"/uploads/images")
		cdn.Static("/download/docs", util.ExPath+"/uploads/docs")
		cdn.Static("/uploads/media", util.ExPath+"/uploads/media")
		cdn.Static("/uploads/images", util.ExPath+"/uploads/images")
		cdn.Static("/uploads/docs", util.ExPath+"/uploads/docs")

		// Legacy endpoints for backward compatibility
		cdn.GET("/doc/all", docHandler.HandleAllDocs)
		cdn.GET("/doc/:filename", dHandlers.HandleDocMetadata)
		cdn.GET("/image/all", imageHandler.HandleAllImages)
		cdn.GET("/image/:filename", iHandlers.HandleImageMetadata)

		cdn.GET("/dashboard", handlers.NewDashboardHandler(
			database.NewDocRepo(database.DB),
			database.NewImageRepo(database.DB),
			database.NewUserRepo(database.DB),
			database.NewConfigRepo(database.DB),
		).GetDashboard)
	}

	// Protected CDN routes (require authentication)
	cdnProtected := cdn.Group("/")
	cdnProtected.Use(authMiddleware.RequireAuth())

	// Unified media endpoints
	upload := cdnProtected.Group("upload")
	{
		upload.POST("/media", mediaHandler.HandleMediaUpload)
		// Legacy endpoints for backward compatibility
		upload.POST("/image", mediaHandler.HandleImageUpload)
		upload.POST("/doc", mediaHandler.HandleDocUpload)
	}

	delete := cdnProtected.Group("delete")
	{
		delete.DELETE("/media/:filename", mediaHandler.HandleMediaDelete)
		// Legacy endpoints for backward compatibility
		delete.DELETE("/image/:filename", mediaHandler.HandleImageDelete)
		delete.DELETE("/doc/:filename", mediaHandler.HandleDocDelete)
	}

	rename := cdnProtected.Group("rename")
	{
		rename.PUT("/media", mediaHandler.HandleMediaRename)
		// Legacy endpoints for backward compatibility
		rename.PUT("/image", mediaHandler.HandleImageRename)
		rename.PUT("/doc", mediaHandler.HandleDocsRename)
	}

	resize := cdnProtected.Group("resize")
	{
		resize.PUT("/media", mediaHandler.HandleMediaResize)
		// Legacy endpoints for backward compatibility
		resize.PUT("/image", iHandlers.HandleImageResize)
	}
	// Admin-only routes
	adminRoutes := api.Group("/admin")
	adminRoutes.Use(authMiddleware.RequireAuth(), authMiddleware.RequireAdmin())
	{
		adminRoutes.POST("/drop/database", dbHandlers.HandleDropDB)

		adminUserHandler := authHandlers.NewAdminUserHandler(database.NewUserRepo(database.DB))
		{
			adminRoutes.GET("/users", adminUserHandler.ListUsers)
			adminRoutes.POST("/users", adminUserHandler.CreateUser)
			adminRoutes.PUT("/users/:id", adminUserHandler.UpdateUser)
			adminRoutes.DELETE("/users/:id", adminUserHandler.DeleteUser)
		}

		// Config endpoints (admin only)
		configHandler := handlers.NewConfigHandler(database.NewConfigRepo(database.DB))
		adminRoutes.GET("/config/registration", configHandler.GetRegistrationEnabled)
		adminRoutes.POST("/config/registration", configHandler.SetRegistrationEnabled)
	}

	// Public config endpoint for registration status
	configHandler := handlers.NewConfigHandler(database.NewConfigRepo(database.DB))
	api.GET("/config/registration", configHandler.GetRegistrationEnabled)
}
