package router

import (
	"log"
	"os"

	"github.com/kevinanielsen/go-fast-cdn/src/database"
	"github.com/kevinanielsen/go-fast-cdn/src/middleware"
	"github.com/kevinanielsen/go-fast-cdn/ui"
)

// Router initializes the router and sets up middleware, routes, etc.
// It returns a *gin.Engine instance configured with the routes, middleware, etc.
func Router() {
	// Ensure configuration is initialized before setting up routes
	if err := database.InitializeDefaultConfigs(); err != nil {
		log.Fatalf("Failed to initialize default configurations: %v", err)
	}

	port := ":" + os.Getenv("PORT")

	s := NewServer(
		WithPort(port),
		WithMiddleware(middleware.CORSMiddleware()),
	)

	// Add all the API routes
	s.AddApiRoutes()

	// Add the embedded ui routes
	ui.AddRoutes(s.Engine)

	s.Run()
}
