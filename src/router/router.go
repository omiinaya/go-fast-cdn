package router

import (
	"os"

	"github.com/kevinanielsen/go-fast-cdn/src/middleware"
	"github.com/kevinanielsen/go-fast-cdn/src/util"
	"github.com/kevinanielsen/go-fast-cdn/ui"
)

// Router initializes the router and sets up middleware, routes, etc.
// It returns a *gin.Engine instance configured with the routes, middleware, etc.
func Router() {
	port := ":" + os.Getenv("PORT")

	s := NewServer(
		WithPort(port),
		WithMiddleware(middleware.CORSMiddleware()),
	)

	// Add all the API routes
	s.AddApiRoutes()

	// Add static file serving for uploads
	s.Engine.Static("/uploads/media", util.ExPath+"/uploads/media")
	s.Engine.Static("/uploads/images", util.ExPath+"/uploads/images")
	s.Engine.Static("/uploads/docs", util.ExPath+"/uploads/docs")
	// Add the embedded ui routes
	ui.AddRoutes(s.Engine)

	s.Run()
}
