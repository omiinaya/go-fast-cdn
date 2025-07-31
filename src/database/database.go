package database

import (
	"fmt"
	"log"
	"os"

	"github.com/glebarez/sqlite"
	"github.com/kevinanielsen/go-fast-cdn/src/models"
	"github.com/kevinanielsen/go-fast-cdn/src/util"
	"gorm.io/gorm"
)

const (
	DbFolder = "db_data"
	DbName   = "main.db"
)

var DB *gorm.DB

func ConnectToDB() {
	dbPath := fmt.Sprintf("%v/%s", util.ExPath, DbFolder)

	_, err := os.Stat(fmt.Sprintf("%v/%s", dbPath, DbName))
	if err != nil {
		os.Mkdir(dbPath, 0o755)
		log.Printf("DB not found, creating at %v/main.db...", dbPath)
		os.Create(fmt.Sprintf("%v/main.db", dbPath))
	}

	database, err := gorm.Open(sqlite.Open(fmt.Sprintf("%v/main.db", dbPath)), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		panic("Failed to connect to database!")
	}
	log.Println("Connected to database!")

	database.AutoMigrate(&models.Media{}, &models.Config{})

	// Create composite unique index on FileName and Type for Media model
	mediaModel := models.Media{}
	if err := mediaModel.AddCompositeUniqueIndex(database); err != nil {
		log.Printf("Warning: Failed to create composite unique index for media model: %v", err)
		// Don't panic here, as the AddCompositeUniqueIndex function now handles duplicates gracefully
		// The application can continue to function even if index creation fails after cleanup attempts
	}

	DB = database
	log.Println("Database initialized!")
}
