package ui

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

// go:embed build
var staticFS embed.FS

// AddRoutes configures the router with middleware to serve the React frontend.
// In development mode, it serves files directly from the filesystem.
// In production mode, it serves files from the embedded build folder.
func AddRoutes(router gin.IRouter) {
	// Check if we're in development mode
	if isDevelopmentMode() {
		serveFromFilesystem(router)
	} else {
		serveFromEmbeddedFS(router)
	}
}

// isDevelopmentMode checks if we're running in development mode
func isDevelopmentMode() bool {
	// Check if the UI build directory exists in the filesystem
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	uiBuildPath := filepath.Join(basepath, "build")

	if _, err := os.Stat(uiBuildPath); os.IsNotExist(err) {
		return false
	}

	return true
}

// serveFromFilesystem serves the UI files directly from the filesystem
func serveFromFilesystem(router gin.IRouter) {
	// Get the path to the UI build directory
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	uiBuildPath := filepath.Join(basepath, "build")

	// Convert IRouter to Engine to access NoRoute method
	engine := router.(*gin.Engine)

	// Add a route to serve static files for specific paths
	engine.GET("/", func(c *gin.Context) {
		c.File(filepath.Join(uiBuildPath, "index.html"))
	})

	// Serve static assets (JS, CSS, images)
	engine.Any("/assets/*filepath", func(c *gin.Context) {
		assetPath := c.Param("filepath")
		// Remove leading slash from assetPath to prevent double slashes
		if len(assetPath) > 0 && assetPath[0] == '/' {
			assetPath = assetPath[1:]
		}
		// Add "assets" prefix since files are in ui/build/assets/
		fullPath := filepath.Join(uiBuildPath, "assets", assetPath)
		c.File(fullPath)
	})

	// Serve other static files
	engine.GET("/favicon.ico", func(c *gin.Context) {
		c.File(filepath.Join(uiBuildPath, "favicon.ico"))
	})

	// Add a fallback route to serve index.html for unknown routes (SPA support)
	// but exclude API routes
	engine.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		fmt.Printf("[DEBUG] NoRoute called for path: %s\n", path)

		// Don't serve index.html for API routes
		if strings.HasPrefix(path, "/api/") {
			fmt.Printf("[DEBUG] NoRoute: API endpoint not found: %s\n", path)
			c.JSON(404, gin.H{"error": "API endpoint not found"})
			return
		}

		// Serve index.html for all other routes (SPA support)
		fmt.Printf("[DEBUG] NoRoute: Serving index.html for path: %s\n", path)
		c.File(filepath.Join(uiBuildPath, "index.html"))
	})
}

// serveFromEmbeddedFS serves the UI files from the embedded filesystem
func serveFromEmbeddedFS(router gin.IRouter) {
	embeddedBuildFolder := newStaticFileSystem()
	log.Printf("Embedded filesystem created successfully")
	fallbackFileSystem := newFallbackFileSystem(embeddedBuildFolder)
	log.Printf("Fallback filesystem created successfully")
	router.Use(static.Serve("/", embeddedBuildFolder))
	router.Use(static.Serve("/", fallbackFileSystem))
}

// ----------------------------------------------------------------------
// staticFileSystem serves files out of the embedded build folder

type staticFileSystem struct {
	http.FileSystem
}

var _ static.ServeFileSystem = (*staticFileSystem)(nil)

func newStaticFileSystem() *staticFileSystem {
	sub, err := fs.Sub(staticFS, "build")
	if err != nil {
		panic(err)
	}

	return &staticFileSystem{
		FileSystem: http.FS(sub),
	}
}

func (s *staticFileSystem) Exists(prefix string, path string) bool {
	buildpath := fmt.Sprintf("build%s", path)

	// support for folders
	if strings.HasSuffix(path, "/") {
		_, err := staticFS.ReadDir(strings.TrimSuffix(buildpath, "/"))
		return err == nil
	}

	// support for files
	f, err := staticFS.Open(buildpath)
	if f != nil {
		_ = f.Close()
	}
	return err == nil
}

// fallbackFileSystem wraps a staticFileSystem and always serves /index.html
type fallbackFileSystem struct {
	staticFileSystem *staticFileSystem
}

var (
	_ static.ServeFileSystem = (*fallbackFileSystem)(nil)
	_ http.FileSystem        = (*fallbackFileSystem)(nil)
)

func newFallbackFileSystem(staticFileSystem *staticFileSystem) *fallbackFileSystem {
	return &fallbackFileSystem{
		staticFileSystem: staticFileSystem,
	}
}

func (f *fallbackFileSystem) Open(path string) (http.File, error) {
	return f.staticFileSystem.Open("/index.html")
}

func (f *fallbackFileSystem) Exists(prefix string, path string) bool {
	return true
}
