package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/kevinanielsen/go-fast-cdn/src/database"
	"github.com/kevinanielsen/go-fast-cdn/src/util"
)

func main() {
	// Load executable path
	util.LoadExPath()

	// Create the database directory in the project root
	dbPath := fmt.Sprintf("%v/%s", util.ExPath, database.DbFolder)
	if err := os.MkdirAll(dbPath, 0755); err != nil {
		log.Fatalf("Failed to create database directory: %v", err)
	}

	// Create the database file
	dbFile := filepath.Join(dbPath, database.DbName)
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		file, err := os.Create(dbFile)
		if err != nil {
			log.Fatalf("Failed to create database file: %v", err)
		}
		file.Close()
		log.Printf("Database file created at: %s", dbFile)
	}

	// Connect to database (this will initialize the schema)
	database.ConnectToDB()

	log.Println("Database initialized successfully")
}
