package initializers

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// LoadEnvVariables loads environment variables from .env file or sets
// hardcoded values based on prod boolean. In prod it sets PORT and DB_SECRET
// to hardcoded values. In dev it loads .env file from current directory.
func LoadEnvVariables(prod bool) {
	if prod {
		os.Setenv("PORT", "8080")
		os.Setenv("DB_SECRET", "secret")
	} else {
		if err := godotenv.Load(); err != nil {
			log.Printf("Warning: .env file not found, using default values: %s", err.Error())
			// Set default PORT for development if not set
			if os.Getenv("PORT") == "" {
				os.Setenv("PORT", "8080")
			}
		}
	}
}
