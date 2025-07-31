package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/kevinanielsen/go-fast-cdn/src/database"
	ini "github.com/kevinanielsen/go-fast-cdn/src/initializers"
	"github.com/kevinanielsen/go-fast-cdn/src/router"
	"github.com/kevinanielsen/go-fast-cdn/src/util"
)

func init() {
	util.LoadExPath()
	gin.SetMode("release")
	ini.LoadEnvVariables(true)
	ini.CreateFolders()
	database.ConnectToDB()
	database.Migrate() // Run database migrations
}

// initializeApp performs synchronous initialization that must complete before the app starts
func initializeApp() error {
	// Initialize default configuration values
	if err := database.InitializeDefaultConfigs(); err != nil {
		return err
	}
	return nil
}

func main() {
	// Perform synchronous initialization before starting the server
	if err := initializeApp(); err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	log.Printf("Starting server on port %v", os.Getenv("PORT"))
	router.Router()
}
